// https://github.com/riscie/websocket-tic-tac-toe/blob/350af6181fed/game.go

package main

type status string

const (
	pairing status = "pairing"
	restart status = "restart"
	started status = "started"
)

type player struct {
	Id int `json: id`
}

type game struct {
	// Exported to JSON
	Move       int
	PlayerTurn int
	Status     status

	// Not exported to JSON
	players []player
}

func newGame() *game {
	g := &game{
		Move:       -1,
		Status:     pairing,
		PlayerTurn: 0,
	}
	return g
}

// Restart game between same players
func (g *game) restartGame(p player) *game {
	g.Status = pairing
	g.players = append(g.players, p)
	g.Move = -1
	return g
}

// In case a player disconnects, then wait for 30 secs in client and do a reset
func (g *game) resetGame() {
	g.players = []player{}
}

// Add player to a game
func (g *game) addPlayer() *player {
	p := player{}
	l := len(g.players)

	if l == 1 {
		g.Status = started
	}

	if l == 0 || l == 1 {
		p.Id = l
		g.players = append(g.players, p)
	}

	return &p
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
	g.switchTurn()
	return g
}

func (g *game) toJSON() []byte {
	return []byte{}
}
