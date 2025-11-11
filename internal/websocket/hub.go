package websocket

import "log"

// MessagePayload defines the structure for messages sent through the hub
type MessagePayload struct {
	SenderID   int    `json:"sender_id"`
	ReceiverID int    `json:"receiver_id"`
	Content    string `json:"content"`
	Timestamp  string `json:"timestamp"`
}

// Hub maintains active clients and broadcasts messages to them
type Hub struct {
	Clients    map[int]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan MessagePayload
}

// Run runs the main event loop for the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.UserID] = client
			log.Printf("User %d connected via WebSocket", client.UserID)

		case client := <-h.Unregister:
			if _, ok := h.Clients[client.UserID]; ok {
				delete(h.Clients, client.UserID)
				close(client.Send)
				log.Printf("User %d disconnected", client.UserID)
			}

		case message := <-h.Broadcast:
			receiver := h.Clients[message.ReceiverID]
			if receiver != nil {
				select {
				case receiver.Send <- message:
				default:
					close(receiver.Send)
					delete(h.Clients, receiver.UserID)
				}
			} else {
				log.Printf("Receiver %d is offline, message persisted only in DB", message.ReceiverID)
			}
		}
	}
}
