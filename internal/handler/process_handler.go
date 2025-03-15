package handler

import (
	"context"
	"net/http"
	"swift/internal/database"
	"swift/internal/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func PostAll(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		parsedData, exists := c.Get("parsedData")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No data to save"})
			return
		}
		var docs []interface{}
		for i, swift := range parsedData.([]model.SwiftData) {
			if i == 0 {
				continue
			}
			docs = append(docs, swift)
		}
		_, err := db.GetCollection().InsertMany(context.TODO(), docs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save data to database"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Data saved to MongoDB"})
	}
}

func GetAll(db database.Service) gin.HandlerFunc {

	return func(c *gin.Context) {
		cursor, err := db.GetCollection().Find(context.TODO(), bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch SWIFT codes"})
			return
		}
		defer cursor.Close(context.TODO())
		var results []model.SwiftData
		if err = cursor.All(context.TODO(), &results); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding results"})
			return
		}
		c.JSON(http.StatusOK, results)
	}
}

func DeleteAll(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		results, err := db.GetCollection().DeleteMany(context.TODO(), bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete all data "})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message":       "Successfully deleted documents",
			"deleted_count": results.DeletedCount,
		})
	}
}
