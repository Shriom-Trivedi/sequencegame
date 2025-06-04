package game

import "log"

var Colors = []string{"Red", "Blue", "Green", "Yellow", "Purple", "Orange"} // Colors for players. Works as avatar

// Assign colors to players
func (g *Game) AssignColor() string {
	idx := len(g.Players) % len(Colors)
	return Colors[idx]
}

// ProcessMove handles a player's move based on the card they played.
// - One-eyed Jacks (J♥, J♠) remove an opponent’s chip.
// - Other moves (including two-eyed Jacks, J♦, J♣) place a chip.
func (g *Game) ProcessMove(p *Player, m Message) {
	if m.Card == "J♥" || m.Card == "J♠" {
		if g.ValidRemoval(m.TargetX, m.TargetY, p) {
			g.Board[m.TargetX][m.TargetY] = ""
			log.Printf("Player %s removed chip at (%d, %d)", p.ID, m.TargetX, m.TargetY)
		} else {
			log.Printf("Invalid removal attempt by %s at (%d, %d)", p.ID, m.TargetX, m.TargetY)
		}
	} else {
		// For a normal move or two-eyed jacks, place a chip if the cell is empty or wild.
		if g.Board[m.X][m.Y] == "" || g.Board[m.X][m.Y] == "WILD" {
			g.Board[m.X][m.Y] = p.Color
			log.Printf("Player %s placed chip at (%d, %d)", p.ID, m.X, m.Y)
		} else {
			log.Printf("Invalid placement by %s at (%d, %d) - cell occupied", p.ID, m.X, m.Y)
		}
	}

	// (Manage player's hand here: remove the played card, then draw a new one.)

	// Check win condition.
	if g.CheckWin(p.Color) {
		log.Printf("Player %s with color %s wins!", p.ID, p.Color)
		// In a complete game, you would notify all players about the win.
	}
}

// ValidRemoval checks if a one-eyed jack removal move is valid.
// Removal is only allowed on an opponent's chip (and not on wild corners).
func (g *Game) ValidRemoval(x, y int, p *Player) bool {
	// Do not allow removal from wild corners.
	if (x == 0 && y == 0) || (x == 0 && y == 9) || (x == 9 && y == 0) || (x == 9 && y == 9) {
		return false
	}
	// Only remove an opponent's chip.
	if g.Board[x][y] != "" && g.Board[x][y] != p.Color {
		return true
	}
	return false
}

// CheckWin scans the board for a sequence (five in a row) of chips (or wild cells) belonging to the same color.
func (g *Game) CheckWin(color string) bool {
	n := 10
	sequenceToWin := 5

	// checkDir checks consecutively in a given direction.
	checkDir := func(x, y, dx, dy int) bool {
		for i := 0; i < sequenceToWin; i++ {
			nx := x + i*dx
			ny := y + i*dy
			if nx < 0 || nx >= n || ny < 0 || ny >= n {
				return false
			}
			cell := g.Board[nx][ny]
			if cell != color && cell != "WILD" {
				return false
			}
		}
		return true
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if g.Board[i][j] == color || g.Board[i][j] == "WILD" {
				if checkDir(i, j, 0, 1) || // Horizontal.
					checkDir(i, j, 1, 0) || // Vertical.
					checkDir(i, j, 1, 1) || // Diagonal down-right.
					checkDir(i, j, 1, -1) { // Diagonal down-left.
					return true
				}
			}
		}
	}
	return false
}
