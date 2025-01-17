package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDB struct {
	*mongo.Database
}

var (
	instance *MongoDB
	once     sync.Once
)

func NewMongo() *MongoDB {
	once.Do(func() {
		host := os.Getenv("MONGO_HOST")
		port := os.Getenv("MONGO_PORT")
		user := os.Getenv("MONGO_USER")
		pass := os.Getenv("MONGO_PASS")
		dbName := os.Getenv("MONGO_DB")

		mongoApi := options.ServerAPI(options.ServerAPIVersion1)
		url := fmt.Sprintf("mongodb://%s:%s", host, port)

		crets := options.Credential{
			AuthSource: "admin",
			Username:   user,
			Password:   pass,
		}

		opts := options.Client().ApplyURI(url)
		opts.SetAuth(crets)
		opts.SetServerAPIOptions(mongoApi)
		opts.SetReplicaSet("mongors")
		opts.SetDirect(true)

		client, err := mongo.Connect(opts)
		if err != nil {
			log.Fatal(err)
		}

		instance = &MongoDB{client.Database(dbName)}
	})

	return instance
}
