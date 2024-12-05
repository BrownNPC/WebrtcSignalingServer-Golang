package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("GET /rooms", handleGetAllRooms)
	http.HandleFunc("POST /create", handleRoomsCreate)
	http.HandleFunc("PATCH /offer", handleUpdateOffers)
	http.HandleFunc("PATCH /answer", handleUpdateAnswers)
	log.Println("server is up on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
