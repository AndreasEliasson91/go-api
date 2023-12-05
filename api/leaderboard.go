package api

import (
	"go-api/db"
	"log"
)

func GetLeaderboard() ([]db.Leaderboard, error) {
	var leaderboard []db.Leaderboard
	rows, err := db.DB.Query("SELECT * FROM Leaderboard ORDER BY position ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var lb db.Leaderboard
		err := rows.Scan(&lb.ID, &lb.ScoreID, &lb.Position)
		if err != nil {
			return nil, err
		}
		leaderboard = append(leaderboard, lb)
	}

	return leaderboard, nil
}

func UpdateLeaderboard(scoreID int64) {
	rows, err := db.DB.Query("SELECT id FROM Score ORDER BY score DESC LIMIT 25")
	if err != nil {
		log.Println("Error getting top 25 scores:", err)
		return
	}
	defer rows.Close()

	_, err = db.DB.Exec("TRUNCATE TABLE Leaderboard")
	if err != nil {
		log.Println("Error clearing leaderboard:", err)
		return
	}

	position := 1
	for rows.Next() {
		var lb db.Leaderboard
		err := rows.Scan(&lb.ScoreID)
		if err != nil {
			log.Println("Error scanning top 25 scores:", err)
			return
		}
		lb.Position = position
		position++
		_, err = db.DB.Exec("INSERT INTO Leaderboard (score_id, position) VALUES (?, ?)", lb.ScoreID, lb.Position)
		if err != nil {
			log.Println("Error inserting leaderboard:", err)
			return
		}
	}

	rows, err = db.DB.Query("SELECT id FROM Score WHERE id = ? ORDER BY score DESC LIMIT 1", scoreID)
	if err != nil {
		log.Println("Error getting new score:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var lb db.Leaderboard
		err := rows.Scan(&lb.ScoreID)
		if err != nil {
			log.Println("Error scanning new score:", err)
			return
		}
		lb.Position = position
		_, err = db.DB.Exec("INSERT INTO Leaderboard (score_id, position) VALUES (?, ?)", lb.ScoreID, lb.Position)
		if err != nil {
			log.Println("Error inserting new score into leaderboard:", err)
			return
		}
	}
}
