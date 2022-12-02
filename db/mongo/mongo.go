package mongodb

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/franciscobonand/seq-matrix/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	uriEnv      = "MONGO_URI"
	databaseEnv = "MONGO_DATABASE"
	collEnv     = "MONGO_COLLECTION"
)

type Mongo struct {
	ctx      context.Context
	db       *mongo.Database
	collName string
}

func Init(ctx context.Context) (db.Database, error) {
	db, ok := os.LookupEnv(databaseEnv)
	if !ok {
		return nil, fmt.Errorf("envvar '%s' not found", databaseEnv)
	}

	coll, ok := os.LookupEnv(collEnv)
	if !ok {
		return nil, fmt.Errorf("envvar '%s' not found", databaseEnv)
	}

	uri, ok := os.LookupEnv(uriEnv)
	if !ok {
		return nil, fmt.Errorf("envvar '%s' not found", uriEnv)
	}

	client, err := openConnection(ctx, uri)
	if err != nil {
		return nil, err
	}

	return &Mongo{
		ctx:      ctx,
		db:       client.Database(db),
		collName: coll,
	}, nil
}

func (m *Mongo) coll() *mongo.Collection {
	return m.db.Collection(m.collName)
}

// Get returns the total number of items registered, and how many of them are valid
func (m *Mongo) Get() (int64, int64, error) {
	var document bson.M
	searchParam := bson.D{{Key: "isValid", Value: true}}

	collStats := m.db.RunCommand(m.ctx, bson.M{"collStats": m.collName})
	err := collStats.Decode(&document)
	if err != nil {
		return -1, -1, fmt.Errorf("failed to get collection stats: %v", err)
	}

	totalItems, ok := document["count"].(int32)
	if !ok {
		return -1, -1, fmt.Errorf("failed to parse collection count")
	}

	valids, err := m.coll().CountDocuments(m.ctx, searchParam)
	if err != nil {
		return -1, -1, fmt.Errorf("failed to get valid sequences: %v", err)
	}

	return int64(totalItems), valids, nil
}

// Set creates a new entry on the database, with an unique ID
func (m *Mongo) Set(seq []string, valid bool) error {
	doc := bson.D{
		{Key: "sequence", Value: seq},
		{Key: "isValid", Value: valid},
	}

	_, err := m.coll().InsertOne(m.ctx, doc)

	return err
}

func openConnection(ctx context.Context, uri string) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(uri)
	c, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	client, err := mongo.Connect(c, opts)
	if err != nil {
		return nil, fmt.Errorf("error connecting to mongo: %v", err)
	}

	if err = client.Ping(c, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("mongo connection timeout: %v", err)
	}

	return client, nil
}
