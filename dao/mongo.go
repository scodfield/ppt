package dao

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"sync"
)

var (
	MongoClient *mongo.Client
	mongoOnce   sync.Once
)

func InitMongo(url string, secondaryPreferred bool) error {
	var err error
	mongoOnce.Do(func() {
		opts := options.Client().ApplyURI(url)
		host, _ := os.Hostname()
		opts.SetAppName(host + "-" + fmt.Sprintf("%d", os.Getpid()))
		if secondaryPreferred {
			opts.SetReadPreference(readpref.SecondaryPreferred())
		} else {
			opts.SetReadPreference(readpref.Primary())
		}
		client, err := mongo.Connect(Ctx, opts)
		if err != nil {
			panic(err)
		}
		err = client.Ping(Ctx, readpref.Primary())
		if err != nil {
			panic(err)
		}
		MongoClient = client
	})
	return err
}

func CloseMongo() {
	if MongoClient != nil {
		_ = MongoClient.Disconnect(Ctx)
	}
}
