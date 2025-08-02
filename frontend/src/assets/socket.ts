import { context, currentRoom } from "./store"
import type { Message } from "./types"

let socket: WebSocket

export function createWebSocket() {
  socket = new WebSocket("ws://localhost:3000/ws")
  socket.addEventListener("open", (event) => {
    console.log("Connected to WebSocket server")
    context.status.code = 0
    context.status.text = "Connected"
  })

  socket.addEventListener("message", (event) => {
    console.log(event.data)
    const data = JSON.parse(event.data)
    handleMessage(data)
  })

  socket.addEventListener("close", () => {
    console.log("Disconnected from WebSocket server")
    context.status.code = 2
    context.status.text = "Disconnected"
  })
}

function handleMessage(data: Message) {
  switch (data.type) {
    case "Clientnum":
      if (!data.string) {
        console.warn("Clientnum message missing string data")
        return
      }
      context.online = parseInt(data.string, 10)
      break
    case "Rooms":
      if (!data.rooms) {
        console.warn("Rooms message missing rooms data")
        return
      }
      context.rooms = data.rooms
      break
    case "NewRoom":
      if (!data.room) {
        console.warn("NewRoom message missing room data")
        return
      }
      context.rooms.push(data.room)
    case "Chat":
      if (!data.chat) {
        console.warn("Chat message missing chat data")
        return
      }
      currentRoom.messages.push({ id: data.chat.id, msg: data.chat.msg })
      break
    default:
      console.warn("Unknown message type:", data.type)
  }
}

export function sendMessage(message: Message) {
  socket.send(JSON.stringify(message))
}
