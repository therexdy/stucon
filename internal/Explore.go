package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type MaterialJSON struct {
	MaterialId int `json:"materialId"`
	SubjectId string `json:"subjectId"`
	Title string `json:"title"`
	UploadedUserId string `json:"uploadedUserId"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (s *Server) ExploreHandler(w http.ResponseWriter, r *http.Request){

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
	scheme := q.Get("scheme")
	branch := q.Get("branch")
	sem := q.Get("sem")
	subject := q.Get("subject")

	switch {
	case scheme == "":
		queryStr := ""
		resut
	}

}
