package main

import (
	"log"
	"sync"
	"time"
)

// timeoutBeforeReBroadcast sets the time in seconds after where we rebroadcast the gameState
// to all clients. This way we see if the opponent is still there
const timeoutBeforeReBroadcast = 5 //TODO: should probably be set higher in real world usages...
// timeoutBeforeConnectionDrop sets the time in seconds after after we drop a connection
// which is not answering
const timeoutBeforeConnectionDrop = 1

// connections handles the update of the gameState between two players
type connections struct {
	// the mutex to protect connections
	connectionsMx sync.RWMutex
	// Registered connections.
	connections map[*connection]struct{}
	// Inbound messages from the connections.
	receiveMove chan bool
	g           *game
}

// newConnections is the constructor for the connectionPair struct
func newConnections() *connections {
	cp := &connections{
		connectionsMx: sync.RWMutex{},
		receiveMove:   make(chan bool),
		connections:   make(map[*connection]struct{}),
		g:             newGame(),
	}

	go func() {
		for {
			select {
			//waiting for an update of one of the clients in the connection pair
			case <-cp.receiveMove:
			case <-time.After(timeoutBeforeReBroadcast * time.Second): //After x seconds we do broadcast again to see if the opp. is still there
			}

			cp.connectionsMx.RLock()
			for c := range cp.connections {
				select {
				case c.doBroadcast <- true:
				// stop trying to send to this connection after trying for 1 second.
				// if we have to stop, it means that a reader died so remove the connection also.
				case <-time.After(timeoutBeforeConnectionDrop * time.Second):
					cp.removeConnection(c)
				}
			}
			cp.connectionsMx.RUnlock()
		}
	}()
	return cp
}

// addConnection adds a players connection to the connectionPair
func (h *connections) addConnection(conn *connection) {
	h.connectionsMx.Lock()
	defer h.connectionsMx.Unlock()
	// TODO: Should be checking if the same user gets paired to himself
	// TODO: by reloading the page. We could achieve that with setting
	// TODO: cookies to re-identify users
	h.connections[conn] = struct{}{}

}

// removeConnection removes a players connection from the connectionPair
func (h *connections) removeConnection(conn *connection) {
	h.connectionsMx.Lock()
	defer h.connectionsMx.Unlock()
	if _, ok := h.connections[conn]; ok {
		delete(h.connections, conn)
		close(conn.doBroadcast)
	}
	log.Println("Player disconnected")
	h.g.resetGame()
}
