// https://github.com/riscie/websocket-tic-tac-toe/blob/350af6181fed/game.go

package main

type status int

const (
	pairing status = iota
	started
	end
)

type game struct {
	// Exported to JSON
	Move       int
	PlayerTurn int
	Status     status

	// Not exported to JSON
	noOfPlayers int
}

func newGame() *game {
	g := &game{
		Move:        -1,
		Status:      pairing,
		noOfPlayers: 0,
		PlayerTurn:  0,
	}
	return g
}

// Restart game between same players
func (g *game) restartGame() *game {
	g.Status = pairing
	g.noOfPlayers += 1
	g.Move = -1
	return g
}

// In case a player disconnects, then wait for 30 secs in client and do a reset
func resetGame() {

}

// Add player to a game
func addPlayer() {

}

// Whenever a player makes a move, switch turn
func (g *game) switchTurn() *game {
	if g.PlayerTurn == 0 {
		g.PlayerTurn = 1
	} else {
		g.PlayerTurn = 0
	}
	return g
}

func (g *game) makeMove(move int) *game {
	g.Move = move
	return g
}

func (g *game) end() *game {
	g.Status = end
	return g
}

func toJSON() {

}
