package chat

import (
	"context"
	"log"

	"github.com/coder/websocket"
)

type Message struct {
	RoomID   int64  `json:"room_id"`
	SenderID int64  `json:"sender_id"`
	Body     string `json:"body"`
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan Message
	name string
	id   int64
	room int64
}

type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan Message
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Message, 64),
	}
}

func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = true
			log.Printf("client %v connected\n", c.id)
		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
				_ = c.conn.CloseNow()
			}
		case msg := <-h.broadcast:
			for c := range h.clients {
				if c.room == msg.RoomID && c.id != msg.SenderID {
					c.send <- msg
				}
			}

		case <-ctx.Done():
			return
		}
	}

}
