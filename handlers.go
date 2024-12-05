package main

import (
	"encoding/json"
	"net/http"
	"regexp"
)

func handleRoomsGet(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetAllRooms(w, r)
	case http.MethodPost:
		handleCreateRoom(w, r)
	case http.MethodPatch:
		handleUpdateOffers(w, r)
	default:
		http.Error(w, "not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetAllRooms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := DatabaseGetAllRooms()
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusExpectationFailed)
	}
}

var isValidName = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString

func handleCreateRoom(w http.ResponseWriter, r *http.Request) {
	RoomName := r.URL.Query().Get("roomId")
	RoomPassword := r.URL.Query().Get("roomPassword")
	Private := r.URL.Query().Has("private")
	Secret := r.URL.Query().Get("secret")
	if !isValidName(RoomName) {
		http.Error(w, "Room name must only contain alphabetic characters (a-z, A-Z)", http.StatusBadRequest)
		return
	}
	if Secret == "" {
		http.Error(w, "Please specify a secret", http.StatusBadRequest)
		return
	}
	if err := DatabaseCreateNewRoom(RoomName, RoomPassword, Private, Secret); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func handleUpdateOffers(w http.ResponseWriter, r *http.Request) {
	payload := ClientProfile{}
	Secret := r.URL.Query().Get("secret")
	RoomName := r.URL.Query().Get("roomId")
	if Secret == "" {
		http.Error(w, "secret was not provided", http.StatusForbidden)
		return
	}
	if RoomName == "" {
		http.Error(w, "roomId was not provided", http.StatusForbidden)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := DatabaseUpdateIceCandidateOffers(RoomName, Secret, payload); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	defer r.Body.Close()
}
