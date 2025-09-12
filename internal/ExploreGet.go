package internal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

type SchemeResposeJSON struct{
	StrArr []string `json:"strArr"`
}

type Branch struct {
	BranchID string `json:"branchId"`
	BranchName string `json:"branchName"`
}

type BranchResposeJSON struct{
	BranchArr []Branch `json:"branchArr"`
}

type Subject struct {
	SubjectID string `json:"subjectId"`
	SubjectName string `json:"subjectName"`
}

type SubjectResposeJSON struct{
	SubjectArr []Subject `json:"subjectArr"`
}

func (s *Server) SchemeGetHandler(w http.ResponseWriter, r *http.Request){

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

	queryStr := "SELECT scheme_id FROM schemes"
	result, err := s.PSQLDB.Query(queryStr)
	if err != nil {
		http.Error(w, "Query Error", http.StatusInternalServerError)
		return
	}
	defer result.Close()

	var respJSON SchemeResposeJSON
	for result.Next() {
		var temp string
		result.Scan(&temp)
		respJSON.StrArr = append(respJSON.StrArr, temp)
	}

	json.NewEncoder(w).Encode(respJSON)

}

func (s *Server) BranchGetHandler(w http.ResponseWriter, r *http.Request){

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

	queryStr := "SELECT (branch_id, branch_name) FROM branches"
	result, err := s.PSQLDB.Query(queryStr)
	if err != nil {
		http.Error(w, "Query Error", http.StatusInternalServerError)
		return
	}
	defer result.Close()

	var respJSON BranchResposeJSON
	for result.Next() {
		var id, name string
		var branch Branch
		result.Scan(&id, &name)
		branch.BranchID, branch.BranchName = id, name 
		respJSON.BranchArr = append(respJSON.BranchArr, branch)
	}

	json.NewEncoder(w).Encode(respJSON)

}


func (s *Server) SubjectGetHandler(w http.ResponseWriter, r *http.Request){

	if r.Method != http.MethodGet {
		http.Error(w, "Wrong Method", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query()
	querySchemeId := q.Get("scheme_id")
	if querySchemeId == "" {
		http.Error(w, "scheme param can't be empty", http.StatusBadRequest)
		return
	}
	queryBranchId := q.Get("branch_id")
	querySem := q.Get("sem")
	

	psqlDB := s.PSQLDB

	err := psqlDB.Ping()
	if err != nil {
		fmt.Println("PSQL DB Ping Failed")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var result *sql.Rows

	switch {
	case queryBranchId == "":
		queryStr := "SELECT (subject_id, subject_name) FROM branches WHERE (scheme_id = $1 AND branch = $2)"
		result, err = s.PSQLDB.Query(queryStr, querySchemeId, queryBranchId)
		if err != nil {
			http.Error(w, "Query Error", http.StatusInternalServerError)
			return
		}

	case querySem == "":
		queryStr := "SELECT (subject_id, subject_name) FROM branches WHERE (scheme_id = $1 AND branch = $2 AND sem = $3)"
		result, err = s.PSQLDB.Query(queryStr, querySchemeId, queryBranchId, querySem)
		if err != nil {
			http.Error(w, "Query Error", http.StatusInternalServerError)
			return
		}
	}

	var respJSON SubjectResposeJSON
	for result.Next() {
		var id, name string
		var subject Subject
		result.Scan(&id, &name)
		subject.SubjectID, subject.SubjectName = id, name 
		respJSON.SubjectArr = append(respJSON.SubjectArr, subject)
	}

	json.NewEncoder(w).Encode(respJSON)
	result.Close()
}
