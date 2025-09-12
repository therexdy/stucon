package internal

import (
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
