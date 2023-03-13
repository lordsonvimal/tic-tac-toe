package main

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// connections stores all the hubs
var playerConnections []*connections

type connection struct {
	// Channel which triggers the connection to update the gameState
	doBroadcast chan bool
	// The connectionPair. Holds up to 2 connections.
	cp *connections
	// playerNum represents the players Slot. Either 0 or 1
	playerNum int
}

// reader reads the moves from the clients ws-connection
func (c *connection) reader(wg *sync.WaitGroup, wsConn *websocket.Conn) {
	defer wg.Done()
	for {
		//Reading next move from connection here
		_, clientMoveMessage, err := wsConn.ReadMessage()
		if err != nil {
			break
		}

		field, _ := strconv.ParseInt(string(clientMoveMessage[:]), 10, 32) //Getting FieldValue From Player Action
		c.cp.g.makeMove(int(field))
		c.cp.receiveMove <- true //telling connectionPair to broadcast the gameState
	}
}

// writer broadcasts the current gameState to the two players in a connectionPair
func (c *connection) writer(wg *sync.WaitGroup, wsConn *websocket.Conn) {
	defer wg.Done()
	for range c.doBroadcast {
		sendGameStateToConnection(wsConn, c)
	}
}

// getConnectionPairWithEmptySlot looks trough all connectionPairs and finds one which has only 1 player
// if there is none a new connectionPair is created and the player is added to that pair
func getConnectionPairWithEmptySlot() (*connections, int) {
	sizeBefore := len(playerConnections)
	// find connections with 1 player first and pair if possible
	for _, h := range playerConnections {
		if len(h.connections) == 1 {
			log.Printf("Players paired")
			return h, len(h.connections)
		}
	}

	//TODO: I need to remove orphaned connectionPairs from the stack

	// if no emtpy slow was found at all, we create a new connectionPair
	h := newConnections()
	playerConnections = append(playerConnections, h)
	log.Printf("Player seated in new connectionPair no. %v", len(playerConnections))
	return playerConnections[sizeBefore], 0
}

// ServeHTTP is the routers HandleFunc for websocket connections
// connections are upgraded to websocket connections and the player is added
// to a connection pair
func ServeHTTP(ctx *gin.Context) {
	w, r := ctx.Writer, ctx.Request

	// upgrader is needed to upgrade the HTTP Connection to a websocket Connection
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true }, //TODO: Remove in production. Needed for gin proxy
	}

	//Upgrading HTTP Connection to websocket connection
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading %s", err)
		return
	}

	//Adding Connection to connectionPair
	cp, pn := getConnectionPairWithEmptySlot()
	c := &connection{doBroadcast: make(chan bool), cp: cp, playerNum: pn}
	c.cp.addConnection(c)

	//If the connectionPair existed before but one player was disconnected
	//we can now reinitialize the gameState after the remaining player has
	//been paired again
	// if c.cp.g.Status == resetWaitPaired {
	// 	c.cp.g = newGame()
	// 	log.Println("gamestate resetted")
	// }

	//inform the gameState about the new player
	c.cp.g.addPlayer()
	//telling connectionPair to broadcast the gameState
	c.cp.receiveMove <- true

	//creating the writer and reader goroutines
	//the websocket connection is open as long as these goroutines are running
	var wg sync.WaitGroup
	wg.Add(2)
	go c.writer(&wg, wsConn)
	go c.reader(&wg, wsConn)
	wg.Wait()
	wsConn.Close()
}

// sendGameStateToConnection broadcasts the current gameState as JSON to all players
// within a connectionPair
func sendGameStateToConnection(wsConn *websocket.Conn, c *connection) {
	err := wsConn.WriteMessage(websocket.TextMessage, c.cp.g.toJSON())
	//removing connection if updating gameState fails
	if err != nil {
		c.cp.removeConnection(c)
	}
}
