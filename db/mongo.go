package db

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"my-instance/config"
	"time"
)

var (
	connectTimeout = 10 * time.Second
	socketTimeout  = 15 * time.Second
)

func InitMongoClient(ctx context.Context) (*mongo.Client, error) {
	opts := options.Client().
		ApplyURI(config.C.Mongo.Conn).
		SetConnectTimeout(connectTimeout).
		SetSocketTimeout(socketTimeout)

	if config.C.Mongo.Password != "" {
		opts.SetAuth(options.Credential{
			AuthMechanism: "SCRAM-SHA-1",
			AuthSource:    config.C.Mongo.Database, //my-instance
			Username:      config.C.Mongo.UserName,
			Password:      config.C.Mongo.Password})
	}
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("Mongo connection error:%s ", err)
	}
	for {
		if err = client.Ping(ctx, nil); err != nil {
			log.WithError(err).Error("mongodb: connecting to broker error, will retry in 1s")
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}
	return client, nil
}
