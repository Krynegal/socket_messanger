package main

import (
	"bufio"
	"fmt"
	"github.com/Krynegal/socket_messanger/internal/conn"
	"github.com/Krynegal/socket_messanger/internal/message"
	"github.com/Krynegal/socket_messanger/internal/room"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

var roomSize = 4

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:9999")
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	//tunnelConn := tunnel.NewTunnel()
	roomConn := room.NewRoom(roomSize)
	connectedUsers := make(map[int]string, roomSize)
	nameChan := make(chan string)

	var wg sync.WaitGroup

	numConn := 0
	for numConn != roomSize {
		con, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		c := conn.NewConnection(con)

		numConn++
		wg.Add(1)
		go handleName(c, nameChan, &wg, numConn, roomConn.Size())

		roomConn.AddNewConnection(c)
		fmt.Printf("new conn: %s\n", con.RemoteAddr().String())
		fmt.Println(roomConn)

		//if tunnelConn.Conn1 == nil {
		//	tunnelConn.Conn1 = con
		//	fmt.Printf("user 1: %s\n", con.RemoteAddr().String())
		//} else {
		//	tunnelConn.Conn2 = con
		//	fmt.Printf("user 2: %s\n", con.RemoteAddr().String())
		//}
	}

	id := 0
	for name := range nameChan {
		connectedUsers[id] = name
		fmt.Println(connectedUsers)
		id++
	}
	wg.Wait()

	fmt.Println("here we go")

	//firstConn := make(chan string)
	//secondConn := make(chan string)

	mainChan := make(chan *message.Message)

	for i := 0; i < roomConn.Size(); i++ {
		go handleClientRequest(roomConn.Connections[i], mainChan)
	}

	//go handleClientRequest(tunnelConn.Conn1, secondConn)
	//go handleClientRequest(tunnelConn.Conn2, firstConn)

	for getMessage := range mainChan {
		for i := 0; i < roomConn.Size(); i++ {
			if _, err := roomConn.Connections[i].Conn.Write([]byte(fmt.Sprintf("%s-> %s\n", getMessage.Name, strings.TrimSpace(getMessage.Text)))); err != nil {
				log.Printf("failed to respond to client: %v\n", err)
			}
		}
	}

	//for {
	//	select {
	//	case message := <-firstConn:
	//		if _, err := tunnelConn.Conn1.Write([]byte(fmt.Sprintf("%s-> %s\n", connectedUsers[1], strings.TrimSpace(message)))); err != nil {
	//			log.Printf("failed to respond to client: %v\n", err)
	//		}
	//	case message := <-secondConn:
	//		if _, err := tunnelConn.Conn2.Write([]byte(fmt.Sprintf("%s-> %s\n", connectedUsers[0], strings.TrimSpace(message)))); err != nil {
	//			log.Printf("failed to respond to client: %v\n", err)
	//		}
	//	}
	//}

}

func handleName(con *conn.Connection, nameChan chan string, wg *sync.WaitGroup, numConn int, roomSize int) {
	if _, err := con.Conn.Write([]byte("Enter your name:\n")); err != nil {
		log.Printf("failed: %v\n", err)
	}
	clientReader := bufio.NewReader(con.Conn)
	clientRequest, _ := clientReader.ReadString('\n')
	nameChan <- strings.TrimSpace(clientRequest)
	con.Name = strings.TrimSpace(clientRequest)
	fmt.Printf("con.Name: %s", con.Name)
	if numConn == roomSize {
		close(nameChan)
	}
	wg.Done()
}

func handleClientRequest(con conn.Connection, mChan chan *message.Message) {
	defer con.Conn.Close()

	clientReader := bufio.NewReader(con.Conn)

	for {
		clientRequest, err := clientReader.ReadString('\n')

		switch err {
		case nil:
			clientRequest = strings.TrimSpace(clientRequest)
			if clientRequest == ":QUIT" {
				log.Println("client requested server to close the connection so closing")
				return
			} else {
				m := message.NewMessage(con.Name, clientRequest)
				mChan <- m
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
