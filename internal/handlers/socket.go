package handlers

import (
	"github.com/spf13/viper"
	"net/http"
)

func GetSocketHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		_, err := w.Write([]byte(viper.GetString("address") + ":" + viper.GetString("port") + "/"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
