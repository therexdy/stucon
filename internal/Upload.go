package internal

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
)

type queryParams struct{
	userId int
	schemeId string
	branchId string
	subjectId string
	sem int
	title string
	fileType string
}

func (s *Server) UploadHandler(w http.ResponseWriter, r *http.Request){

	if r.Method != http.MethodPost {
		http.Error(w, "Wrong Method", http.StatusMethodNotAllowed)
		return
	}

	psqlDB := s.PSQLDB

	err := psqlDB.Ping()
	if err != nil {
		fmt.Println("PSQL DB Ping Failed")
		http.Error(w, "Internal Server Error " + err.Error(), http.StatusInternalServerError)
		return
	}

	var params queryParams
	q := r.URL.Query()
	params.userId, err = strconv.Atoi(q.Get("user_id"))
	if err != nil {
		http.Error(w, "user_id param not of valid form "+err.Error(), http.StatusBadRequest)
		return
	}
	params.schemeId = q.Get("scheme_id")
	params.branchId = q.Get("branch_id")
	params.subjectId = q.Get("subject_id")
	params.sem, err = strconv.Atoi(q.Get("sem"))
	if err != nil {
		http.Error(w, "sem param not of valid form"+err.Error(), http.StatusBadRequest)
		return
	}
	params.title = q.Get("title")
	params.fileType = q.Get("file_type")

	switch "" {
	case params.subjectId, params.branchId, params.schemeId, params.title, params.fileType:
		http.Error(w, "Missing Params", http.StatusBadRequest)
		return
	}

	preHash := params.schemeId+params.branchId+string(params.sem)+params.subjectId+string(params.userId)+params.title+params.fileType
	hash := sha256.Sum256([]byte(preHash))

	fileName := hex.EncodeToString(hash[:])+"."+params.fileType
	
	queryStr := `
		INSERT INTO materials
		(uploaded_by_user, scheme_id, branch_id, subject_id, sem, title, file_path, file_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	result , err := psqlDB.Exec(queryStr, params.userId, params.schemeId, params.branchId, params.subjectId, params.sem, params.title, fileName, params.fileType)
	if err != nil {
		http.Error(w, "DB Exec Error " + err.Error(), http.StatusInternalServerError)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected != 1{
		http.Error(w, "DB Return Rows Error " + err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = s.Minio.PutObject(
		s.Ctx,
		"main",
		fileName,
		r.Body,
		r.ContentLength,
		minio.PutObjectOptions{ContentType: "application/octet-stream"},
		)
	if err != nil {
		http.Error(w, "Upload Failed " + err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Upload Successful"))
}

