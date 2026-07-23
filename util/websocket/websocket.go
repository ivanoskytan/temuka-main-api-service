package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSMessage struct {
	ConversationID int    `json:"conversation_id"`
	ParticipantID  int    `json:"participant_id"`
	SenderID       int    `json:"sender_id"`
	Text           string `json:"text"`
}

type Client struct {
	Hub            *Hub
	Conn           *websocket.Conn
	Send           chan []byte
	UserID         int
	ConversationID int
}

type Hub struct {
	Rooms      map[int]map[*Client]bool
	Broadcast  chan WSMessage
	Register   chan *Client
	Unregister chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[int]map[*Client]bool),
		Broadcast:  make(chan WSMessage),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if h.Rooms[client.ConversationID] == nil {
				h.Rooms[client.ConversationID] = make(map[*Client]bool)
			}
			h.Rooms[client.ConversationID][client] = true
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			if clients, ok := h.Rooms[client.ConversationID]; ok {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.Send)
					if len(clients) == 0 {
						delete(h.Rooms, client.ConversationID)
					}
				}
			}
			h.mu.Unlock()

		case msg := <-h.Broadcast:
			h.mu.RLock()
			clients := h.Rooms[msg.ConversationID]
			payload, err := json.Marshal(msg)
			if err == nil {
				for client := range clients {
					select {
					case client.Send <- payload:
					default:
						close(client.Send)
						delete(clients, client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (c *Client) ReadPump(handleIncomingMsg func(msg WSMessage)) {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		var wsMsg WSMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			log.Println("Invalid WS message payload:", err)
			continue
		}

		wsMsg.ConversationID = c.ConversationID
		wsMsg.SenderID = c.UserID

		handleIncomingMsg(wsMsg)

		c.Hub.Broadcast <- wsMsg
	}
}

func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
	}()

	for message := range c.Send {
		w, err := c.Conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		w.Write(message)

		if err := w.Close(); err != nil {
			return
		}
	}
}
