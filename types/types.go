package types

type Message struct {
	Type   string  `json:"type"`
	String *string `json:"string,omitempty"`
	Chat   *Chat   `json:"chat,omitempty"`
	Chats  *[]Chat `json:"chats,omitempty"`
	Room   *Room   `json:"room,omitempty"`
	Rooms  *[]Room `json:"rooms,omitempty"`
}

type Chat struct {
	ID  string `json:"id"`
	Msg string `json:"msg"`
}

type Room struct {
	RoomID string `json:"roomid"`
	Title  string `json:"title"`
}
