package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"url-shortener/controllers"
)

func SetupRoutes(r *gin.Engine) {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	// Initialize the MongoDB connection
	connectionString := os.Getenv("MONGO_CONNECTION_STRING")
	dbName := os.Getenv("MONGO_DB_NAME")
	collectionName := os.Getenv("MONGO_COLLECTION_NAME")
	controllers.InitDB(connectionString, dbName, collectionName)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/shorten", controllers.AddURL)
	r.GET("/:shortURL", controllers.RedirectToURL)
}
