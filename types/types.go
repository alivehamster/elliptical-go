package types

type Message struct {
	Type   string  `json:"type"`
	String *string `json:"string,omitempty"`
	Chat   *Chat   `json:"chat,omitempty"`
}

type Chat struct {
	ID  string `json:"id"`
	Msg string `json:"msg"`
}
