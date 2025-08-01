import { context } from "./store"
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
      context.online = parseInt(data.content, 10)
      break
    default:
      console.warn("Unknown message type:", data.type)
  }
}

export function sendMessage(message: Message) {
  socket.send(JSON.stringify(message))
}

