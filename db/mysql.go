package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"time"
)

// mysql配置
var (
	MysqlMaxIdle = 5
	MysqlMaxOpen = 10
)

func OpenDatabase(dsn string) (*gorm.DB, error) {
	log.WithField("dsn", dsn).Info("OpenDatabase")
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.WithError(err).Error("database connection error")
		return nil, fmt.Errorf("database connection error: %s", err)
	}
	for {
		if err = db.DB().Ping(); err != nil {
			log.WithError(err).Error("ping database error, will retry in 2s")
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}
	db.LogMode(true)
	db.DB().SetMaxIdleConns(MysqlMaxIdle)
	db.DB().SetMaxOpenConns(MysqlMaxOpen)
	db.DB().SetConnMaxLifetime(6 * time.Hour)
	return db, nil
}
