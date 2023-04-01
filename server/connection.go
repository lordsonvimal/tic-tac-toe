package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type connection struct {
	wsConn *websocket.Conn
	roomId uuid.UUID
	id     uuid.UUID // playerId
}

// read reads the moves from the clients ws-connection
func (c *connection) read() {
	for {
		// Reading next move from connection here
		_, clientMessage, err := c.wsConn.ReadMessage()

		if err != nil {
			log.Fatalln(clientMessage)
			log.Fatalln(err)
			c.close()
			break
		}

		// field, _ := strconv.ParseInt(string(clientMoveMessage[:]), 10, 32) //Getting FieldValue From Player Action
		// c.cp.g.makeMove(int(field))
		// c.cp.receiveMove <- true //telling connectionPair to broadcast the gameState
	}
}

// write something to the connection
func (c *connection) write() {
	// for range c.doBroadcast {
	// 	sendGameStateToConnection(wsConn, c)
	// }
}

func (c *connection) close() {
	if r, err := GetRoom(c); err == nil {
		r.removeConnection(c)
	}
	c.wsConn.Close()
}

// getConnectionPairWithEmptySlot looks trough all connectionPairs and finds one which has only 1 player
// if there is none a new connectionPair is created and the player is added to that pair
// func getConnectionPairWithEmptySlot() (*room, int) {
// 	sizeBefore := len(rooms)
// 	// find connections with 1 player first and pair if possible
// 	for _, h := range rooms {
// 		if len(h.connections) == 1 {
// 			log.Printf("Players paired")
// 			return h, len(h.connections)
// 		}
// 	}

//TODO: I need to remove orphaned connectionPairs from the stack

// if no emtpy slow was found at all, we create a new connectionPair
// h := newRoom()
// rooms = append(rooms, h)
// log.Printf("Player seated in new connectionPair no. %v", len(rooms))
// return rooms[sizeBefore], 0
// }

func getId(ctx *gin.Context, key string) uuid.UUID {
	var id uuid.UUID

	paramId, ok := ctx.GetQuery(key)

	if ok {
		id, _ = uuid.Parse(paramId)
	} else {
		id = uuid.New()
	}

	return id
}

// createWS is the routers HandleFunc for websocket connections
// connections are upgraded to websocket connections and the player is added
// to a connection pair
func createWS(ctx *gin.Context) {
	w, r := ctx.Writer, ctx.Request

	// upgrader is needed to upgrade the HTTP Connection to a websocket Connection
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// CheckOrigin:     func(r *http.Request) bool { return true }, //TODO: Remove in production. Needed for gin proxy
	}

	//Upgrading HTTP Connection to websocket connection
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading %s", err)
		return
	}

	c := &connection{
		id:     getId(ctx, "player_id"),
		roomId: getId(ctx, "room_id"),
		wsConn: wsConn,
	}

	JoinRoom(c)

	//the websocket connection is always open. Close it from a client request / response
	go c.read()
}

// sendGameStateToConnection broadcasts the current gameState as JSON to all players
// within a connectionPair
// func sendGameStateToConnection(wsConn *websocket.Conn, c *connection) {
// 	err := wsConn.WriteMessage(websocket.TextMessage, c.cp.g.toJSON())
// 	//removing connection if updating gameState fails
// 	if err != nil {
// 		c.cp.removeConnection(c)
// 	}
// }
