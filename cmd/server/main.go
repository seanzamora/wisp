package main

import (
	"log"
	"net/http"
	"github.com/seanzamora/wisp/internal/server"
)

func main() {
	s := server.NewServer()
	go s.Run()

	http.HandleFunc("/ws", s.HandleWebSocket)

	log.Println("Server starting on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
