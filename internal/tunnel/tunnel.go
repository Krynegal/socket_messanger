package tunnel

import "net"

type tunnel struct {
	Conn1 net.Conn
	Conn2 net.Conn
}

func NewTunnel() *tunnel {
	return &tunnel{
		Conn1: nil,
		Conn2: nil,
	}
}
