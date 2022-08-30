package room

import (
	"fmt"
	"github.com/Krynegal/socket_messanger/internal/conn"
)

type Room struct {
	Connections []conn.Connection
	size        int
}

func NewRoom(size int) *Room {
	return &Room{
		make([]conn.Connection, 0, size),
		size,
	}
}

func (r *Room) AddNewConnection(conn *conn.Connection) {
	r.Connections = append(r.Connections, *conn)
}

func (r Room) Size() int {
	return r.size
}

func (r *Room) String() string {
	return fmt.Sprintf("Number of current connections: %v", len(r.Connections))
}
