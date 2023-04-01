package main

import (
	"errors"
	"log"
	"sync"

	"github.com/google/uuid"
)

// Room manages the connections
type room struct {
	id            uuid.UUID
	connectionsMx sync.RWMutex              // mutex to protect connections
	connections   map[uuid.UUID]*connection // Registered connections
}

type Rooms map[uuid.UUID]*room

var rooms = make(Rooms)

// newRoom is the constructor for storing connections in a room
func newRoom(id uuid.UUID) *room {
	r := &room{
		id:            id,
		connectionsMx: sync.RWMutex{},
		connections:   make(map[uuid.UUID]*connection),
	}

	rooms[id] = r

	// go func() {
	// 	for {
	// 		select {
	// 		//waiting for an update of one of the clients in the connection pair
	// 		case <-cp.receiveMove:
	// 		case <-time.After(timeoutBeforeReBroadcast * time.Second): //After x seconds we do broadcast again to see if the opp. is still there
	// 		}

	// 		cp.connectionsMx.RLock()
	// 		for c := range cp.connections {
	// 			select {
	// 			case c.doBroadcast <- true:
	// 			// stop trying to send to this connection after trying for 1 second.
	// 			// if we have to stop, it means that a reader died so remove the connection also.
	// 			case <-time.After(timeoutBeforeConnectionDrop * time.Second):
	// 				cp.removeConnection(c)
	// 			}
	// 		}
	// 		cp.connectionsMx.RUnlock()
	// 	}
	// }()
	return r
}

func JoinRoom(conn *connection) *room {
	if gr, ok := rooms[conn.roomId]; ok {
		return gr
	}

	r := newRoom(conn.roomId)
	r.addConnection(conn)

	return r
}

func GetRoom(conn *connection) (*room, error) {
	if r, ok := rooms[conn.roomId]; ok {
		return r, nil
	}

	return &room{}, errors.New("room was not found")
}

// addConnection adds a players connection to the connectionPair
func (r *room) addConnection(conn *connection) {
	r.connectionsMx.Lock()
	defer r.connectionsMx.Unlock()
	r.connections[conn.id] = conn
	// fmt.Println("[Player] connected")
	log.Println("[Player] connected")
}

// removeConnection removes a player
func (r *room) removeConnection(conn *connection) {
	r.connectionsMx.Lock()
	defer r.connectionsMx.Unlock()
	// if _, ok := r.connections[conn.id]; ok {
	// 	// close(conn.doBroadcast)
	// }
	delete(r.connections, conn.id)

	if len(r.connections) == 0 {
		log.Printf("[Room] cleared: %s", r.id)
		delete(rooms, r.id)
	}

	log.Println("[Player] disconnected")
}
