package usersservice

import (
	"net/http"
	"encoding/json"
)

func usersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		json.NewEncoder(w).Encode(map[string]string{"message": "GET users"})
	case http.MethodPost:
		json.NewEncoder(w).Encode(map[string]string{"message": "POST users"})
	case http.MethodPut:
		json.NewEncoder(w).Encode(map[string]string{"message": "PUT users"})
	case http.MethodDelete:
		json.NewEncoder(w).Encode(map[string]string{"message": "DELETE users"})
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
