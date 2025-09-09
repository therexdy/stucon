package main

import (
	"fmt"
	"net/http"

	"stucon.ramanalabs.in/internal"
)

func main(){
	s, err := internal.InitConn()
	if err != nil {
		fmt.Println("InitConn Failed")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/user", s.LoginHandler)

	http.ListenAndServe(":8080", mux)
}
