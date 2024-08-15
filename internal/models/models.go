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
	User    string `json:"user"`
	Message string `json:"message"`
}

type Messages struct {
	Messages []Message `json:"messages"`
}
