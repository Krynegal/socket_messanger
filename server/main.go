package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
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
	connectedUsers := make(map[int]string, 2)
	nameChan := make(chan string)

	var wg sync.WaitGroup

	numConn := 0
	for numConn != 2 {
		con, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		numConn++
		wg.Add(1)
		go handleName(con, nameChan, &wg, numConn)

		if tunnel.conn1 == nil {
			tunnel.conn1 = con
			fmt.Printf("user 1: %s\n", con.RemoteAddr().String())
		} else {
			tunnel.conn2 = con
			fmt.Printf("user 2: %s\n", con.RemoteAddr().String())
		}
	}

	id := 0
	for name := range nameChan {
		connectedUsers[id] = name
		fmt.Println(connectedUsers)
		id++
	}
	wg.Wait()

	fmt.Println("here we go")

	firstConn := make(chan string)
	secondConn := make(chan string)

	go handleClientRequest(tunnel.conn1, secondConn)
	go handleClientRequest(tunnel.conn2, firstConn)

	for {
		select {
		case message := <-firstConn:
			if _, err := tunnel.conn1.Write([]byte(fmt.Sprintf("%s-> %s\n", connectedUsers[1], strings.TrimSpace(message)))); err != nil {
				log.Printf("failed to respond to client: %v\n", err)
			}
		case message := <-secondConn:
			if _, err := tunnel.conn2.Write([]byte(fmt.Sprintf("%s-> %s\n", connectedUsers[0], strings.TrimSpace(message)))); err != nil {
				log.Printf("failed to respond to client: %v\n", err)
			}
		}
	}

}

func handleName(con net.Conn, nameChan chan string, wg *sync.WaitGroup, numConn int) {
	if _, err := con.Write([]byte("Enter your name:\n")); err != nil {
		log.Printf("failed: %v\n", err)
	}
	clientReader := bufio.NewReader(con)
	clientRequest, _ := clientReader.ReadString('\n')
	nameChan <- strings.TrimSpace(clientRequest)
	if numConn == 2 {
		close(nameChan)
	}
	wg.Done()
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
