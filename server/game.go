package main

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

const (
	GAME_PENDING       status = "GAME_PENDING"
	GAME_STARTED       status = "GAME_STARTED"
	PLAYER_TURN        status = "PLAYER_TURN"
	PLAYER_TURN_CHANGE status = "PLAYER_TURN_CHANGE"
)

const SENDER_GAME sender = "GAME"

type ticTacToe string

const (
	X ticTacToe = "X"
	O ticTacToe = "O"
)

type player map[uuid.UUID]ticTacToe

type game struct {
	Turn   ticTacToe
	Data   int
	Player player
	Status status
}

func SwapTurn(r *room) {
	if r.Game.Turn == X {
		r.Game.Turn = O
	} else {
		r.Game.Turn = X
	}
	r.Game.Status = PLAYER_TURN_CHANGE

	r.Broadcast(r.ToJSON())
}

func ReadGameState(conn *connection, data []byte) {
	newRoom := room{}
	if err := json.Unmarshal(data, &newRoom); err != nil {
		log.Printf("[ERROR] Unmarshalling data to room: %s", string(data))
		return
	}

	r, err := GetRoom(conn)

	if err != nil {
		log.Printf("[ERROR] Fetching Room with id: %s from connection : %s", conn.roomId, conn.id)
		return
	}

	switch s := newRoom.Game.Status; s {
	case PLAYER_TURN:
		r.Game.Status = PLAYER_TURN
		r.Sender = SENDER_GAME
		r.Game.Data = newRoom.Game.Data
		r.Broadcast(r.ToJSON())

		SwapTurn(r)

		return
	}
}

func NewGame() game {
	return game{
		Data:   -1,
		Player: make(player),
		Status: GAME_PENDING,
	}
}

func StartGame(r *room, conn *connection) bool {
	if len(r.connections) == 1 {
		r.Game.Player[conn.id] = X
	}

	if len(r.connections) == 2 && r.Game.Status == GAME_PENDING {
		r.Game.Turn = O
		r.Game.Status = GAME_STARTED
		r.Game.Player[conn.id] = O
		r.Sender = SENDER_GAME
		return true
	}

	return false
}
