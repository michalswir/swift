package middleware

import (
	"encoding/csv"
	"io"
	"net/http"
	"strings"
	"swift/internal/model"
	"swift/internal/service"

	"github.com/gin-gonic/gin"
)

func DownloadAsCSV(service *service.GoogleDocsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		fileID := c.Param("fileID")
		csvFile, err := service.DownloadCSV(fileID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to download CSV file"})
			return
		}
		defer csvFile.Close()
		c.Set("csvData", csvFile)
		c.Next()
	}
}

func ParseCSV(c *gin.Context) {
	csvData, exists := c.Get("csvData")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "CSV data not found"})
		return
	}
	reader := csv.NewReader(csvData.(io.Reader))
	var swiftCodes []model.SwiftData
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "CSV parsing error"})
			return
		}
		swiftCodes = append(swiftCodes, model.SwiftData{
			CountryISO2:   strings.ToUpper(record[0]),
			SwiftCode:     record[1],
			CodeType:      record[2],
			BankName:      record[3],
			Address:       record[4],
			TownName:      record[5],
			CountryName:   strings.ToUpper(record[6]),
			TimeZone:      record[7],
			IsHeadquarter: strings.HasSuffix(record[1], "XXX"),
		})
	}
	c.Set("parsedData", swiftCodes)
	c.Next()
}
