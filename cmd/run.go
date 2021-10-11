package cmd

import (
	"github.com/julienschmidt/httprouter"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"my-instance/controller"
	"my-instance/db"
	"my-instance/storage"
	"net"
	"net/http"

	"golang.org/x/net/context"
	"my-instance/config"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"
)

func run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	ctx, cancle := context.WithCancel(ctx)
	defer cancle()
	tasks := []func() error{
		setLogLevel,
		setMysqlSQLConnection,
		runDatabaseMigrations,
		setRedisPool,
		setMongoClient,
		setHttpSever,
	}

	for _, t := range tasks {
		if err := t(); err != nil {
			log.Fatal(err)
		}
	}
	sigChan := make(chan os.Signal)
	exitChan := make(chan struct{})
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	log.WithField("signal", <-sigChan).Info("signal received")
	go func() {
		log.Warning("stopping my-instance")
		exitChan <- struct{}{}
	}()
	select {
	case <-exitChan:
	case s := <-sigChan:
		log.WithField("signal", s).Info("signal received, stopping immediately")
	}

	return nil
}

func setLogLevel() error {
	log.SetFormatter(
		&log.TextFormatter{
			DisableTimestamp:       false,
			DisableColors:          false,
			DisableLevelTruncation: false,
			DisableSorting:         true,
			FullTimestamp:          true,
			ForceColors:            true,
			QuoteEmptyFields:       true,
		},
	)
	log.SetLevel(log.Level(uint8(config.C.General.LogLevel)))
	// 开启日志写入文件
	//if config.C.General.LogPath != "" {
	//	configLocalFilesystemLogger(config.C.General.LogPath, "my-instance", time.Hour*24, time.Hour*24)
	//}
	log.Info("上线版本:%s", version)

	return nil
}

// 连接mysql
func setMysqlSQLConnection() error {
	log.Info("connecting to mysql")
	if config.C.MySQL.MaxIdle != 0 {
		db.MysqlMaxIdle = config.C.MySQL.MaxIdle
	}
	if config.C.MySQL.MaxIdle != 0 {
		db.MysqlMaxOpen = config.C.MySQL.MaxIdle
	}

	timezone := "Asia%2FShanghai"
	if config.C.MySQL.Timezone != "" {
		timezone = config.C.MySQL.Timezone
	}
	dsn := config.C.MySQL.User + ":" + config.C.MySQL.Password + "@tcp(" + config.C.MySQL.Host + ":" + config.C.MySQL.Port + ")/" + config.C.MySQL.Database + "?charset=utf8mb4&loc=" + timezone + "&parseTime=true"
	pool, err := db.OpenDatabase(dsn)
	if err != nil {
		return errors.Wrap(err, "database connection error")
	}
	config.C.MySQL.DB = pool
	return nil
}

// gorm 自动建表
func runDatabaseMigrations() error {
	// 全局设置表名不可以为复数形式
	config.C.MySQL.DB.SingularTable(true)
	config.C.MySQL.DB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(
		&storage.Account{},
		&storage.Device{},
	)
	return nil
}

// setRedisPool init redis
func setRedisPool() error {
	if config.C.Redis.URL != "" {
		log.WithField("RedisHost", config.C.Redis.URL).Info("setup redis connection pool")
		if config.C.Redis.MaxIdle != 0 {
			db.RedisMaxIdle = config.C.Redis.MaxIdle
		}
		if config.C.Redis.MaxActive != 0 {
			db.RedisMaxActive = config.C.Redis.MaxActive
		}
		config.C.Redis.Pool = db.NewRedisPool(config.C.Redis.URL)
		return config.C.Redis.Pool.TestOnBorrow(config.C.Redis.Pool.Get(), time.Now())
	}
	return nil
}

// setMongoSession init mongodb
func setMongoClient() error {
	var err error
	if config.C.Mongo.Conn == ""{
		log.Info("mongo connection  is nil")
		return nil
	}
	if config.C.Mongo.Conn != "" {
		log.WithField("MongoHost", config.C.Mongo.Conn).Info("connecting to mongo ")
		ctx ,cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		config.C.Mongo.Client, err = db.InitMongoClient(ctx)
		return err
	}
	return nil
}

func setHttpSever() error {
	log.Infof("start http server listen address: %v", config.C.General.Listen)
	router := httprouter.New()
	controller.InitRouter(router)
	host, port, err := net.SplitHostPort(config.C.General.Listen)
	if err != nil {
		log.Warnf("ip:port is error %s,%v", config.C.General.Listen, err)
		return err
	}
	address := net.JoinHostPort(host, port)
	h := &http.Server{
		Addr:         address,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
		Handler:      router,
	}
	go func() {
		for {
			err := h.ListenAndServe()
			if err != nil {
				log.Errorf("start http server %v", err)
				continue
			}
			break
		}
	}()
	return nil
}

// 配置日志文件输出
func configLocalFilesystemLogger(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) {
	var exist bool
	if _, err := os.Stat(logPath); err == nil {
		exist = true
	}
	if !exist {
		if err := os.Mkdir(logPath, os.ModePerm); err != nil {
			log.Errorf("config local file system logger error. %+v", errors.WithStack(err))
		}
	}
	baseLogPath := path.Join(logPath, logFileName)
	writerDebug, err := rotatelogs.New(
		baseLogPath+".debug.%Y%m%d.log",
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		log.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}
	writerInfo, err := rotatelogs.New(
		baseLogPath+".info.%Y%m%d.log",
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		log.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}
	writerWarn, err := rotatelogs.New(
		baseLogPath+".warn.%Y%m%d.log",
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		log.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}
	writerError, err := rotatelogs.New(
		baseLogPath+".error.%Y%m%d.log",
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		log.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writerDebug, // 为不同级别设置不同的输出目的
		log.InfoLevel:  writerInfo,
		log.WarnLevel:  writerWarn,
		log.ErrorLevel: writerError,
		log.FatalLevel: writerError,
		log.PanicLevel: writerError,
	}, &log.JSONFormatter{})
	log.AddHook(lfHook)
}
