package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
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
	if materialId == "" {
		http.Error(w, "Missing material_id Parameter", http.StatusBadRequest)
		return
	}

	queryStr := "SELECT file_path FROM materials WHERE material_id = $1"
	result, err := psqlDB.Query(queryStr, materialId)
	if err != nil {
		http.Error(w, "File Path Query Error", http.StatusInternalServerError)
		return
	}
	defer result.Close()
	var file_path string
	result.Scan(&file_path)

	obj, err := s.Minio.GetObject(s.Ctx, "main", file_path, minio.GetObjectOptions{})
	if err != nil {
			http.Error(w, "object not found", http.StatusNotFound)
			return
	}
	defer obj.Close()

	info, err := obj.Stat()
	if err==nil {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size))
		w.Header().Set("Content-Type", info.ContentType)
	}

	_, err = io.Copy(w, obj)
	if err != nil {
		http.Error(w, "IO error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) FileMetaHandler(w http.ResponseWriter, r *http.Request){

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

	queryStr := `
	SELECT material_id, uploaded_by, scheme_id, branch_id, subject_id, sem, title, file_type, uploaded_at
	FROM materials_with_user
	WHERE subject_id = $1;
	`
	result, err := psqlDB.Query(queryStr, materialId)
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

	if len(rowsArray) != 1 {
		http.Error(w, "Num of rows != 1", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(rowsArray[0])
	if err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
		return
	}
}
