package main

import (
	"fmt"
	"net/http"

	"stucon.ramanalabs.in/internal"
)

func main(){
	s, err := internal.InitConn()
	if err != nil {
		fmt.Println("InitConn Failed", err)
		return
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/user", s.LoginHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Allowed", http.StatusForbidden)
	})

	port := "8080"
	fmt.Println("Listening ", port)
	http.ListenAndServe(":"+port, mux)
}
