package websocket

import (
	"SimpleMessenger/internal/db"
	"SimpleMessenger/internal/models"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func Handler(w http.ResponseWriter, r *http.Request) {
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
		_, err = db.DB.Exec("INSERT INTO messages (chatID, user, message) VALUES(?, ?, ?)", chatID, user, string(p))
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

func broadcastMessage(chatID string, user string, message string) {
	msg := models.Message{
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
