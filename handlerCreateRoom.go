package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

// Request structure for creating a room
type CreateRoomRequest struct {
	Name      string `json:"name"`               // Name of the room (must be unique)
	Password  string `json:"password,omitempty"` // Optional password
	IsPrivate bool   `json:"isPrivate"`          // Whether the room should be private

}

// Response structure for room creation
type CreateRoomResponse struct {
	RoomID     string `json:"roomId"`     // The unique room ID
	HostSecret string `json:"hostSecret"` // Secret key for host authentication
}

// HTTP handler for creating a room
func (rs *RoomStore) CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body
	var req CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the input
	if req.Name == "" {
		http.Error(w, "Room name is required", http.StatusBadRequest)
		return
	}

	// Check if the room name is already taken
	if rs.IsRoomNameTaken(req.Name) {
		http.Error(w, "Room name already exists", http.StatusConflict)
		return
	}

	// Generate a host secret
	hostSecret := uuid.New().String()

	// Create the room in the store
	_, err := rs.CreateRoom(req.Name, hostSecret, req.Password, req.IsPrivate)
	if err != nil {
		http.Error(w, "Failed to create room: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare the response
	response := CreateRoomResponse{
		RoomID:     req.Name,
		HostSecret: hostSecret,
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
