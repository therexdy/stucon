package internal

import (
	"encoding/json"
	"fmt"
	"net/http"

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

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request){

	if r.Method != http.MethodPut {
		http.Error(w, "Wrong Method", http.StatusMethodNotAllowed)
		return
	}

	psqlDB := s.PSQLDB

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

	pwdQuery := "SELECT password_hash FROM normal_users WHERE email=(?)"
	result, err := psqlDB.Query(pwdQuery, )
	defer result.Close()

	var pwdHashFromDb string
	err = result.Scan(&pwdHashFromDb)
	if err != nil {
		http.Error(w, "DB Row Scan Error", http.StatusInternalServerError)
		return
	}

	if pwdHashFromDb != reqJson.PasswordHash {
		http.Error(w, "pwd Hash did not match", http.StatusInternalServerError)
		return
	}

	respJson := ResponseJSON{
		Valid: true,
		Token: "temp_token",
	}

	err = json.NewEncoder(w).Encode(respJson)
	if err != nil {
		http.Error(w, "JSON Encoder Error", http.StatusInternalServerError)
		return
	}

	return
}
