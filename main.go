package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3" // Импортируем SQLite драйвер
	"html/template"
	"log"
	"net/http"
	"net/url"
)

var db *sql.DB

type User struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Инициализация базы данных
func initDB() error {
	var err error
	db, err = sql.Open("sqlite3", "./users.db")
	if err != nil {
		return err
	}

	createTableSQL := `
  CREATE TABLE IF NOT EXISTS users (
    login TEXT PRIMARY KEY, 
    password TEXT NOT NULL
  );
  `
	createMessagesTableSQL := `
  CREATE TABLE IF NOT EXISTS messages (
    chatID TEXT NOT NULL,
    user TEXT NOT NULL,
    message TEXT NOT NULL
  );
  `
	_, err = db.Exec(createTableSQL)
	_, err = db.Exec(createMessagesTableSQL)
	return err
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM users WHERE login = ?", user.Login).Scan(&count)
		if err != nil {
			log.Fatal(err)
			return
		}
		if count > 0 {
			http.Error(w, "already registered", http.StatusBadRequest)
			return
		}
		fmt.Println(user.Login, user.Password)
		_, err = db.Exec("INSERT INTO users(login, password) VALUES(?, ?)", user.Login, user.Password)
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
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		var count int
		err = db.QueryRow("SELECT count(*) FROM users WHERE login = ?", user.Login).Scan(&count)
		if err != nil {
			log.Fatal(err)
			return
		}
		if count == 0 {
			http.Error(w, "login isn't exist", http.StatusBadRequest)
			return
		}

		var password string
		err = db.QueryRow("SELECT password FROM users WHERE login = ?", user.Login).Scan(&password)
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

type Response struct {
	Users []string `json:"users"`
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		login := r.URL.Query().Get("user_login")
		fmt.Println(login)

		rows, err := db.Query("SELECT login FROM users WHERE login <> ?", login)

		if err != nil {
			log.Fatal(err)
		}

		var users []string
		defer rows.Close()
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
		respose := Response{
			Users: users,
		}
		json.NewEncoder(w).Encode(respose)
	}
}

func ChatPageHandler(w http.ResponseWriter, r *http.Request) {
	chatID := r.URL.Query().Get("chatID")
	if chatID == "" {
		http.Error(w, "chatID missing", http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("static/chat.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	data := struct {
		ChatID string
	}{
		ChatID: chatID,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func handler(w http.ResponseWriter, r *http.Request) {
	chatID := r.URL.Query().Get("chatID")
	if chatID == "" {
		http.Error(w, "chatID missing", http.StatusBadRequest)
		return
	}
	userLogin, err := r.Cookie("user_login")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := url.QueryUnescape(userLogin.Value)
	fmt.Println(user)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade to WebSocket:", err)
		return
	}
	fmt.Println(chatConnections)
	defer conn.Close()

	handleConnection(conn, chatID, user)
}

var chatConnections = make(map[string][]*websocket.Conn)

func handleConnection(conn *websocket.Conn, chatID string, user string) {
	chatConnections[chatID] = append(chatConnections[chatID], conn)
	fmt.Println(chatConnections[chatID])
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			removeConnection(conn, chatID)
			break
		}
		log.Printf("Received message in chat %s from user %s: %s", chatID, user, string(p))
		_, err = db.Exec("INSERT INTO messages (chatID, user, message) VALUES(?, ?, ?)", chatID, user, string(p))
		if err != nil {
			log.Fatal(err)
		}
		broadcastMessage(chatID, user, string(p))
	}
}

func removeConnection(conn *websocket.Conn, chatID string) {
	connections := chatConnections[chatID]
	for i, c := range connections {
		if c == conn {
			chatConnections[chatID] = append(connections[:i], connections[i+1:]...)
			break
		}
	}
}

type Message struct {
	User    string `json:"user"`
	Message string `json:"message"`
}
type Messages struct {
	Messages []Message `json:"messages"`
}

func broadcastMessage(chatID string, user string, message string) {
	msg := Message{
		User:    user,
		Message: message,
	}
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		log.Println("JSON marshal error:", err)
		return
	}

	for _, conn := range chatConnections[chatID] {
		err := conn.WriteMessage(websocket.TextMessage, jsonMsg)
		if err != nil {
			log.Println("Write error:", err)
			removeConnection(conn, chatID)
		}
	}
}

func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		chatID := r.URL.Query().Get("chatID")
		rows, err := db.Query("SELECT user, message FROM messages WHERE chatID = ?", chatID)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var messages []Message
		for rows.Next() {
			var msg Message
			if err := rows.Scan(&msg.User, &msg.Message); err != nil {
				log.Fatal(err)
			}
			messages = append(messages, msg)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}
		response := Messages{
			Messages: messages,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	err := initDB()
	if err != nil {
		log.Fatal(err)
		return
	}
	http.HandleFunc("/chat", ChatPageHandler)
	http.HandleFunc("/ws", handler)
	// Обработчик для корневого пути, чтобы обслуживать index.html
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/login.html")
	})

	http.HandleFunc("/welcome", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/main.html")
	})

	http.HandleFunc("/register", RegisterHandler)
	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/get_users", GetUsersHandler)
	http.HandleFunc("/get_messages", GetMessagesHandler)
	err = http.ListenAndServe("192.168.1.14:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
