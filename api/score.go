package api

import (
	"go-api/db"
	"time"

	"github.com/gin-gonic/gin"
)

func CreateScore(c *gin.Context) {
	var input db.Score
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	result, err := db.DB.Exec("INSERT INTO Score (username, score, created_at) VALUES (?, ?, ?)", input.Username, input.Score, time.Now())
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	scoreID, _ := result.LastInsertId()
	UpdateLeaderboard(scoreID)

	c.JSON(200, gin.H{"message": "Score created successfully"})
}

func GetScore(c *gin.Context) {
	id := c.Param("id")

	var score db.Score
	err := db.DB.QueryRow("SELECT * FROM Score WHERE id = ?", id).Scan(&score.ID, &score.Username, &score.Score, &score.CreatedAt)
	if err != nil {
		c.JSON(404, gin.H{"error": "Score not found"})
		return
	}

	c.JSON(200, score)
}

func DeleteScore(c *gin.Context) {
	id := c.Param("id")

	_, err := db.DB.Exec("DELETE FROM Score WHERE id = ?", id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Score deleted successfully"})
}

func GetAllScores(c *gin.Context) ([]db.Score, error) {
	rows, err := db.DB.Query("SELECT * FROM Score")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var scores []db.Score

	for rows.Next() {
		var score db.Score
		err := rows.Scan(&score.ID, &score.Username, &score.Score, &score.CreatedAt)
		if err != nil {
			return nil, err
		}
		scores = append(scores, score)
	}

	return scores, nil
}
