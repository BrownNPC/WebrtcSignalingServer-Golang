package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	roomStore := NewRoomStore()
	mux := http.NewServeMux()
	// GET
	mux.HandleFunc("GET /get-rooms", roomStore.GetAllRoomsHandler)
	mux.HandleFunc("GET /get-rooms-scary", roomStore.GetAllRoomsScaryHandler)
	mux.HandleFunc("POST /get-peers", roomStore.GetPeersHandler)
	// post on purpose cuz we need to pass some shit
	mux.HandleFunc("POST /get-answer", roomStore.GetAnswersHandler)

	// POST
	mux.HandleFunc("POST /create-room", roomStore.CreateRoomHandler)
	mux.HandleFunc("POST /set-answer", roomStore.SetAnswerHandler)

	mux.HandleFunc("POST /add-peer", roomStore.AddPeerHandler)

	// Start the HTTP server
	fmt.Println("starting server on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(mux)))
}
