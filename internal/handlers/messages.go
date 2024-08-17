package handlers

import (
	"SimpleMessenger/internal/db"
	"SimpleMessenger/internal/models"
	"encoding/json"
	"log"
	"net/http"
)

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		chatID := r.URL.Query().Get("chatID")
		rows, err := db.DB.Query("SELECT id, user, message, time FROM messages WHERE chatID = ?", chatID)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var messages []models.Message
		for rows.Next() {
			var msg models.Message
			if err := rows.Scan(&msg.Id, &msg.User, &msg.Message, &msg.Time); err != nil {
				log.Fatal(err)
			}
			messages = append(messages, msg)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}
		response := models.Messages{
			Messages: messages,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Fatal(err)
		}
	}
}
