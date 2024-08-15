package handlers

import (
	"SimpleMessenger/internal/db"
	"SimpleMessenger/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		var count int
		err = db.DB.QueryRow("SELECT COUNT(*) FROM users WHERE login = ?", user.Login).Scan(&count)
		if err != nil {
			log.Fatal(err)
			return
		}
		if count > 0 {
			http.Error(w, "already registered", http.StatusBadRequest)
			return
		}
		fmt.Println(user.Login, user.Password)
		_, err = db.DB.Exec("INSERT INTO users(login, password) VALUES(?, ?)", user.Login, user.Password)
		if err != nil {
			http.Error(w, "internal error", http.StatusBadRequest)
			log.Fatal(err)
			return
		}
		w.WriteHeader(http.StatusCreated)
	} else {
		http.Error(w, "wrong method", http.StatusMethodNotAllowed)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		var count int
		err = db.DB.QueryRow("SELECT count(*) FROM users WHERE login = ?", user.Login).Scan(&count)
		if err != nil {
			log.Fatal(err)
			return
		}
		if count == 0 {
			http.Error(w, "login isn't exist", http.StatusBadRequest)
			return
		}

		var password string
		err = db.DB.QueryRow("SELECT password FROM users WHERE login = ?", user.Login).Scan(&password)
		if err != nil {
			log.Fatal(err)
			return
		}
		if password != user.Password {
			http.Error(w, "wrong password", http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "http://192.168.1.14:8080/welcome", 302)
	} else {
		http.Error(w, "wrong method", http.StatusBadRequest)
	}
}
