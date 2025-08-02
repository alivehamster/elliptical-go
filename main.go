package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/uuid"

	"github.com/alivehamster/elliptical-go/types"
	"github.com/alivehamster/elliptical-go/utils"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DB_DSN")

	if dsn == "" {
		dsn = "root:password@tcp(localhost:3306)/?parseTime=true"
	}

	db, err := utils.NewDB(dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	app := fiber.New()

	app.Use(logger.New())

	app.Static("/", "./frontend/dist")

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		utils.AddClient(c)
		defer utils.RemoveClient(c)

		rooms, err := utils.GetRooms(db)
		if err != nil {
			log.Printf("Error getting rooms: %v", err)
		}
		c.WriteJSON(types.Message{Type: "Rooms", Rooms: &rooms})

		var msg []byte

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
			case "CreateRoom":
				if message.String != nil {
					roomID, err := utils.CreateRoom(db, *message.String)
					if err != nil {
						log.Printf("Error creating room: %v", err)
					} else {
						room := types.Room{RoomID: strconv.FormatInt(roomID, 10), Title: *message.String}
						utils.BroadcastJSON(types.Message{Type: "NewRoom", Room: &room})
					}
				} else {
					log.Printf("CreateRoom message missing string field")
				}
			default:
				log.Printf("Unknown message type: %s", message.Type)
			}
		}
	}))

	log.Println("Server starting on :3000")
	log.Fatal(app.Listen(":3000"))
}
