// controllers/dbController.go

package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"url-shortener/utils"
)

var collection *mongo.Collection

type URLPayload struct {
	URL string `json:"url"`
}

// InitDB initializes the MongoDB connection
func InitDB(connectionString string, dbName string, collectionName string) {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database(dbName).Collection(collectionName)
}

func AddURL(c *gin.Context) {
	var payload URLPayload

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(400, gin.H{
			"error": "Invalid JSON payload",
		})
		return
	}
	baseURL := payload.URL

	// Ensure the URL starts with "https://"
	if strings.HasPrefix(baseURL, "http://") {
		baseURL = strings.Replace(baseURL, "http://", "https://", 1)
	} else if !strings.HasPrefix(baseURL, "https://") {
		baseURL = "https://" + baseURL
	}

	// Ensure the URL has a valid TLD
	validURLPattern := `^https://(?:www\.)?[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+.*$`
	match, _ := regexp.MatchString(validURLPattern, baseURL)
	if !match {
		c.JSON(400, gin.H{
			"error": "Invalid URL format. Must contain a top-level domain.",
		})
		return
	}

	shortURL := utils.CreateShortURL(5)
	var data = bson.M{
		shortURL: baseURL,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if a record with the given value exists
	pipeline := mongo.Pipeline{
		{{"$match", bson.M{
			"$expr": bson.M{
				"$ne": bson.A{
					bson.M{
						"$filter": bson.M{
							"input": bson.M{"$objectToArray": "$$ROOT"},
							"as":    "field",
							"cond": bson.M{
								"$eq": bson.A{"$$field.v", baseURL},
							},
						},
					},
					bson.A{},
				},
			},
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Error checking for existing record",
		})
		return
	}

	var existingRecords []bson.M
	err = cursor.All(ctx, &existingRecords)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Error decoding existing records",
		})
		return
	}

	// If the record exists, return an error
	if len(existingRecords) > 0 {
		c.JSON(409, gin.H{
			"error": "Record already exists",
		})
		return
	}

	// Insert the record if it doesn't exist
	result, err := collection.InsertOne(ctx, data)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Failed to insert URL",
		})
		return
	}

	c.JSON(200, gin.H{
		"result":   result,
		"message":  "URL inserted successfully",
		"shortUrl": shortURL,
	})
}

func RedirectToURL(c *gin.Context) {
	shortURL := c.Param("shortURL")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{shortURL: bson.M{"$exists": true}}
	var result bson.M
	err := collection.FindOne(ctx, filter).Decode(&result)

	if err == nil && result != nil {
		destinationURL := result[shortURL].(string)
		c.Redirect(http.StatusMovedPermanently, destinationURL)
		return
	}

	// If no match is found, redirect to google.com
	c.Redirect(http.StatusMovedPermanently, "https://mun.day/urlshortener")
}
