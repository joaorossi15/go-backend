package chat

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/coder/websocket"
	"github.com/joaorossi15/gobh/internal/middleware"
	"github.com/joaorossi15/gobh/internal/user"
)

func ChatHandler(hub *Hub, repo *user.UserR) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet ||
			r.Header.Get("Upgrade") != "websocket" {
			http.Error(w, "websocket only", http.StatusUpgradeRequired)
			return
		}

		v := r.Context().Value(middleware.UserIDKey)
		userName, ok := v.(string)
		if !ok {
			http.Error(w, "error getting user: "+userName, http.StatusBadRequest)
			return
		}
		userID, _, err := repo.Get(r.Context(), userName)

		if err != nil {
			http.Error(w, "error getting user: "+userName, http.StatusBadRequest)
			return
		}

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
			name: userName,
			room: int64(roomID),
		}

		hub.register <- client

		go client.writeMessagesToClients()

		if err := client.readMessagesFromClients(r.Context()); err != nil {
			log.Printf("error reading: %v", err)
		}

		hub.unregister <- client
	}
}

// read frames written in the websocket and writes to hub.broadcast
func (c *Client) readMessagesFromClients(ctx context.Context) error {
	c.conn.SetReadLimit(4096)

	for {
		_, body, err := c.conn.Read(ctx)
		if err != nil {
			return err
		}

		c.hub.broadcast <- Message{
			RoomID:   c.room,
			SenderID: c.id,
			Body:     fmt.Sprintf("%s: %s", c.name, string(body)),
		}
	}
}

// read channel client.send and writes to the websocket for the client
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
			var err error

			err = c.conn.Write(ctx, websocket.MessageText, []byte(msg.Body))

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
			log.Printf("GOOD PING TO CLIENT %v", c.id)
			cancel()
		}
	}
}
