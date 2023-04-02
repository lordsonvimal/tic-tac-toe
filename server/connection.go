package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type connection struct {
	wsConn *websocket.Conn
	roomId uuid.UUID
	id     uuid.UUID // playerId
}

// reads the moves from the clients ws-connection
func (c *connection) read(wg *sync.WaitGroup) {
	for {
		// Reading next move from connection here
		messageType, clientMessage, err := c.wsConn.ReadMessage()

		log.Printf("[Type: %d][Message]: %s", messageType, clientMessage)

		if err != nil {
			log.Println("[ERROR] while read", err)
			wg.Done()

			if r, err := GetRoom(c); err == nil {
				log.Println("[SUCCESS] removing connection with id: ", c.id)
				r.RemoveConnection(c)
			}

			c.close()
			break
		}

		ReadGameState(c, clientMessage)
	}
}

// write something to the connection
func (c *connection) write(data []byte) {
	if err := c.wsConn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Println("[ERROR] while write", err)

		if r, err := GetRoom(c); err == nil {
			log.Println("[SUCCESS] removing connection with id: ", c.id)
			r.RemoveConnection(c)
		}

		c.close()
	}
}

func (c *connection) close() {
	if r, err := GetRoom(c); err == nil {
		r.RemoveConnection(c)
	}
	c.wsConn.Close()
}

func getId(ctx *gin.Context, key string) uuid.UUID {
	paramId, ok := ctx.GetQuery(key)

	if ok {
		if id, err := uuid.Parse(paramId); err == nil {
			return id
		}
	}

	return uuid.New()
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
		CheckOrigin:     func(r *http.Request) bool { return true }, //TODO: Remove in production. Needed for gin proxy
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

	ro := JoinRoom(c)
	c.roomId = ro.Id

	//the websocket connection is always open. Close it from a client request / response
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go c.read(wg)
	wg.Wait()
}
