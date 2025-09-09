package internal

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

type RequestJSON struct {
	Email string `json:"email"`
	PasswordHash string `json:"passwordHash"`
}

type ResponseJSON struct {
	Valid bool `json:"valid"`
	Token string `json:"token"`
}

func generateSessionToken(n int) (token string, err error) {
	b := make([]byte, n)

	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request){

	if r.Method != http.MethodPut {
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

	var reqJson RequestJSON
	err = json.NewDecoder(r.Body).Decode(&reqJson)
	if err != nil {
		http.Error(w, "JSON Decoder Error", http.StatusInternalServerError)
		return
	}

	pwdQuery := "SELECT password_hash FROM normal_users WHERE email=($1)"
	result, err := psqlDB.Query(pwdQuery, reqJson.Email)
	if err != nil {
		http.Error(w, "Query Error: " + err.Error(), http.StatusInternalServerError)
		return
	}
	defer result.Close()

	var pwdHashFromResult []string
	for result.Next() {
		var temp string
		if err = result.Scan(&temp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pwdHashFromResult = append(pwdHashFromResult, temp)
	}

	if len(pwdHashFromResult) != 1 {
		http.Error(w, "Number of Rows is not 1", http.StatusInternalServerError)
		return
	}

	pwdHashFromDb := pwdHashFromResult[0]

	if pwdHashFromDb != reqJson.PasswordHash {
		http.Error(w, "pwd Hash did not match", http.StatusInternalServerError)
		return
	}

	token, err := generateSessionToken(128)
	if err != nil {
		http.Error(w, "Could not Generate Session Token", http.StatusInternalServerError)
		return
	}

	err = rDB.Set(s.Ctx, token, reqJson.Email, time.Hour).Err()
	if err != nil {	
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respJson := ResponseJSON{
		Valid: true,
		Token: token,
	}

	err = json.NewEncoder(w).Encode(respJson)
	if err != nil {
		http.Error(w, "JSON Encoder Error", http.StatusInternalServerError)
		return
	}
}
