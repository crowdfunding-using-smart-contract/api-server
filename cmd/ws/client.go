package ws

import (
	"encoding/json"
	"fmt"
	"fund-o/api-server/internal/entity"
	"fund-o/api-server/internal/http/middleware"
	"fund-o/api-server/pkg/token"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 10000
)

var newline = []byte{'\n'}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client represents the websockets client at the server
type Client struct {
	// The actual websockets connection.
	ID    string
	conn  *websocket.Conn
	hub   *Hub
	send  chan []byte
	rooms map[*Room]bool
}

func newClient(conn *websocket.Conn, hub *Hub, id string) *Client {
	return &Client{
		ID:    id,
		conn:  conn,
		hub:   hub,
		send:  make(chan []byte, 256),
		rooms: make(map[*Room]bool),
	}
}

func (client *Client) readPump() {
	defer func() {
		client.disconnect()
	}()

	client.conn.SetReadLimit(maxMessageSize)

	_ = client.conn.SetReadDeadline(time.Now().Add(pongWait))

	client.conn.SetPongHandler(func(string) error {
		_ = client.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// Start endless read loop, waiting for messages from client
	for {
		_, jsonMessage, err := client.conn.ReadMessage()
		if err != nil {
			break
		}
		client.handleNewMessage(jsonMessage)
	}

}

func (client *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = client.conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.send:
			_ = client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				_ = client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)

			// Attach queued chat messages to the current websockets message.
			n := len(client.send)
			for i := 0; i < n; i++ {
				_, _ = w.Write(newline)
				_, _ = w.Write(<-client.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			_ = client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (client *Client) disconnect() {
	client.hub.unregister <- client
	for room := range client.rooms {
		room.unregister <- client
	}
	close(client.send)
	_ = client.conn.Close()
}

// ServeWs handles websockets requests from clients requests.
func ServeWs(hub *Hub, ctx *gin.Context) {
	userID := ctx.MustGet(middleware.AuthorizationPayloadKey).(*token.Payload).UserID
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := newClient(conn, hub, userID)

	go client.writePump()
	go client.readPump()

	hub.register <- client
}

func (client *Client) handleNewMessage(jsonMessage []byte) {
	var message entity.ReceivedMessage
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
	}

	fmt.Printf("Received message: %+v\n", message)

	switch message.Action {
	case JoinUserAction:
		client.handleJoinRoomMessage(message)
	case LeaveUserAction:
		client.handleLeaveRoomMessage(message)
	}
}

func (client *Client) handleJoinRoomMessage(message entity.ReceivedMessage) {
	roomName := message.Room

	room := client.hub.findRoomById(roomName)
	if room == nil {
		room = client.hub.createRoom(roomName)
	}

	client.rooms[room] = true

	room.register <- client
}

func (client *Client) handleLeaveRoomMessage(message entity.ReceivedMessage) {
	roomName := message.Room

	room := client.hub.findRoomById(roomName)
	if room == nil {
		return
	}

	delete(client.rooms, room)

	room.unregister <- client
}
