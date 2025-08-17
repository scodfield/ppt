package dao

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"os"
	"ppt/log"
	"ppt/nacos/wrapper"
	"sync"
)

var (
	MongoClient *mongo.Client
	mongoOnce   sync.Once
)

func InitMongo(nacosDBCfg *wrapper.DBConfig) error {
	var err error
	mongoOnce.Do(func() {
		MongoClient, err = initMongoByUrl(nacosDBCfg.MongoConfig, false)
		if err != nil {
			panic(err)
		}
	})
	return err
}

func initMongoByUrl(url string, isSecondPreferred bool) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(url)
	host, _ := os.Hostname()
	opts.SetAppName(host + "-" + fmt.Sprintf("%d", os.Getpid()))
	if isSecondPreferred {
		opts.SetReadPreference(readpref.SecondaryPreferred()) // 从库优先
	} else {
		opts.SetReadPreference(readpref.PrimaryPreferred())
	}
	client, err := mongo.Connect(Ctx, opts)
	if err != nil {
		log.Error("initMongoByUrl connect error", zap.String("mongo_url", url), zap.Bool("second_preferred", isSecondPreferred), zap.Error(err))
		return nil, err
	}
	err = client.Ping(Ctx, readpref.PrimaryPreferred())
	if err != nil {
		log.Error("initMongoByUrl ping error", zap.String("mongo_url", url), zap.Error(err))
		return nil, err
	}
	return client, nil
}

func CloseMongo() {
	if MongoClient != nil {
		_ = MongoClient.Disconnect(Ctx)
		MongoClient = nil
	}
}
