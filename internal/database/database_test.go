package database

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
)

// Sets up the container with the specified environment variables and exposed ports.
func mustStartMongoContainer() (func(context.Context, ...testcontainers.TerminateOption) error, error) {
	dbContainer, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mongo:latest",
			Env:          map[string]string{"MONGO_INITDB_ROOT_USERNAME": "melkey", "MONGO_INITDB_ROOT_PASSWORD": "password1234"},
			ExposedPorts: []string{"27017/tcp"},
		},
		Started: true,
	})
	if err != nil {
		return nil, err
	}
	dbHost, err := dbContainer.Host(context.Background())
	if err != nil {
		return dbContainer.Terminate, err
	}
	dbPort, err := dbContainer.MappedPort(context.Background(), "27017/tcp")
	if err != nil {
		return dbContainer.Terminate, err
	}
	time.Sleep(5 * time.Second)
	host = dbHost
	port = dbPort.Port()
	return dbContainer.Terminate, err
}

// Entry point for tests, starts the MongoDB container and ensures it's cleaned up afterward.
func TestMain(m *testing.M) {
	teardown, err := mustStartMongoContainer()
	if err != nil {
		log.Fatalf("could not start mongodb container: %v", err)
	}
	m.Run()
	if teardown != nil && teardown(context.Background()) != nil {
		log.Fatalf("could not teardown mongodb container: %v", err)
	}
}

// Ensures that the New function creates a valid service instance.
func TestNew(t *testing.T) {
	srv := New()
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

// Checks that the Health function returns the expected "It's healthy" message.
func TestHealth(t *testing.T) {
	srv := New()
	stats := srv.Health()
	if stats["message"] != "It's healthy" {
		t.Fatalf("expected message to be 'It's healthy', got %s", stats["message"])
	}
}
