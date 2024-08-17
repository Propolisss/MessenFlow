package main

import (
	"SimpleMessenger/internal/db"
	"SimpleMessenger/internal/handlers"
	"SimpleMessenger/internal/websocket"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func main() {
	err := db.InitDB()

	if err != nil {
		log.Fatal(err)
		return
	}
	http.HandleFunc("/chat", handlers.ChatPageHandler)
	http.HandleFunc("/ws", websocket.Handler)

	http.Handle("/static/css/", http.StripPrefix("/static/css/", http.FileServer(http.Dir("static/css"))))
	http.Handle("/static/js/", http.StripPrefix("/static/js/", http.FileServer(http.Dir("static/js"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/html/login.html")
	})

	http.HandleFunc("/welcome", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/html/main.html")
	})

	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/get_users", handlers.GetUsersHandler)
	http.HandleFunc("/get_messages", handlers.GetMessagesHandler)
	err = http.ListenAndServe("192.168.1.14:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
