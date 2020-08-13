package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//Database serves to connect and manage mongodb driver
type Database struct {
	Client *mongo.Client
	Root   *mongo.Database
	Ctx    context.Context
}

//NewDatabase creates a new MongoDB instance
func NewDatabase() *Database {
	return &Database{}
}

//Connect connects to the mongo online db
func (db *Database) Connect(uri, name string) context.CancelFunc {

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	err = client.Connect(ctx)

	if err != nil {
		panic(err)
	}



	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	db.Ctx = ctx
	db.Client = client
	db.Root = client.Database(name)

	fmt.Println("Connected to MongoDB Cloud Cluster")

	return cancel
}
