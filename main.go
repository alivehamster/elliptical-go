package main

import (
	"encoding/json"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/uuid"

	"github.com/alivehamster/elliptical-go/types"
	"github.com/alivehamster/elliptical-go/utils"
)

func main() {
	app := fiber.New()

	app.Use(logger.New())

	app.Static("/", "./frontend/dist")

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		utils.AddClient(c)
		defer utils.RemoveClient(c)

		var (
			msg []byte
			err error
		)
		for {
			if _, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}
			var message types.Message
			if err := json.Unmarshal(msg, &message); err != nil {
				log.Printf("Error parsing JSON: %v", err)
				continue
			}
			switch message.Type {
			case "SendChat":
				if message.String != nil {
					utils.BroadcastJSON(types.Message{Type: "Chat", Chat: &types.Chat{ID: uuid.New().String(), Msg: *message.String}})
				} else {
					log.Printf("Chat message missing string field")
				}
			default:
				log.Printf("Unknown message type: %s", message.Type)
			}
		}
	}))

	log.Println("Server starting on :3000")
	log.Fatal(app.Listen(":3000"))
}
