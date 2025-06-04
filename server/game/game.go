package game

import "sync"

type BoardType [][]string

type Game struct {
	Board       BoardType
	Players     []*Player
	CurrentTurn int
	Mu          sync.Mutex // Maintaining the game state
}

func NewGame() *Game {
	g := &Game {
		Board: make(BoardType, 10),
		Players: []*Player{},
		CurrentTurn: 0,
	}
	for i:=0; i<10; i++ {
		g.Board[i] = make([]string, 10)
		for j := 0; j < 10; j++ {
			g.Board[i][j] = ""
		}
	}

	// Set WILD corners
	g.Board[0][0] = "WILD"
	g.Board[0][9] = "WILD"
	g.Board[9][0] = "WILD"
	g.Board[9][9] = "WILD"

	return g
}