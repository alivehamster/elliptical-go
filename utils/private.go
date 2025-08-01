package utils

import (
	"log"

	"github.com/gofiber/contrib/websocket"
)

func broadcast(messageType int, message []byte) {
	mu.RLock()
	defer mu.RUnlock()

	for client := range clients {
		if err := client.WriteMessage(messageType, message); err != nil {
			log.Printf("Error broadcasting to client: %v", err)
			// Remove client if write fails
			go func(c *websocket.Conn) {
				c.Close()
				RemoveClient(c)
			}(client)
		}
	}
}
