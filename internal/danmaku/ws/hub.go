package ws

import (
	"github.com/gorilla/websocket"
)

// RoomMessage carries payloads scoped to a single video room.
type RoomMessage struct {
	VideoID uint
	Payload []byte
}

type Hub struct {
	rooms      map[uint]map[*Client]bool
	Broadcast  chan RoomMessage
	Register   chan *Client
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan RoomMessage),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		rooms:      make(map[uint]map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if _, ok := h.rooms[client.VideoID]; !ok {
				h.rooms[client.VideoID] = make(map[*Client]bool)
			}
			h.rooms[client.VideoID][client] = true
		case client := <-h.Unregister:
			if clients, ok := h.rooms[client.VideoID]; ok {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.Send)
					if len(clients) == 0 {
						delete(h.rooms, client.VideoID)
					}
				}
			}
		case message := <-h.Broadcast:
			if clients, ok := h.rooms[message.VideoID]; ok {
				for client := range clients {
					select {
					case client.Send <- message.Payload:
					default:
						close(client.Send)
						delete(clients, client)
					}
				}
				if len(clients) == 0 {
					delete(h.rooms, message.VideoID)
				}
			}
		}
	}
}

type Client struct {
	Hub     *Hub
	Conn    *websocket.Conn
	Send    chan []byte
	VideoID uint
}
