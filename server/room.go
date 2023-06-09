package main

import (
	"encoding/json"
	"errors"
	"log"
	"sync"

	"github.com/google/uuid"
)

const (
	CONNECTION_CONNECTED    status = "CONNECTION_CONNECTED"
	CONNECTION_DISCONNECTED status = "CONNECTION_DISCONNECTED"
	ROOM_REMOVED            status = "ROOM_REMOVED"
)

// Room manages the connections
type room struct {
	Connection    uuid.UUID // Player that modified the status
	Game          game
	Id            uuid.UUID
	Status        status
	connectionsMx sync.RWMutex              // mutex to protect connections
	connections   map[uuid.UUID]*connection // Registered connections
}

type Rooms map[uuid.UUID]*room

var rooms = make(Rooms)

// NewRoom is the constructor for storing connections in a room
func NewRoom(id uuid.UUID) *room {
	r := &room{
		Id:            id,
		Game:          NewGame(),
		connectionsMx: sync.RWMutex{},
		connections:   make(map[uuid.UUID]*connection),
	}

	rooms[id] = r
	return r
}

func JoinRoom(conn *connection) *room {
	if r, ok := rooms[conn.roomId]; ok {
		log.Printf("[Player] rejoined room %s", conn.id)
		return r
	}

	for k := range rooms {
		rooms[k].AddConnection(conn)
		return rooms[k]
	}

	r := NewRoom(conn.roomId)
	r.AddConnection(conn)

	return r
}

func GetRoom(conn *connection) (*room, error) {
	if r, ok := rooms[conn.roomId]; ok {
		return r, nil
	}

	return &room{}, errors.New("room was not found")
}

func (r *room) UpdateConnection(s status, conn *connection) {
	r.Status = s
	r.Connection = conn.id
}

// adds a players connection to the room
func (r *room) AddConnection(conn *connection) {
	r.connectionsMx.Lock()
	defer r.connectionsMx.Unlock()
	r.connections[conn.id] = conn
	r.UpdateConnection(CONNECTION_CONNECTED, conn)
	log.Println("[Player] connected")
	r.Broadcast(r.ToJSON())
	if isStarted := StartGame(r, conn); isStarted {
		r.Broadcast(r.ToJSON())
	}
}

// removes a player
func (r *room) RemoveConnection(conn *connection) {
	r.connectionsMx.Lock()
	defer r.connectionsMx.Unlock()

	delete(r.connections, conn.id)
	r.UpdateConnection(CONNECTION_DISCONNECTED, conn)

	if len(r.connections) == 0 {
		log.Printf("[Room] removed: %s", r.Id)
		delete(rooms, r.Id)
		r.UpdateConnection(ROOM_REMOVED, conn)
	}

	log.Printf("[Player] disconnected: %s", conn.id)
	r.Broadcast(r.ToJSON())
}

func (r *room) Broadcast(data []byte) {
	log.Printf("Broadcasting data to room: %s", r.Id)
	for _, conn := range r.connections {
		conn.write(data)
	}
}

func (r *room) ToJSON() []byte {
	j, err := json.Marshal(r)
	if err != nil {
		log.Printf("[Error] in marshalling json: %s", err)
	}
	return j
}
