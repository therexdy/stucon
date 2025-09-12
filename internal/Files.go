package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

func (s *Server) FileHandler(w http.ResponseWriter, r *http.Request){

	if r.Method != http.MethodGet {
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

	q := r.URL.Query()
	materialId := q.Get("material_id")

}
