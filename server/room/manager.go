package room

import (
	"sync"

	"sequencegame/server/game"
)

type Room struct {
	Game *game.Game
	// TODO: add more room specific data here...
}

var (
	rooms = make(map[string]*Room)
	roomsMu sync.Mutex
)

func GetRoom(roomID string) *Room {
	roomsMu.Lock()
	defer roomsMu.Unlock()

	if room, exists := rooms[roomID]; exists {
		return room
	}

	newRoom := &Room{
		Game: game.NewGame(),
	}

	rooms[roomID] = newRoom
	return newRoom
}