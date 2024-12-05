package main

import (
	"errors"
	"fmt"
)

// ClientProfile represents the profile of a client (player) in a room.
// It includes their ICE candidates for offer and answer, their player number,
// and a flag to indicate whether the client is unavailable.
type ClientProfile struct {
	// Offers represents the list of ICE candidates offered by the client.
	Offers []string `json:"offers"`
	// Answers represents the list of ICE candidates the client has answered with.
	Answers []string `json:"answers"`
	// PlayerNumber uniquely identifies the client in a room.
	PlayerNumber int `json:"playerNumber"`
	// Unavailable indicates whether the client is currently unavailable (used for peer management).
	unavailable bool
}

// MinimalRoomProfile contains the basic information of a room, including the number of players
// and whether the room is password-protected or not.
type MinimalRoomProfile struct {
	// Players is the number of players currently in the room.
	Players int
	// PasswordProtected indicates whether the room requires a password to join.
	PasswordProtected bool
}

// RoomProfile represents a full room profile with detailed information about the room.
// It includes the list of player peers, room name, password, and whether the room is private.
type RoomProfile struct {
	// MinimalRoomProfile is embedded to include basic room details like players and password protection.
	MinimalRoomProfile
	// Name is the name of the room.
	Name string
	// Password is the password required to join the room.
	Password string
	// PlayerPeers is an array of 4 client profiles representing the players in the room.
	PlayerPeers [4]ClientProfile
	// Private indicates whether the room is private (requires a secret to join).
	Private bool
	// secret is a password used by the room creator to manage the room.
	secret string
}

// Database is a map where the key is the room name, and the value is a pointer to the RoomProfile.
// This represents all rooms in the system.
var database = map[string]*RoomProfile{} // Use pointers to RoomProfile to allow in-place updates.

// DatabaseCreateNewRoom creates a new room in the system with the provided name, password, privacy status, and secret.
// It returns an error if the room already exists.
func DatabaseCreateNewRoom(RoomName string, RoomPassword string, Private bool, Secret string) error {
	// Check if the room already exists in the database.
	_, ok := database[RoomName]
	if ok {
		return errors.New("this room already exists") // Return an error if the room is already in the database.
	}

	// Create a new RoomProfile and populate its fields.
	var newRoom RoomProfile
	newRoom.Password = RoomPassword
	newRoom.Private = Private
	newRoom.secret = Secret
	newRoom.Name = RoomName
	if RoomPassword != "" {
		// If the room has a password, set the PasswordProtected flag to true.
		newRoom.PasswordProtected = true
	}

	// Add the new room to the database.
	database[RoomName] = &newRoom
	return nil
}

// DatabaseGetAllRooms returns a map of all public rooms (non-private) with minimal room details (players count).
// This map can be used to display available rooms.
func DatabaseGetAllRooms() map[string]MinimalRoomProfile {
	// Create a response map to hold the minimal details for each room.
	response := make(map[string]MinimalRoomProfile, len(database))

	// Loop through all rooms in the database.
	for roomName, roomProfile := range database {
		// Skip private rooms since we only want public rooms in the response.
		if roomProfile.Private {
			continue
		}

		// Add the room's minimal details (like number of players) to the response map.
		room := MinimalRoomProfile{
			Players: roomProfile.Players,
		}
		response[roomName] = room
	}

	return response
}

// DatabaseGetAvailablePeerFromRoom checks for an available peer in the specified room.
// If a peer is available, it returns the index of the peer and the client profile, else an error is returned.
func DatabaseGetAvailablePeerFromRoom(roomName string, secret string) (int, ClientProfile, error) {
	// Retrieve the room by name from the database.
	room, ok := database[roomName]
	if !ok {
		return -1, ClientProfile{}, errors.New("this room does not exist") // Return an error if the room is not found.
	}

	// Check if the secret provided matches the room's secret.
	if room.secret != secret {
		return -1, ClientProfile{}, errors.New("invalid secret was provided") // Return an error if the secret doesn't match.
	}

	// Find the first available peer in the room. We assume that availability is determined by the 'unavailable' flag.
	for i := range room.PlayerPeers {
		if !room.PlayerPeers[i].unavailable {
			// Return the index of the available peer and their profile.
			return i, room.PlayerPeers[i], nil
		}
	}

	// Return an error if no available peer is found.
	return -1, ClientProfile{}, errors.New("no available peer")
}

// DatabaseUpdateIceCandidateOffers updates the ICE candidate offers of an available peer in a room.
// It first finds an available peer and appends the new offers to the peer's existing offers.
func DatabaseUpdateIceCandidateOffers(roomName string, secret string, profile ClientProfile) error {
	// Use DatabaseGetAvailablePeerFromRoom to find an available peer in the specified room.
	index, availablePeer, err := DatabaseGetAvailablePeerFromRoom(roomName, secret)
	if err != nil {
		return err // Return the error if no available peer is found or if there's another issue.
	}

	// Append the new ICE offers to the peer's existing offers.
	availablePeer.Offers = append(availablePeer.Offers, profile.Offers...)

	// Save the updated offers back into the room's PlayerPeers array.
	// Dereference the pointer to update the peer's profile in the room.
	database[roomName].PlayerPeers[index].Offers = availablePeer.Offers
	// For debugging, print the updated offers to the console.
	fmt.Println(database[roomName].PlayerPeers[index].Offers)
	return nil
}

func DatabaseUpdateIceCandidateAnswers(roomName string, password string, profile ClientProfile) error {
	return nil
}
