package main

import (
	"MessenFlow/internal/db"
	"MessenFlow/internal/handlers"
	"MessenFlow/internal/websocket"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

func main() {
	pflag.String("address", "localhost", "Server address")
	pflag.String("port", "8080", "Server port")
	pflag.Parse()

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		log.Fatal(err)
	}

	viper.SetEnvPrefix("MESSENFLOW")
	viper.AutomaticEnv()

	serverAddress := viper.GetString("address")
	serverPort := viper.GetString("port")

	err = db.InitDB()
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
	http.HandleFunc("/delete_message", handlers.DeleteMessageHandler)
	http.HandleFunc("/update_message", handlers.UpdateMessageHandler)
	http.HandleFunc("/get_socket", handlers.GetSocketHandler)

	fmt.Println("Server started on:", serverAddress+":"+serverPort)
	log.Fatal(http.ListenAndServe(serverAddress+":"+serverPort, nil))
}
