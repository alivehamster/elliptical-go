package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

// eventually save roomid instead of bool
var (
	clients = make(map[*websocket.Conn]bool)
	mu      sync.RWMutex
)

func addClient(conn *websocket.Conn) {
	mu.Lock()
	defer func() {
		broadcastJSON(Message{Type: "Clientnum", Content: fmt.Sprintf("%d", len(clients))})
	}()
	defer mu.Unlock()
	clients[conn] = true
	log.Printf("Client connected. Total clients: %d", len(clients))
}

func removeClient(conn *websocket.Conn) {
	mu.Lock()
	defer func() {
		broadcastJSON(Message{Type: "Clientnum", Content: fmt.Sprintf("%d", len(clients))})
	}()
	defer mu.Unlock()
	delete(clients, conn)
	log.Printf("Client disconnected. Total clients: %d", len(clients))
}

func broadcastJSON(msg Message) {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return
	}
	log.Printf("Broadcasting message: %s", jsonData)

	broadcast(websocket.TextMessage, jsonData)
}

func broadcast(messageType int, message []byte) {
	mu.RLock()
	defer mu.RUnlock()

	for client := range clients {
		if err := client.WriteMessage(messageType, message); err != nil {
			log.Printf("Error broadcasting to client: %v", err)
			// Remove client if write fails
			go func(c *websocket.Conn) {
				c.Close()
				removeClient(c)
			}(client)
		}
	}
}

func main() {
	app := fiber.New()

	app.Use(logger.New())

	app.Static("/", "./frontend/dist")

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		addClient(c)
		defer removeClient(c)

		var (
			msg []byte
			err error
		)
		for {
			if _, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}
			var message Message
			if err := json.Unmarshal(msg, &message); err != nil {
				log.Printf("Error parsing JSON: %v", err)
				continue
			}

			log.Printf("Received message - Type: %s, Content: %s", message.Type, message.Content)
		}
	}))

	log.Println("Server starting on :3000")
	log.Fatal(app.Listen(":3000"))
}
