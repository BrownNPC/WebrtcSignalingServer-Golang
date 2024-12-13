package main

import (
	"encoding/json"
	"net/http"
)

// Response structure for retrieving all rooms
type GetAllRoomsScaryResponse struct {
	Rooms []Room `json:"rooms"` // List of room summaries
}

// HTTP handler for retrieving all rooms
func (rs *RoomStore) GetAllRoomsScaryHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve public room summaries
	rooms := rs.GetAllRoomsScary()

	// Prepare the response
	response := GetAllRoomsScaryResponse{
		Rooms: rooms,
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
