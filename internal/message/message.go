package message

type Message struct {
	SenderID int
	Name     string
	Text     string
}

func NewMessage(senderID int, name, text string) *Message {
	return &Message{
		SenderID: senderID,
		Name:     name,
		Text:     text,
	}
}
