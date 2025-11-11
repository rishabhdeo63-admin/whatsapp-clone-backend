package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Client represents a connected user
type Client struct {
	UserID int
	Conn   *websocket.Conn
	Send   chan MessagePayload
	Hub    *Hub
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		var msg MessagePayload
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading from user %d: %v", c.UserID, err)
			break
		}
		msg.SenderID = c.UserID
		msg.Timestamp = time.Now().Format(time.RFC3339)
		c.Hub.Broadcast <- msg
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()

	for msg := range c.Send {
		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Marshal error for user %d: %v", c.UserID, err)
			continue
		}
		if err := c.Conn.WriteMessage(websocket.TextMessage, jsonMsg); err != nil {
			log.Printf("Write error for user %d: %v", c.UserID, err)
			break
		}
	}
}
