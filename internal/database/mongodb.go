package database

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

var (
	mongoInstance *MongoDB
	mongoOnce     sync.Once
	mongoInitErr  error
)

// NewMongoDB returns a singleton MongoDB connection.
// Subsequent calls reuse the same client instance.
func NewMongoDB(uri string) *MongoDB {
	mongoOnce.Do(func() {
		log.Println("[MongoDB] Initializing MongoDB connection...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		log.Printf("[MongoDB] Connecting to MongoDB...")
		clientOptions := options.Client().ApplyURI(uri).
			SetConnectTimeout(30 * time.Second).
			SetServerSelectionTimeout(30 * time.Second).
			SetSocketTimeout(30 * time.Second)
		
		client, err := mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Printf("[MongoDB] Error connecting to MongoDB: %v", err)
			mongoInitErr = err
			return
		}
		log.Println("[MongoDB] MongoDB client connected successfully")

		// Ping the database to verify connection
		log.Println("[MongoDB] Pinging database to verify connection...")
		pingCtx, pingCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer pingCancel()
		if err := client.Ping(pingCtx, nil); err != nil {
			log.Printf("[MongoDB] Error pinging database: %v", err)
			mongoInitErr = err
			return
		}
		log.Println("[MongoDB] Database ping successful")

		database := client.Database("astroway")
		log.Printf("[MongoDB] Using database: %s", database.Name())
		mongoInstance = &MongoDB{
			Client:   client,
			Database: database,
		}
		log.Println("[MongoDB] MongoDB instance initialized successfully")
	})

	if mongoInitErr != nil {
		log.Printf("[MongoDB] Fatal error initializing MongoDB: %v", mongoInitErr)
		panic(mongoInitErr)
	}

	return mongoInstance
}

func (m *MongoDB) GetCollection(name string) *mongo.Collection {
	return m.Database.Collection(name)
}

func (m *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.Client.Disconnect(ctx)
}
