package game

import "github.com/gorilla/websocket"

type Player struct {
	ID string `json:"id"`
	Cards []string `json:"cards"`
	Color string `json:"color"`
	Conn *websocket.Conn
}