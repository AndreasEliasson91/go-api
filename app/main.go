package main

import (
	"go-api/api"
	"go-api/db"
	"time"

	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func CleanupOldScores() error {
	cutoffTime := time.Now().AddDate(0, -1, 0)

	_, err := db.DB.Exec("DELETE FROM Score WHERE CreatedAt < ?", cutoffTime)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbConnectionString := os.Getenv("DB_CONNECTION_STRING")
	db.InitDB(dbConnectionString)

	err = CleanupOldScores()
	if err != nil {
		log.Fatal("Error cleaning up old scores:", err)
	}

	router := gin.Default()

	router.POST("/score", api.CreateOrUpdateScore)
	router.GET("/score/:id", api.GetScore)
	router.GET("/scores", func(c *gin.Context) {
		scores, err := api.GetAllScores(c)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, scores)
	})
	router.DELETE("/score/:id", api.DeleteScore)
	router.Run(":8080")
}
