package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

type SignUpRequestJSON struct {
	Email string `json:"email"`
	Name string `json:"name"`
	Password string `json:"password"`
}

type SignUpResponseJSON struct {
	Valid bool `json:"valid"`
	Token string `json:"token"`
}


func (s *Server) SignUpHandler(w http.ResponseWriter, r *http.Request){

	if r.Method != http.MethodPost {
		http.Error(w, "Wrong Method", http.StatusMethodNotAllowed)
		return
	}

	psqlDB := s.PSQLDB
	rDB := s.RedisDB

	err := psqlDB.Ping()
	if err != nil {
		fmt.Println("PSQL DB Ping Failed")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var reqJson SignUpRequestJSON
	err = json.NewDecoder(r.Body).Decode(&reqJson)
	if err != nil {
		http.Error(w, "JSON Decoder Error", http.StatusInternalServerError)
		return
	}

	var count int
	countQuery := "SELECT COUNT(*) FROM users WHERE email=($1)"
	err = psqlDB.QueryRow(countQuery, reqJson.Email).Scan(&count)
	if err != nil {
		http.Error(w, "Query Error: " + err.Error(), http.StatusInternalServerError)
		return
	}

	if count != 0 {
		http.Error(w, "Email already exists", http.StatusInternalServerError)
		return
	}

	passwordHash, err := HashPassword(reqJson.Password)
	if err != nil {
		http.Error(w, "pwd Hash Error", http.StatusInternalServerError)
		return
	}

	insertQuery := "INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)"
	result, err := psqlDB.Exec(insertQuery, reqJson.Name, reqJson.Email, passwordHash)

	if rowsAffected, _ := result.RowsAffected(); rowsAffected != 1 {
		http.Error(w, "RowsAffected not equal to 1", http.StatusInternalServerError)
		return
	}

	token, err := generateSessionToken(128)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}
	
	err = rDB.Set(s.Ctx, token, reqJson.Email, 24*time.Hour).Err()
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	var respJson SignUpResponseJSON
	respJson.Valid = true
	respJson.Token = token

	json.NewEncoder(w).Encode(respJson)

}
