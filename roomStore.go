package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Room represents a game room with metadata and peers.
type Room struct {
	ID         string           // Unique identifier for the room
	HostSecret string           // Secret key for host authentication
	Password   string           // Optional room password
	IsPrivate  bool             // Whether the room is private
	Peers      map[string]*Peer // Connected peers in the room
	CreatedAt  int64            // Room creation timestamp
	LastActive int64            // Timestamp of the last activity
}

// Peer represents a connected player with their ICE candidates.
type Peer struct {
	Secret              string
	OfferICECandidates  []string `json:"offerIceCandidates"`  // ICE candidates from the peer
	AnswerICECandidates []string `json:"answerIceCandidates"` // ICE candidates from the room creator
	OfferSDP            string   `json:"offerSdp"`
	AnswerSDP           string   `json:"answerSdp"`
}

// RoomSummary provides public information about a room.
type RoomSummary struct {
	ID             string `json:"id"`             // Room name
	NumPeers       int    `json:"numPeers"`       // Number of connected peers
	AvailableSlots int    `json:"availableSlots"` // Number of available slots
	IsPrivate      bool   `json:"isPrivate"`      // Whether the room is private
}

// RoomStore manages rooms in memory with thread-safe access.
type RoomStore struct {
	rooms map[string]*Room // Map of room IDs to rooms
	mutex sync.RWMutex     // Mutex for thread-safe access
}

func (rs *RoomStore) CreateRoom(id, hostSecret, password string, isPrivate bool) (*Room, error) {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	if _, exists := rs.rooms[id]; exists {
		return nil, errors.New("room with the given ID already exists")
	}

	room := &Room{
		ID:         id,
		HostSecret: hostSecret,
		Password:   password,
		IsPrivate:  isPrivate,
		Peers:      make(map[string]*Peer),
		CreatedAt:  time.Now().Unix(),
		LastActive: time.Now().Unix(),
	}

	rs.rooms[id] = room
	return room, nil
}
func (rs *RoomStore) GetAllRoomsScary() []Room {

	rs.mutex.RLock()
	defer rs.mutex.RUnlock()
	rooms := []Room{}
	for _, room := range rs.rooms {
		rooms = append(rooms, *room)
	}
	return rooms
}
func (rs *RoomStore) GetAllRooms() []RoomSummary {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	summaries := []RoomSummary{}
	for _, room := range rs.rooms {
		if room.IsPrivate {
			continue // Skip private rooms
		}
		summaries = append(summaries, RoomSummary{
			ID:        room.ID,
			NumPeers:  len(room.Peers),
			IsPrivate: room.IsPrivate,
		})
	}
	return summaries
}

func (rs *RoomStore) RemovePeerFromRoom(roomID, peerID string) error {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	room, exists := rs.rooms[roomID]
	if !exists {
		return errors.New("room not found")
	}

	if _, exists := room.Peers[peerID]; !exists {
		return errors.New("peer not found in the room")
	}

	delete(room.Peers, peerID)
	room.LastActive = time.Now().Unix()
	return nil
}

// GetRoom retrieves a room by its name.
func (rs *RoomStore) GetRoom(name string) (*Room, error) {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	room, exists := rs.rooms[name]
	if !exists {
		return nil, fmt.Errorf("room not found")
	}
	return room, nil
}

func (rs *RoomStore) UpdateRoom(roomID, hostSecret, password string, isPrivate bool) error {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	room, exists := rs.rooms[roomID]
	if !exists {
		return errors.New("room not found")
	}

	if room.HostSecret != hostSecret {
		return errors.New("invalid host secret")
	}

	room.Password = password
	room.IsPrivate = isPrivate
	room.LastActive = time.Now().Unix()
	return nil
}

func (rs *RoomStore) DeleteRoom(roomID, hostSecret string) error {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	room, exists := rs.rooms[roomID]
	if !exists {
		return errors.New("room not found")
	}

	if room.HostSecret != hostSecret {
		return errors.New("invalid host secret")
	}

	delete(rs.rooms, roomID)
	return nil
}
func (rs *RoomStore) CleanupInactiveRooms(timeout int64) {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()

	now := time.Now().Unix()
	for id, room := range rs.rooms {
		if now-room.LastActive > timeout {
			delete(rs.rooms, id)
		}
	}
}

func (rs *RoomStore) IsRoomNameTaken(name string) bool {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()

	for _, room := range rs.rooms {
		if room.ID == name { // Using room name as ID
			return true
		}
	}
	return false
}

func NewRoomStore() *RoomStore {
	return &RoomStore{
		rooms: make(map[string]*Room),
	}
}
