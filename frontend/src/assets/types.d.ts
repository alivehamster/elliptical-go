export interface Message {
  type: string
  string?: string
  chat?: Chat
  chats?: Chat[]
  room?: Room
  rooms?: Room[]
}

interface Chat {
  id: string
  msg: string
}

interface Room {
  roomid: string
  title: string
}
