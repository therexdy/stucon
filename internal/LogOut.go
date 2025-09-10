package internal

import (
	"encoding/json"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type LogOutRequestJSON struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func (s *Server) LogOutHandler(w http.ResponseWriter, r *http.Request){

	if r.Method != http.MethodPut {
		http.Error(w, "Wrong Method", http.StatusMethodNotAllowed)
		return
	}

	rDB := s.RedisDB

	var reqJson LogOutRequestJSON
	err := json.NewDecoder(r.Body).Decode(&reqJson)
	if err != nil {
		http.Error(w, "JSON Decoder Error", http.StatusInternalServerError)
		return
	}

	tokenFromRedis, err := rDB.Get(s.Ctx, reqJson.Email).Result()
	if err == redis.Nil {
		http.Error(w, "Bad LogIn State", http.StatusUnauthorized)
		return
	}else if err != nil {
		http.Error(w, "Redis Get Error", http.StatusInternalServerError)
		return
	}else if tokenFromRedis != reqJson.Token {
		http.Error(w, "Token did not Match", http.StatusUnauthorized)
	}

	result, err := rDB.Del(s.Ctx, reqJson.Email).Result()
	if err != nil {
		http.Error(w, "Could not log out", http.StatusInternalServerError)
		return
	}

	if result == 0 {
		http.Error(w, "Key not found", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
