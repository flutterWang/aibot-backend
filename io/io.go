package bot_io

import (
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"go.mongodb.org/mongo-driver/bson"

	"context"
	"time"
)

type Io struct {
	MongoClient *mongo.Database
}

func NewIo() *Io {
	return &Io{}
}

func (io *Io) GetMongoClient() *mongo.Database {
	return io.MongoClient
}

func (io *Io) RegisterMongoClient(uri, dbName string, connectTimeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()
	if client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri)); err != nil {
		return err
	} else if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return err
	} else {
		io.MongoClient = client.Database(dbName)
		return nil
	}
}

func (io *Io) CreateIndex(collName string, indexName string, unique bool) error {
	coll := io.MongoClient.Collection(collName)
	index := mongo.IndexModel{
		Keys:    bson.M{indexName: -1},
		Options: options.Index().SetUnique(unique),
	}
	_, err := coll.Indexes().CreateOne(context.Background(), index)
	return err
}
