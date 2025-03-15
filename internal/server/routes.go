package server

import (
	"net/http"
	"swift/internal/handler"
	"swift/internal/middleware"
	"swift/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

// RegisterRoutes defines all the routes and their associated handlers for the server.
func (s *Server) RegisterRoutes() http.Handler {

	// Create a new Gin router
	r := gin.Default()

	// Set up CORS middleware to handle cross-origin requests.
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// Health check and basic routes group
	healthGroup := r.Group("/")
	{
		// Basic hello world response
		healthGroup.GET("/", handler.HelloWorldHandler())
		// Health status, checks if the server and database are functioning correctly
		healthGroup.GET("/health", handler.HealthHandler(s.db))
	}

	// Group for processing data and show/delete all
	processGroup := r.Group("/my")
	{
		// Download, process and insert data
		processGroup.POST("/insertAll/:fileID", middleware.DownloadAsCSV(service.NewGoogleDocsService()), middleware.ParseCSV, handler.PostAll(s.db))
		// Get all data
		processGroup.GET("/getAll", handler.GetAll(s.db))
		// Delete all data
		processGroup.DELETE("/deleteAll", handler.DeleteAll(s.db))
	}

	// Group for handling SWIFT code-related endpoints
	swiftGroup := r.Group("/v1/swift-codes")
	{
		// Endpoint 1: Retrieve details of a single SWIFT code whether for a headquarters or branches
		swiftGroup.GET("/:swift-code", handler.GetSwift(s.db))
		// Endpoint 2: Return all SWIFT codes with details for a specific country (both headquarters and branches)
		swiftGroup.GET("/countries/:countryISO2code", handler.GetCountry(s.db))
		// Endpoint 3: Adds new SWIFT code entries to the database for a specific country
		swiftGroup.POST("/", handler.PostSwift(s.db))
		// Endpoint 4: Deletes swift-code data if swiftCode matches the one in the database
		swiftGroup.DELETE("/:swift-code", handler.DeleteSwift(s.db))
	}
	return r
}
