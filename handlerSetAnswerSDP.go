package main

import (
	"encoding/json"
	"net/http"
)

type SetAnswerRequest struct {
	RoomName            string   `json:"roomName"`
	PeerId              string   `json:"peerId"`
	HostSecret          string   `json:"hostSecret"`
	AnswerSDP           string   `json:"answerSdp"`
	AnswerIceCandidates []string `json:"answerIceCandidates"`
}

func (rs *RoomStore) SetAnswerHandler(w http.ResponseWriter, r *http.Request) {

	var req SetAnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.PeerId == "" || req.HostSecret == "" || req.RoomName == "" || req.AnswerSDP == "" || len(req.AnswerIceCandidates) == 0 {
		http.Error(w, "answerSdp, roomName, peerId, answerIceCandiates, and hostSecret are required", http.StatusBadRequest)
		return
	}

	room, err := rs.GetRoom(req.RoomName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	if room.HostSecret != req.HostSecret {
		http.Error(w, "wrong host secret was sent", http.StatusForbidden)
		return
	}
	rs.mutex.Lock()
	defer rs.mutex.Unlock()
	peer, ok := room.Peers[req.PeerId]
	if !ok {
		http.Error(w, "This peer is not in this room", http.StatusForbidden)
		return
	}

	peer.AnswerSDP = req.AnswerSDP
	peer.AnswerICECandidates = req.AnswerIceCandidates

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
