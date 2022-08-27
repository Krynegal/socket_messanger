package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type tunnel struct {
	conn1 net.Conn
	conn2 net.Conn
}

func NewTunnel() *tunnel {
	return &tunnel{
		conn1: nil,
		conn2: nil,
	}
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:9999")
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	tunnel := NewTunnel()

	for {
		con, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		if tunnel.conn1 == nil {
			tunnel.conn1 = con
			fmt.Printf("conn1: %s\n", con.RemoteAddr().String())
		} else {
			tunnel.conn2 = con
			fmt.Printf("conn2: %s\n", con.RemoteAddr().String())
		}

		if tunnel.conn1 != nil && tunnel.conn2 != nil {
			go handleClientRequest(tunnel.conn1, tunnel.conn2)
		}
	}
}

func handleClientRequest(con1 net.Conn, con2 net.Conn) {
	defer con1.Close()
	defer con2.Close()

	clientReader := bufio.NewReader(con1)

	for {
		clientRequest, err := clientReader.ReadString('\n')

		switch err {
		case nil:
			clientRequest := strings.TrimSpace(clientRequest)
			if clientRequest == ":QUIT" {
				log.Println("client requested server to close the connection so closing")
				return
			} else {
				log.Println(clientRequest)
			}
		case io.EOF:
			log.Println("client closed the connection by terminating the process")
			return
		default:
			log.Printf("error: %v\n", err)
			return
		}

		if _, err = con2.Write([]byte(fmt.Sprintf("%s\n", strings.TrimSpace(clientRequest)))); err != nil {
			log.Printf("failed to respond to client: %v\n", err)
		}
	}
}
