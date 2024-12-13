package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Request structure for adding a peer
type AddPeerRequest struct {
	PeerID            string   `json:"peerId"`   // Unique identifier for the peer
	RoomName          string   `json:"roomName"` // Name of the room
	Password          string   `json:"password"` // Password for the room (optional)
	OfferSDP          string   `json:"offerSdp"`
	OfferIceCandiates []string `json:"offerIceCandidates"`
}

// Response structure for adding a peer
type AddPeerResponse struct {
	Secret string `json:"peerSecret"` // Confirmation message
}

// HTTP handler for adding a peer to a room
func (rs *RoomStore) AddPeerHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON request body
	var req AddPeerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.PeerID == "" || req.RoomName == "" || req.OfferSDP == "" || len(req.OfferIceCandiates) == 0 {
		http.Error(w, "Peer ID, SDP, offerIceCandidates[], and room name are required", http.StatusBadRequest)
		return
	}

	// Get the room from the store
	room, err := rs.GetRoom(req.RoomName)
	if err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	// Validate the room's password
	if room.Password != "" && room.Password != req.Password {
		http.Error(w, "Invalid room password", http.StatusUnauthorized)
		return
	}

	// Check if the room has space for another peer
	rs.mutex.Lock()
	defer rs.mutex.Unlock()
	if len(room.Peers) >= 4 {
		http.Error(w, "Room is full", http.StatusConflict)
		return
	}

	// Check if the peer already exists
	if _, exists := room.Peers[req.PeerID]; exists {
		http.Error(w, "Peer already exists in the room", http.StatusConflict)
		return
	}

	// Add the peer to the room
	newPeer := &Peer{
		Secret:              uuid.New().String(),
		AnswerICECandidates: make([]string, 0),
		OfferICECandidates:  req.OfferIceCandiates,
		OfferSDP:            req.OfferSDP,
	}
	room.Peers[req.PeerID] = newPeer
	room.LastActive = time.Now().Unix() // Update last activity time

	// Prepare the response
	response := AddPeerResponse{Secret: newPeer.Secret}
	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
