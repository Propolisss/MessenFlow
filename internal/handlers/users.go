package handlers

import (
	"MessenFlow/internal/db"
	"MessenFlow/internal/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		login := r.URL.Query().Get("user_login")
		fmt.Println(login)

		rows, err := db.DB.Query("SELECT login FROM users WHERE login <> ?", login)

		if err != nil {
			log.Fatal(err)
		}

		var users []string
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(rows)
		for rows.Next() {
			var login string
			err = rows.Scan(&login)
			if err != nil {
				log.Println(err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			users = append(users, login)
		}
		w.Header().Set("Content-Type", "application/json")
		response := models.Response{
			Users: users,
		}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
