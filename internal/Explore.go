package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

type MaterialJSON struct {
	MaterialId int `json:"materialId"`
	Title string `json:"title"`
	UploadedUser string `json:"uploadedUser"`
	UploadedAt time.Time `json:"uploadedAt"`
	FileType string `json:"fileType"`
	SchemeId string `json:"schemeId"`
	BranchId string `json:"branchId"`
	Sem string `json:"sem"`
	SubjectId string `json:"subjectId"`
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
	limit := q.Get("limit")
	offset := q.Get("offset")
	scheme := q.Get("scheme")
	branch := q.Get("branch")
	sem := q.Get("sem")
	subject := q.Get("subject")

	switch {
	case limit == "" || offset == "":
		http.Error(w, "Limit and Offset Needed", http.StatusBadRequest)
		return
	case scheme == "":
		queryStr := `
		SELECT material_id, uploaded_by, scheme_id, branch_id, subject_id, sem, title, file_type, uploaded_at
		FROM materials_with_user
		ORDER BY uploaded_at DESC
		LIMIT $1 OFFSET $2;
		`
		result, err := psqlDB.Query(queryStr, limit, offset)
		if err != nil {
			http.Error(w, "Query Error", http.StatusInternalServerError)
			return
		}
		defer result.Close()

		var rowsArray []MaterialJSON
		for result.Next(){
			var temp MaterialJSON
			err = result.Scan(&temp.MaterialId, &temp.UploadedUser, &temp.SchemeId, &temp.BranchId, &temp.SubjectId, &temp.Sem, &temp.Title, &temp.FileType, &temp.UploadedAt)
			if err != nil {
				http.Error(w, "Row Scan Error", http.StatusInternalServerError)
				return
			}
			rowsArray = append(rowsArray, temp)
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(rowsArray)
		if err != nil {
			http.Error(w, "JSON encoding error", http.StatusInternalServerError)
			return
		}

	case branch == "":
		queryStr := `
		SELECT material_id, uploaded_by, scheme_id, branch_id, subject_id, sem, title, file_type, uploaded_at
		FROM materials_with_user
		WHERE scheme_id = $1
		ORDER BY uploaded_at DESC
		LIMIT $2 OFFSET $3;
		`
		result, err := psqlDB.Query(queryStr, scheme, limit, offset)
		if err != nil {
			http.Error(w, "Query Error", http.StatusInternalServerError)
			return
		}
		defer result.Close()

		var rowsArray []MaterialJSON
		for result.Next(){
			var temp MaterialJSON
			err = result.Scan(&temp.MaterialId, &temp.UploadedUser, &temp.SchemeId, &temp.BranchId, &temp.SubjectId, &temp.Sem, &temp.Title, &temp.FileType, &temp.UploadedAt)
			if err != nil {
				http.Error(w, "Row Scan Error", http.StatusInternalServerError)
				return
			}
			rowsArray = append(rowsArray, temp)
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(rowsArray)
		if err != nil {
			http.Error(w, "JSON encoding error", http.StatusInternalServerError)
			return
		}
	case sem == "":
		queryStr := `
		SELECT material_id, uploaded_by, scheme_id, branch_id, subject_id, sem, title, file_type, uploaded_at
		FROM materials_with_user
		WHERE (scheme_id = $1 AND branch_id = $2)
		ORDER BY uploaded_at DESC
		LIMIT $3 OFFSET $4;
		`
		result, err := psqlDB.Query(queryStr, scheme, branch, limit, offset)
		if err != nil {
			http.Error(w, "Query Error", http.StatusInternalServerError)
			return
		}
		defer result.Close()

		var rowsArray []MaterialJSON
		for result.Next(){
			var temp MaterialJSON
			err = result.Scan(&temp.MaterialId, &temp.UploadedUser, &temp.SchemeId, &temp.BranchId, &temp.SubjectId, &temp.Sem, &temp.Title, &temp.FileType, &temp.UploadedAt)
			if err != nil {
				http.Error(w, "Row Scan Error", http.StatusInternalServerError)
				return
			}
			rowsArray = append(rowsArray, temp)
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(rowsArray)
		if err != nil {
			http.Error(w, "JSON encoding error", http.StatusInternalServerError)
			return
		}
	case subject == "":
		queryStr := `
		SELECT material_id, uploaded_by, scheme_id, branch_id, subject_id, sem, title, file_type, uploaded_at
		FROM materials_with_user
		WHERE (scheme_id = $1 AND branch_id = $2 AND sem = $3)
		ORDER BY uploaded_at DESC
		LIMIT $4 OFFSET $5;
		`
		result, err := psqlDB.Query(queryStr, scheme, branch, sem, limit, offset)
		if err != nil {
			http.Error(w, "Query Error", http.StatusInternalServerError)
			return
		}
		defer result.Close()

		var rowsArray []MaterialJSON
		for result.Next(){
			var temp MaterialJSON
			err = result.Scan(&temp.MaterialId, &temp.UploadedUser, &temp.SchemeId, &temp.BranchId, &temp.SubjectId, &temp.Sem, &temp.Title, &temp.FileType, &temp.UploadedAt)
			if err != nil {
				http.Error(w, "Row Scan Error", http.StatusInternalServerError)
				return
			}
			rowsArray = append(rowsArray, temp)
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(rowsArray)
		if err != nil {
			http.Error(w, "JSON encoding error", http.StatusInternalServerError)
			return
		}
	default:
		queryStr := `
		SELECT material_id, uploaded_by, scheme_id, branch_id, subject_id, sem, title, file_type, uploaded_at
		FROM materials_with_user
		WHERE (scheme_id = $1 AND branch_id = $2 AND sem = $3 AND subject = $4)
		ORDER BY uploaded_at DESC
		LIMIT $5 OFFSET $6;
		`
		result, err := psqlDB.Query(queryStr, scheme, branch, sem, subject, limit, offset)
		if err != nil {
			http.Error(w, "Query Error", http.StatusInternalServerError)
			return
		}
		defer result.Close()

		var rowsArray []MaterialJSON
		for result.Next(){
			var temp MaterialJSON
			err = result.Scan(&temp.MaterialId, &temp.UploadedUser, &temp.SchemeId, &temp.BranchId, &temp.SubjectId, &temp.Sem, &temp.Title, &temp.FileType, &temp.UploadedAt)
			if err != nil {
				http.Error(w, "Row Scan Error", http.StatusInternalServerError)
				return
			}
			rowsArray = append(rowsArray, temp)
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(rowsArray)
		if err != nil {
			http.Error(w, "JSON encoding error", http.StatusInternalServerError)
			return
		}
	}
}
