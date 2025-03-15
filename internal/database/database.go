package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Defines the methods for interacting with MongoDB
type Service interface {
	Health() map[string]string
	GetCollection() *mongo.Collection
}

// Holds the MongoDB client and implements the Service interface
type service struct {
	db *mongo.Client
}

// Store the MongoDB host and port
var (
	host = os.Getenv("BLUEPRINT_DB_HOST")
	port = os.Getenv("BLUEPRINT_DB_PORT")
	// username = os.Getenv("BLUEPRINT_DB_USERNAME")
	// password = os.Getenv("BLUEPRINT_DB_ROOT_PASSWORD")
)

// Initializes a MongoDB client and returns a new Service instance
func New() Service {
	credential := options.Credential{Username: "melkey", Password: "password1234"}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", host, port)).SetAuth(credential))
	if err != nil {
		log.Fatal(err)
	}
	return &service{db: client}
}

// Checks the connection to MongoDB
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := s.db.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("db down: %v", err)
	}
	return map[string]string{"message": "It's healthy"}
}

// Retrieves a MongoDB collection by name
func (s *service) GetCollection() *mongo.Collection {
	return s.db.Database("swift_db").Collection("swift_data")
}
