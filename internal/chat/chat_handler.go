package chat

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type ctxKey string

const userIDKey ctxKey = "userID"

func ChatHandler(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet ||
			r.Header.Get("Upgrade") != "websocket" {
			http.Error(w, "websocket only", http.StatusUpgradeRequired)
			return
		}

		userID, _ := strconv.Atoi(r.PathValue("userID"))
		roomID, _ := strconv.Atoi(r.PathValue("roomID"))

		// handshake
		c, err := websocket.Accept(w, r, nil)
		if err != nil {
			log.Printf("upgrade: %v", err)
			return
		}

		client := &Client{
			hub:  hub,
			conn: c,
			send: make(chan Message, 32),
			id:   int64(userID),
			room: int64(roomID),
		}

		hub.register <- client

		// another goroutine
		go client.writeMessagesToClients()

		if err := client.readMessagesFromClients(); err != nil {
			log.Printf("error reading: %v", err)
		}

		hub.unregister <- client
	}
}

func (c *Client) readMessagesFromClients() error {
	c.conn.SetReadLimit(4096)

	for {
		var body struct {
			Body string `json:"body"`
		}
		err := wsjson.Read(context.Background(), c.conn, &body)
		if err != nil {
			return err
		}

		c.hub.broadcast <- Message{
			RoomID:   c.room,
			SenderID: c.id,
			Body:     body.Body,
		}
	}
}

func (c *Client) writeMessagesToClients() {
	ticker := time.NewTicker(10 * time.Second)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close(websocket.StatusNormalClosure, "")
	}()

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			err := wsjson.Write(ctx, c.conn, msg)
			cancel()
			if err != nil {
				return
			}

		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			if err := c.conn.Ping(ctx); err != nil {
				cancel()
				return
			}
			cancel()
		}
	}
}

