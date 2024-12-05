package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/rooms", handleRoomsGet)
	log.Println("server is up on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
