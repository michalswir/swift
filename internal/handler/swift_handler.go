package handler

import (
	"context"
	"net/http"
	"swift/internal/database"
	"swift/internal/model"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type SwiftResponse struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	CountryName   string `json:"countryName"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}

type SwiftResponseBranches struct {
	Address       string   `json:"address"`
	BankName      string   `json:"bankName"`
	CountryISO2   string   `json:"countryISO2"`
	CountryName   string   `json:"countryName"`
	IsHeadquarter bool     `json:"isHeadquarter"`
	SwiftCode     string   `json:"swiftCode"`
	Branches      []Branch `json:"branches"`
}

type Branch struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}

type Country struct {
	CountryISO2 string       `json:"countryISO2"`
	CountryName string       `json:"countryName"`
	SwiftCodes  []SwiftCodes `json:"swiftCodes"`
}

type SwiftCodes struct {
	Address       string `json:"address"`
	BankName      string `json:"bankName"`
	CountryISO2   string `json:"countryISO2"`
	IsHeadquarter bool   `json:"isHeadquarter"`
	SwiftCode     string `json:"swiftCode"`
}

func GetSwift(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		swiftCode := c.Param("swift-code")
		var result model.SwiftData
		err := db.GetCollection().FindOne(context.TODO(), bson.M{"swiftCode": swiftCode}).Decode(&result)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "SWIFT Code not found"})
			return
		}
		if result.IsHeadquarter {
			response := SwiftResponseBranches{
				Address:       result.Address,
				BankName:      result.BankName,
				CountryISO2:   result.CountryISO2,
				CountryName:   result.CountryName,
				IsHeadquarter: result.IsHeadquarter,
				SwiftCode:     swiftCode,
				Branches:      []Branch{},
			}
			cursor, err := db.GetCollection().Find(context.TODO(), bson.M{"bankName": result.BankName})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch branches"})
				return
			}
			defer cursor.Close(context.TODO())
			for cursor.Next(context.TODO()) {
				var branch model.SwiftData
				if err := cursor.Decode(&branch); err != nil {
					continue
				}
				if branch.SwiftCode != response.SwiftCode {
					response.Branches = append(response.Branches, Branch{
						Address:     branch.Address,
						BankName:    branch.BankName,
						CountryISO2: branch.CountryISO2,
						SwiftCode:   branch.SwiftCode,
					})
				}
			}
			c.JSON(http.StatusOK, response)
		} else {
			response := SwiftResponse{
				Address:       result.Address,
				BankName:      result.BankName,
				CountryISO2:   result.CountryISO2,
				CountryName:   result.CountryName,
				IsHeadquarter: result.IsHeadquarter,
				SwiftCode:     swiftCode,
			}
			c.JSON(http.StatusOK, response)
		}

	}
}

func GetCountry(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		CountryISO2 := c.Param("countryISO2code")
		cursor, err := db.GetCollection().Find(context.TODO(), bson.M{"countryIso2": CountryISO2})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		defer cursor.Close(context.TODO())
		var response Country
		found := false
		for cursor.Next(context.TODO()) {
			var branch model.SwiftData
			if err := cursor.Decode(&branch); err != nil {
				continue
			}
			if !found {
				response.CountryISO2 = branch.CountryISO2
				response.CountryName = branch.CountryName
				response.SwiftCodes = []SwiftCodes{}
				found = true
			}
			response.SwiftCodes = append(response.SwiftCodes, SwiftCodes{
				Address:       branch.Address,
				BankName:      branch.BankName,
				CountryISO2:   branch.CountryISO2,
				IsHeadquarter: branch.IsHeadquarter,
				SwiftCode:     branch.SwiftCode,
			})
		}
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"error": "No records found"})
			return
		}
		c.JSON(http.StatusOK, response)
	}
}

func PostSwift(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newSwiftData model.SwiftData
		if err := c.ShouldBindJSON(&newSwiftData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}
		_, err := db.GetCollection().InsertOne(context.TODO(), bson.M{
			"countryIso2":   newSwiftData.CountryISO2,
			"swiftCode":     newSwiftData.SwiftCode,
			"codeType":      "BIC11",
			"bankName":      newSwiftData.BankName,
			"address":       newSwiftData.Address,
			"townName":      "TEST",
			"countryName":   newSwiftData.CountryName,
			"timeZone":      "Europe/Warsaw",
			"isHeadquarter": newSwiftData.IsHeadquarter,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert data into MongoDB"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "SWIFT Code created successfully"})
	}
}

func DeleteSwift(db database.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		swiftCode := c.Param("swift-code")
		result, err := db.GetCollection().DeleteOne(context.TODO(), bson.M{"swiftCode": swiftCode})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete SWIFT code"})
			return
		}
		if result.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "SWIFT code not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "SWIFT code deleted successfully"})
	}
}
