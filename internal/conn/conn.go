package conn

import (
	"github.com/Krynegal/socket_messanger/internal/message"
	"net"
)

type Connection struct {
	net.Conn
	ID   int32
	Name string
	Ch   chan message.Message
}

func NewConnection(con net.Conn) *Connection {
	return &Connection{
		Conn: con,
		Ch:   make(chan message.Message),
	}
}
