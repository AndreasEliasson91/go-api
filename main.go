package main

import (
	"go-api/api"
	"go-api/db"

	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbConnectionString := os.Getenv("DB_CONNECTION_STRING")
	db.InitDB(dbConnectionString)

	router := gin.Default()

	router.POST("/score", api.CreateScore)
	router.GET("/score/:id", api.GetScore)
	router.GET("/scores", func(c *gin.Context) {
		scores, err := api.GetAllScores(c)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, scores)
	})
	router.GET("/leaderboard", func(c *gin.Context) {
		leaderboard, err := api.GetLeaderboard()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, leaderboard)
	})
	router.DELETE("/score/:id", api.DeleteScore)

	router.Run(":8080")
}
