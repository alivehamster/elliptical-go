package main

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"

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

	app.Static("/", "./dist")

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
				if message.Chat != nil {
					msgid, err := utils.StoreChat(db, message.Chat.ID, message.Chat.Msg)
					if err != nil {
						log.Printf("Error storing chat message: %v", err)
					} else {
						utils.SendRoomMessage(types.Chat{ID: strconv.FormatInt(msgid, 10), Msg: message.Chat.Msg}, message.Chat.ID)
					}

				} else {
					log.Printf("Chat message missing chat field")
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
			case "JoinRoom":
				if message.String != nil {
					roomID := *message.String
					if chats, err := utils.GetChats(db, roomID); err == nil {
						utils.SetClientRoom(c, roomID)
						c.WriteJSON(types.Message{Type: "JoinedRoom", Chats: &chats, String: &roomID})
					} else {
						log.Printf("Error joining room: %v", err)
					}
				} else {
					log.Printf("JoinRoom message missing string field")
				}
			default:
				log.Printf("Unknown message type: %s", message.Type)
			}
		}
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Println("Server starting on", port)
	log.Fatal(app.Listen(":" + port))
}
