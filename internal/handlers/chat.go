package handlers

import (
	"html/template"
	"net/http"
)

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
