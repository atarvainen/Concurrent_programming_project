/*
 * Original Program by Mikko Neijonen
 * Refactoring and expanding by Hannu Oksman, Antti Tarvainen
 */

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

// server maintains rooms
type Server struct {
	rooms map[string]time.Time
}

// rooms maintains clients
// they have unique names to allow changing rooms
type Room struct {
	name        string
	clientConns map[net.Conn]time.Time
}

func NewServer() *Server {
	srv := &Server{
		rooms: map[string]time.Time{},
	}
	return srv
}

func NewRoom(n string) *Room {
	room := &Room{
		name:        n,
		clientConns: map[net.Conn]time.Time{},
	}
	return room
}

// adding clients to default room
func (srv *Server) ListenAndServe(address string, room *Room) error {
	// channels for room communications
	connInfo := make(chan net.Conn)
	msgChan := make(chan string)
	go room.handleClientConn(connInfo, msgChan)
	fmt.Println("Listening")
	listener, err := net.Listen("tcp4", address)
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		// new client connection, passing info to a room
		connInfo <- conn
	}
}

// room functionality
func (room *Room) handleClientConn(connInfo chan net.Conn, msgChan chan string) {
	fmt.Println("Handling conns")
	delChan := make(chan net.Conn)
	for {
		select {
		case newConnection := <-connInfo:
			room.clientConns[newConnection] = time.Now()
			go room.handleMessage(newConnection, delChan, msgChan)
			fmt.Printf("%s connected\n", newConnection.LocalAddr)
		case deleteConnenction := <-delChan:
			delete(room.clientConns, deleteConnenction)
			fmt.Printf("%s disconnected\n", deleteConnenction.LocalAddr)
		case newMessage := <-msgChan:
			for peer := range room.clientConns {
				if _, err := fmt.Fprintf(peer, "%s\n", newMessage); err != nil {
					log.Printf("error writing to client %v: %v", peer.RemoteAddr(), err)
				}
			}
		}
	}
}

// message handling
func (room *Room) handleMessage(conn net.Conn, delChan chan net.Conn, msgChan chan string) {
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			// error, deleting connection from the room
			delChan <- conn
			return
		}
		if err != nil {
			log.Printf("client error: %v", err)
			return
		}
		line = strings.TrimSpace(line)

		log.Printf("%s: %v", conn.RemoteAddr(), line)

		// sending message to the room for publishing
		msgChan <- line
	}
}

// creating server and main room, where clients are initially connected
func main() {
	server := NewServer()
	fmt.Println("Server created")
	mainRoom := NewRoom("main")
	fmt.Println("Main room created")
	if err := server.ListenAndServe(":9801", mainRoom); err != nil {
		log.Fatal(err)
	}
}
