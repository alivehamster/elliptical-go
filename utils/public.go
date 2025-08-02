package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/alivehamster/elliptical-go/types"
	"github.com/gofiber/contrib/websocket"
)

// eventually save roomid instead of bool
var (
	clients = make(map[*websocket.Conn]string)
	mu      sync.RWMutex
)

func AddClient(conn *websocket.Conn) {
	mu.Lock()
	defer func() {
		clientCount := fmt.Sprintf("%d", len(clients))
		BroadcastJSON(types.Message{Type: "Clientnum", String: &clientCount})
	}()
	defer mu.Unlock()
	clients[conn] = "home"
	log.Printf("Client connected. Total clients: %d", len(clients))
}

func RemoveClient(conn *websocket.Conn) {
	mu.Lock()
	defer func() {
		clientCount := fmt.Sprintf("%d", len(clients))
		BroadcastJSON(types.Message{Type: "Clientnum", String: &clientCount})
	}()
	defer mu.Unlock()
	delete(clients, conn)
	log.Printf("Client disconnected. Total clients: %d", len(clients))
}

func BroadcastJSON(msg types.Message) {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return
	}
	log.Printf("Broadcasting message: %s", jsonData)

	broadcast(websocket.TextMessage, jsonData)
}
