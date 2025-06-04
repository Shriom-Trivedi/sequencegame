package game

// Message defines the structure for communications between the client and server.
type Message struct {
    Action   string `json:"action"`
    PlayerID string `json:"player_id,omitempty"`
    Card     string `json:"card,omitempty"`
    X        int    `json:"x,omitempty"`         // For placement moves.
    Y        int    `json:"y,omitempty"`
    TargetX  int    `json:"target_x,omitempty"`  // For one-eyed jack removal moves.
    TargetY  int    `json:"target_y,omitempty"`
}