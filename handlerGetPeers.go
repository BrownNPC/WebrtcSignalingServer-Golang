package main

import (
	"encoding/json"
	"net/http"
)

type GetPeersRequest struct {
	RoomName   string `json:"roomName"`
	HostSecret string `json:"hostSecret"`
}

func (rs *RoomStore) GetPeersHandler(w http.ResponseWriter, r *http.Request) {

	var req GetPeersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if req.HostSecret == "" || req.RoomName == "" {
		http.Error(w, "hostSecret and roomName are required", http.StatusBadRequest)
		return
	}

	room, err := rs.GetRoom(req.RoomName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	if len(room.Peers) == 0 {
		http.Error(w, "no peers yet", http.StatusTooEarly)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(room.Peers)
}
