package main

import (
	"encoding/json"
	"net/http"
)

// Response structure for retrieving all rooms
type GetAllRoomsResponse struct {
	Rooms []RoomSummary `json:"rooms"` // List of room summaries
}

// HTTP handler for retrieving all rooms
func (rs *RoomStore) GetAllRoomsHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve public room summaries
	rooms := rs.GetAllRooms()

	// Add available slots to the response
	for i := range rooms {
		rooms[i].AvailableSlots = 4 - rooms[i].NumPeers // Max 4 players
	}

	// Prepare the response
	response := GetAllRoomsResponse{
		Rooms: rooms,
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
