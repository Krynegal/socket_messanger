package room

import (
	"fmt"
	"github.com/Krynegal/socket_messanger/internal/conn"
)

type Room struct {
	Connections []conn.Connection
	size        int
	capacity    int
}

func NewRoom(cap int) *Room {
	return &Room{
		Connections: make([]conn.Connection, 0, cap),
		size:        0,
		capacity:    cap,
	}
}

func (r *Room) AddNewConnection(conn *conn.Connection) {
	r.Connections = append(r.Connections, *conn)
	r.size++
}

func (r Room) GetCapacity() int {
	return r.capacity
}

func (r Room) GetSize() int {
	return r.size
}

func (r *Room) String() string {
	return fmt.Sprintf("Number of current connections: %v", len(r.Connections))
}
