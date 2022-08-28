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

	numConn := 0
	for numConn != 2 {
		con, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		numConn++
		if tunnel.conn1 == nil {
			tunnel.conn1 = con
			fmt.Printf("user 1: %s\n", con.RemoteAddr().String())
		} else {
			tunnel.conn2 = con
			fmt.Printf("user 2: %s\n", con.RemoteAddr().String())
		}
	}

	firstConn := make(chan string)
	secondConn := make(chan string)

	go handleClientRequest(tunnel.conn1, secondConn)
	go handleClientRequest(tunnel.conn2, firstConn)

	for {
		select {
		case message := <-firstConn:
			if _, err := tunnel.conn1.Write([]byte(fmt.Sprintf("-> %s\n", strings.TrimSpace(message)))); err != nil {
				log.Printf("failed to respond to client: %v\n", err)
			}
		case message := <-secondConn:
			if _, err := tunnel.conn2.Write([]byte(fmt.Sprintf("-> %s\n", strings.TrimSpace(message)))); err != nil {
				log.Printf("failed to respond to client: %v\n", err)
			}
		}
	}

}

func handleClientRequest(con net.Conn, oChan chan string) {
	defer con.Close()

	clientReader := bufio.NewReader(con)

	for {
		clientRequest, err := clientReader.ReadString('\n')

		switch err {
		case nil:
			clientRequest = strings.TrimSpace(clientRequest)
			if clientRequest == ":QUIT" {
				log.Println("client requested server to close the connection so closing")
				return
			} else {
				oChan <- clientRequest
			}
		case io.EOF:
			log.Println("client closed the connection by terminating the process")
			return
		default:
			log.Printf("error: %v\n", err)
			return
		}
	}
}
