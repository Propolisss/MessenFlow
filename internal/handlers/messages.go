package handlers

import (
	"MessenFlow/internal/db"
	"MessenFlow/internal/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		chatID := r.URL.Query().Get("chatID")
		rows, err := db.DB.Query("SELECT id, user, message, time FROM messages WHERE chatID = ?", chatID)
		if err != nil {
			log.Fatal(err)
		}
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(rows)

		var messages []models.Message
		for rows.Next() {
			var msg models.Message
			if err := rows.Scan(&msg.Id, &msg.User, &msg.Message, &msg.Time); err != nil {
				log.Fatal(err)
			}
			msg.Type = "add"
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

func DeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}
	fmt.Println(idInt)
	result, err := db.DB.Exec("DELETE FROM messages WHERE id = ?", idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Message deleted"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func UpdateMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	messageId := r.URL.Query().Get("message_id")
	new_text := r.URL.Query().Get("new_text")
	if messageId == "" || new_text == "" {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}
	idInt, err := strconv.Atoi(messageId)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}
	_, err = db.DB.Exec("UPDATE messages SET message = ? WHERE id = ?", new_text, idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Message updated"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
