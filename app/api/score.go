package api

import (
	"database/sql"
	"go-api/db"
	"time"

	"github.com/gin-gonic/gin"
)

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

func CreateOrUpdateScore(c *gin.Context) {
	var input db.Score
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var existingScore db.Score
	err := db.DB.QueryRow("SELECT * FROM Score WHERE Username = ?", input.Username).Scan(&existingScore.ID, &existingScore.Username, &existingScore.Score, &existingScore.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			_, err := db.DB.Exec("INSERT INTO Score (Username, Score, CreatedAt) VALUES (?, ?, ?)", input.Username, input.Score, time.Now())
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			c.JSON(200, gin.H{"message": "Score created successfully"})
			return
		}

		// Handle other errors
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if input.Score > existingScore.Score {
		_, err := db.DB.Exec("UPDATE Score SET Score = ?, CreatedAt = ? WHERE Username = ?", input.Score, time.Now(), input.Username)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Score updated successfully"})
		return
	}

	c.JSON(200, gin.H{"message": "Score not updated. New score is not better than the existing score."})
}
