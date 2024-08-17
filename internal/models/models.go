package models

type Response struct {
	Users []string `json:"users"`
}

type User struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Message struct {
	Id      int64  `json:"id"`
	User    string `json:"user"`
	Message string `json:"message"`
	Time    string `json:"time"`
	Type    string `json:"type"`
}

type Messages struct {
	Messages []Message `json:"messages"`
}
