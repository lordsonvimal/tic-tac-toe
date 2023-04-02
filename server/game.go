package main

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
)

const (
	GAME_PENDING status = "GAME_PENDING"
	GAME_STARTED status = "GAME_STARTED"
	PLAYER_MOVED status = "PLAYER_MOVED"
)

const SENDER_GAME sender = "GAME"

type ticTacToe string

const (
	X ticTacToe = "X"
	O ticTacToe = "O"
)

type player map[uuid.UUID]ticTacToe

type game struct {
	Turn   uuid.UUID
	Data   int
	Player player
	Status status
}

func ReadGameState(conn *connection, data []byte) {
	newRoom := room{}
	if err := json.Unmarshal(data, &newRoom); err != nil {
		log.Printf("[ERROR] Unmarshalling data to room: %s", string(data))
		return
	}

	r, err := GetRoom(conn)

	if err != nil {
		log.Printf("[ERROR] Fetching Room from connection : %s", conn.id)
		return
	}

	switch s := newRoom.Status; s {
	case PLAYER_MOVED:
		r.Status = newRoom.Status
		r.Sender = SENDER_GAME
		r.Data = newRoom.Data
		r.Broadcast(r.ToJSON())
	}
}

func NewGame() game {
	return game{
		Data:   -1,
		Player: make(player),
		Status: GAME_PENDING,
	}
}

func StartGame(r *room, conn *connection) {
	if len(r.connections) == 1 {
		r.Game.Player[conn.id] = X
	}

	if len(r.connections) == 2 && r.Game.Status == GAME_PENDING {
		r.Game.Turn = conn.id
		r.Game.Status = GAME_STARTED
		r.Game.Player[conn.id] = O
		r.Sender = SENDER_GAME
		r.Broadcast(r.ToJSON())
	}
}
