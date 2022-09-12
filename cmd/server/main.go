package main

import (
	"bufio"
	"fmt"
	"github.com/Krynegal/socket_messanger/configs"
	"github.com/Krynegal/socket_messanger/internal/conn"
	"github.com/Krynegal/socket_messanger/internal/message"
	"github.com/Krynegal/socket_messanger/internal/room"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

var roomSize = 2

func main() {
	conf := configs.Get()

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", "0.0.0.0", conf.ServerPort))
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	roomConn := room.NewRoom(roomSize)
	connectedUsers := make(map[int]string, roomSize)
	nameChan := make(chan string)

	var wg sync.WaitGroup

	var numConn int32 = 0
	for numConn != int32(roomSize) {
		con, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		c := conn.NewConnection(con)

		numConn++
		wg.Add(1)

		go handleName(c, nameChan, &wg, numConn, roomConn)

		fmt.Printf("new conn: %s\n", con.RemoteAddr().String())
	}

	id := 0
	for name := range nameChan {
		connectedUsers[id] = name
		fmt.Println(connectedUsers)
		id++
	}
	wg.Wait()

	//fmt.Println("here we go")

	mainChan := make(chan *message.Message)

	for i := 0; i < roomConn.Size(); i++ {
		go handleClientRequest(roomConn.Connections[i], mainChan)
	}

	for getMessage := range mainChan {
		for i := 0; i < roomConn.Size(); i++ {
			if getMessage.SenderID == roomConn.Connections[i].ID {
				continue
			}
			if _, err := roomConn.Connections[i].Conn.Write([]byte(fmt.Sprintf("%s-> %s\n", getMessage.Name, strings.TrimSpace(getMessage.Text)))); err != nil {
				log.Printf("failed to respond to client: %v\n", err)
			}
		}
	}
}

func handleName(con *conn.Connection, nameChan chan string, wg *sync.WaitGroup, numConn int32, room *room.Room) {
	if _, err := con.Conn.Write([]byte("Enter your name:\n")); err != nil {
		log.Printf("failed: %v\n", err)
	}

	clientReader := bufio.NewReader(con.Conn)
	clientRequest, _ := clientReader.ReadString('\n')
	nameChan <- strings.TrimSpace(clientRequest)
	con.ID = numConn

	con.Name = strings.TrimSpace(clientRequest)

	room.AddNewConnection(con)

	if numConn == int32(room.Size()) {
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
				mChan <- message.NewMessage(con.ID, con.Name, clientRequest)
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
