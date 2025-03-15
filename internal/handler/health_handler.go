package handler

import (
	"net/http"
	"swift/internal/database"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func HelloWorldHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello World"})
	}
}

func HealthHandler(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, db.Health())
	}
}
