package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"swift/internal/database"
)

// Server represents the HTTP server and holds configuration related to the server
type Server struct {
	port int
	db   database.Service
}

// NewServer initializes a new instance of the server with necessary configurations and routes.
func NewServer() *http.Server {

	// Read the PORT environment variable
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	// Create a new Server instance
	NewServer := &Server{
		port: port,
		db:   database.New(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
