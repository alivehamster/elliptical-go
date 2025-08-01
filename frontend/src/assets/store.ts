import { reactive } from "vue"

export const context = reactive({
  username: "",
  online: 0,
  roomid: 1,
  status: {
    code: 1,
    text: "Connecting..."
  },
  rooms: {
    length: 0,
  }
})

export const usernameModal = reactive({
  open: false,
  username: "",
})

export const currentRoom = reactive({
  title: "test",
  messages: [] as { id: string; msg: string }[],
})