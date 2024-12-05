package main

import (
	"errors"
	"fmt"
)

type ClientProfile struct {
	// offer ice candidates
	Offers []string `json:"offers"`
	// answers ice candidates
	Answers      []string `json:"answers"`
	PlayerNumber int      `json:"playerNumber"`
	unavailable  bool
}

type MinimalRoomProfile struct {
	Players int
}
type RoomProfile struct {
	MinimalRoomProfile
	Name        string
	Password    string
	PlayerPeers [4]ClientProfile
	Private     bool
	// secret is a password only the room creator knows
	// can be useful for changing the room to private, or public
	secret string
}

var database = map[string]*RoomProfile{} // Use pointers to RoomProfile
// returns an error if room already exists
func DatabaseCreateNewRoom(RoomName string, RoomPassword string, Private bool, Secret string) error {
	_, ok := database[RoomName]
	if ok {
		return errors.New("this room already exists")
	}

	var newRoom RoomProfile
	newRoom.Password = RoomPassword
	newRoom.Private = Private
	newRoom.secret = Secret
	newRoom.Name = RoomName
	for i := range newRoom.PlayerPeers {
		newRoom.PlayerPeers[i].PlayerNumber = i
	}
	database[RoomName] = &newRoom

	return nil
}
func DatabaseGetAllRooms() map[string]MinimalRoomProfile {
	response := make(map[string]MinimalRoomProfile, len(database))
	for roomName, roomProfile := range database {
		if roomProfile.Private {
			continue
		}
		room := MinimalRoomProfile{
			Players: roomProfile.Players,
		}
		response[roomName] = room
	}
	return response
}
func DatabaseGetAvailablePeerFromRoom(roomName string, secret string) (int, ClientProfile, error) {
	room, ok := database[roomName]
	if !ok {
		return -1, ClientProfile{}, errors.New("this room does not exist")
	}

	if room.secret != secret {
		return -1, ClientProfile{}, errors.New("invalid secret was provided")
	}

	// Find the first available peer and return the index and profile

	for i := range room.PlayerPeers {
		if !room.PlayerPeers[i].unavailable {
			return i, room.PlayerPeers[i], nil
		}
	}

	return -1, ClientProfile{}, errors.New("no available peer")
}
func DatabaseUpdateIceCandidateOffers(roomName string, secret string, profile ClientProfile) error {
	// Use DatabaseGetAvailablePeerFromRoom to find an available peer
	index, availablePeer, err := DatabaseGetAvailablePeerFromRoom(roomName, secret)
	if err != nil {
		return err
	}

	// Update the offers of the available peer
	availablePeer.Offers = append(availablePeer.Offers, profile.Offers...)

	// Save the updated peer back into the room
	database[roomName].PlayerPeers[index].Offers = availablePeer.Offers // Dereference the pointer to update the array
	fmt.Println(database[roomName].PlayerPeers[index].Offers)
	return nil
}
