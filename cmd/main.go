package main

import (
	"fmt"
	"net/http"

	"stucon.ramanalabs.in/internal"
)

func main(){
	s, err := internal.InitConn()
	defer s.CloseConn()
	if err != nil {
		fmt.Println("InitConn Failed", err)
		return
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/user/login", s.LogInHandler)
	mux.HandleFunc("/api/user/validate", s.ValidateSession)
	mux.HandleFunc("/api/user/logout", s.LogOutHandler)
	mux.HandleFunc("/api/user/signup", s.SignUpHandler)
	mux.HandleFunc("/api/explore", s.ExploreHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Allowed", http.StatusForbidden)
	})

	port := "8080"
	fmt.Println("Listening ", port)
	http.ListenAndServe(":"+port, mux)
}
