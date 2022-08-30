package message

type Message struct {
	Name string
	Text string
}

func NewMessage(name, text string) *Message {
	return &Message{
		Name: name,
		Text: text,
	}
}
