package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func clientSend(clientReader *bufio.Reader, iChan chan string) {
	for {
		//fmt.Print("->")
		clientRequest, err := clientReader.ReadString('\n')

		switch err {
		case nil:
			clientRequest = strings.TrimSpace(clientRequest)
			iChan <- clientRequest
		case io.EOF:
			log.Println("client closed the connection")
			return
		default:
			log.Printf("client error: %v\n", err)
			return
		}
	}
}

func clientGet(serverReader *bufio.Reader, oChan chan string) {
	for {
		serverResponse, err := serverReader.ReadString('\n')

		switch err {
		case nil:
			oChan <- strings.TrimSpace(serverResponse)
		case io.EOF:
			log.Println("server closed the connection")
			return
		default:
			log.Printf("server error: %v\n", err)
			return
		}
	}
}

func main() {
	con, err := net.Dial("tcp", "0.0.0.0:9999")
	if err != nil {
		log.Fatalln(err)
	}
	defer con.Close()

	clientReader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(con)

	iChan := make(chan string)
	oChan := make(chan string)

	go clientGet(serverReader, oChan)
	go clientSend(clientReader, iChan)

	for {
		select {
		case clientRequest := <-iChan:
			if _, err = con.Write([]byte(clientRequest + "\n")); err != nil {
				log.Printf("failed to send the client request: %v\n", err)
			}
		case serverResponse := <-oChan:
			fmt.Println(serverResponse)
		}
	}
}
