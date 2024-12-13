package main

import (
	"encoding/json"
	"net/http"
)

type GetAnswersRequest struct {
	RoomName   string `json:"roomName"`
	PeerId     string `json:"peerId"`
	PeerSecret string `json:"peerSecret"`
}

type GetAnswersResponse struct {
	AnswerSDP           string   `json:"answerSdp"`
	AnswerIceCandidates []string `json:"answerIceCandidates"`
}

func (rs *RoomStore) GetAnswersHandler(w http.ResponseWriter, r *http.Request) {

	var req GetAnswersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.PeerId == "" || req.PeerSecret == "" || req.RoomName == "" {
		http.Error(w, "roomName, peerId and peerSecret are required", http.StatusBadRequest)
		return
	}

	room, err := rs.GetRoom(req.RoomName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return

	}

	rs.mutex.RLock()
	defer rs.mutex.RUnlock()
	peer, ok := room.Peers[req.PeerId]
	if !ok {
		http.Error(w, "This peer is not in this room", http.StatusForbidden)
		return
	}
	if peer.AnswerSDP == "" || len(peer.AnswerICECandidates) == 0 {
		http.Error(w, "no answer yet", http.StatusTooEarly)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	respBody := GetAnswersResponse{
		AnswerSDP:           peer.AnswerSDP,
		AnswerIceCandidates: peer.AnswerICECandidates,
	}
	json.NewEncoder(w).Encode(respBody)
}
