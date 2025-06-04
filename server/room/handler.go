package room

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/gorilla/websocket"
    "sequencegame/server/game"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

// RoomWsHandler returns an HTTP handler function to handle WebSocket connections for a room.
func RoomWsHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Extract roomId from the URL.
        vars := mux.Vars(r)
        roomID := vars["roomId"]

        // Retrieve or create the room associated with roomID.
        roomInstance := GetRoom(roomID)
        
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            log.Println("WebSocket upgrade error:", err)
            return
        }
        defer conn.Close()

        var currentPlayer *game.Player

        // Read messages from the client.
        for {
            _, msg, err := conn.ReadMessage()
            if err != nil {
                log.Println("Error reading message:", err)
                break
            }

            var m game.Message
            if err := json.Unmarshal(msg, &m); err != nil {
                log.Println("JSON unmarshal error:", err)
                continue
            }

            switch m.Action {
            case "join":
                currentPlayer = &game.Player{
                    ID:    m.PlayerID,
                    Color: roomInstance.Game.AssignColor(),
                    Conn:  conn,
                    // Initial hand is a placeholder.
                    Cards: []string{"7♥", "8♦", "2♣", "J♦", "A♠"},
                }
                roomInstance.Game.Mu.Lock()
                roomInstance.Game.Players = append(roomInstance.Game.Players, currentPlayer)
                roomInstance.Game.Mu.Unlock()
                reply := map[string]string{
                    "message": "joined",
                    "id":      currentPlayer.ID,
                    "color":   currentPlayer.Color,
                }
                replyMsg, _ := json.Marshal(reply)
                conn.WriteMessage(websocket.TextMessage, replyMsg)

            case "move":
                if currentPlayer == nil {
                    log.Println("Player not joined; cannot process move")
                    continue
                }

                roomInstance.Game.Mu.Lock()
                if roomInstance.Game.Players[roomInstance.Game.CurrentTurn].ID != currentPlayer.ID {
                    roomInstance.Game.Mu.Unlock()
                    errMsg := map[string]string{"error": "Not your turn"}
                    replyMsg, _ := json.Marshal(errMsg)
                    conn.WriteMessage(websocket.TextMessage, replyMsg)
                    continue
                }

                roomInstance.Game.ProcessMove(currentPlayer, m)
                // Advance the turn.
                roomInstance.Game.CurrentTurn = (roomInstance.Game.CurrentTurn + 1) % len(roomInstance.Game.Players)
                roomInstance.Game.Mu.Unlock()

                BroadcastGameState(roomInstance.Game)
            default:
                log.Println("Unknown action:", m.Action)
            }
        }
    }
}

// BroadcastGameState sends the updated game state to every connected player in the room.
func BroadcastGameState(g *game.Game) {
    state := struct {
        Board       [][]string `json:"board"`
        CurrentTurn int        `json:"current_turn"`
    }{
        Board:       g.Board,
        CurrentTurn: g.CurrentTurn,
    }
    msg, err := json.Marshal(state)
    if err != nil {
        log.Println("Error marshaling game state:", err)
        return
    }
    for _, p := range g.Players {
        if p.Conn != nil {
            p.Conn.WriteMessage(websocket.TextMessage, msg)
        }
    }
}
