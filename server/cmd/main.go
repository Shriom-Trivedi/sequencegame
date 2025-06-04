package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"sequencegame/room"
)

func main() {
	router := mux.NewRouter()

	// Register the room WebSocket endpoint.
	// When a client connects via, for example, ws://URL/room/ABC123,
	// they will join (or create) room with ID "ABC123".

	router.HandleFunc("/room/{roomId}", room.RoomWsHandler())

	fmt.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
