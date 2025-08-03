import { reactive } from "vue"
import type { Room } from "./types"

export const context = reactive({
  username: "",
  online: 0,
  status: {
    code: 1,
    text: "Connecting...",
  },
  rooms: [] as Room[],
})

export const usernameModal = reactive({
  open: false,
  username: "",
})

export const newRoom = reactive({
  open: false,
})

export const currentRoom = reactive({
  title: "",
  roomid: null as string | null,
  messages: [] as { id: string; msg: string }[],
})

export function reset() {
  context.online = 0
  context.rooms = []
  currentRoom.title = ""
  currentRoom.roomid = null
  currentRoom.messages = []
}