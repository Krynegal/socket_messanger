package message

type Message struct {
	SenderID int32
	Name     string
	Text     string
}

func NewMessage(senderID int32, name, text string) *Message {
	return &Message{
		SenderID: senderID,
		Name:     name,
		Text:     text,
	}
}
