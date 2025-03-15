package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"swift/internal/server"
)

// Listens for OS interrupt signals (SIGINT, SIGTERM) and allows the server to finish processing requests before shutting down
func gracefulShutdown(apiServer *http.Server, done chan bool) {

	// Create context that listens for the interrupt signal from the OS
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Wait for the interrupt signal to be received.
	<-ctx.Done()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// Create a new context with a 5-second timeout to allow the server to finish any active requests.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to shut down the server gracefully using the context.
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {

	// Create new server
	server := server.NewServer()

	// Create a done channel to signal when the graceful shutdown is complete
	done := make(chan bool, 1)

	// Run the gracefulShutdown function in a separate goroutine to allow the server to continue running.
	go gracefulShutdown(server, done)

	// Start the HTTP server and begin listening for requests.
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	// Wait for the graceful shutdown to complete
	<-done

	log.Println("Graceful shutdown complete.")
}
