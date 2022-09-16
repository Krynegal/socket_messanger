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

var roomCapacity = 3

func main() {
	conf := configs.Get()

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", "0.0.0.0", conf.ServerPort))
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()

	roomConn := room.NewRoom(roomCapacity)
	connectedUsers := make(map[int]string, roomCapacity)
	nameChan := make(chan string)
	//newConnChan := make(chan conn.Connection)

	var wg sync.WaitGroup

	for roomConn.GetSize() != roomConn.GetCapacity() {
		con, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		c := conn.NewConnection(con)

		roomConn.Size++
		wg.Add(1)
		go handleName(c, nameChan, &wg, roomConn)

		fmt.Printf("new conn: %s\n", con.RemoteAddr().String())
	}

	id := 0
	for i := 0; i < roomConn.GetCapacity(); i++ {
		connectedUsers[id] = <-nameChan
		fmt.Println(connectedUsers)
		id++
	}
	wg.Wait()

	//fmt.Println("here we go")

	mainChan := make(chan *message.Message)

	for i := 0; i < roomConn.GetCapacity(); i++ {
		if roomConn.Connections[i].ID != roomConn.GetLastConnID() {
			if _, err = roomConn.Connections[i].Conn.Write([]byte("All users are here!\n")); err != nil {
				log.Printf("failed to respond to client: %v\n", err)
			}
		}
	}

	for i := 0; i < roomConn.GetCapacity(); i++ {
		go handleClientRequest(roomConn.Connections[i], mainChan)
	}

	for getMessage := range mainChan {
		for i := 0; i < roomConn.GetCapacity(); i++ {
			if getMessage.SenderID == roomConn.Connections[i].ID {
				continue
			}
			if _, err := roomConn.Connections[i].Conn.Write([]byte(fmt.Sprintf("%s-> %s\n", getMessage.Name, strings.TrimSpace(getMessage.Text)))); err != nil {
				log.Printf("failed to respond to client: %v\n", err)
			}
		}
	}
}

func handleName(con *conn.Connection, nameChan chan string, wg *sync.WaitGroup, room *room.Room) {
	defer wg.Done()
	if _, err := con.Conn.Write([]byte("Enter your name:\n")); err != nil {
		log.Printf("failed: %v\n", err)
	}
	con.ID = room.Size
	clientReader := bufio.NewReader(con.Conn)
	clientRequest, _ := clientReader.ReadString('\n')
	con.Name = strings.TrimSpace(clientRequest)
	room.AddNewConnection(con)
	nameChan <- con.Name
	fmt.Println("roomCap", room.GetCapacity())
	fmt.Println("roomSize", room.GetSize())
	if len(room.Connections) != room.GetCapacity() {
		if _, err := con.Conn.Write([]byte("Waiting for other users...\n")); err != nil {
			log.Printf("failed: %v\n", err)
		}
	}
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
