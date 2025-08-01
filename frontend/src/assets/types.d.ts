export interface Message {
  type: string
  string?: string
  chat?: Chat
}

interface Chat {
  id: string
  msg: string
}