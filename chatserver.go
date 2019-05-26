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
	rooms map[string]*Room
}

// rooms maintains clients
// they have unique names to allow changing rooms
type Room struct {
	name        string
	clientConns map[net.Conn]time.Time
}

func NewServer() *Server {
	srv := &Server{
		rooms: make(map[string]*Room),
	}
	return srv
}

func (srv *Server) NewRoom(n string) *Room {
	room := &Room{
		name:        n,
		clientConns: map[net.Conn]time.Time{},
	}
	srv.rooms[n] = room
	return room
}

// adding clients to default room
func (srv *Server) ListenAndServe(address string, room *Room) error {
	// channels for room communications
	connInfo := make(chan net.Conn)
	msgChan := make(chan string)

	fmt.Println("Listening")
	listener, err := net.Listen("tcp4", address)

	go room.handleClientConn(connInfo, msgChan, srv)
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
func (room *Room) handleClientConn(connInfo chan net.Conn, msgChan chan string, srv *Server) {
	fmt.Println("Handling conns")
	delChan := make(chan net.Conn)
	for {
		select {
		case newConnection := <-connInfo:
			room.clientConns[newConnection] = time.Now()
			go room.handleMessage(newConnection, delChan, msgChan, srv)
			fmt.Printf("%s connected\n", newConnection.LocalAddr)

			// welcome message for connecting users including room name
			if _, err := fmt.Fprintf(newConnection, "Welcome to room %v\n", room.name); err != nil {
				log.Printf("error writing to client %v: %v", newConnection.RemoteAddr(), err)
			}
		case deleteConnenction := <-delChan:
			delete(room.clientConns, deleteConnenction)
			fmt.Printf("%s disconnected\n", deleteConnenction.LocalAddr)
		case newMessage := <-msgChan:
			for peer := range room.clientConns {
				fmt.Println("Printing to room %v", room.name)
				if _, err := fmt.Fprintf(peer, "%s\n", newMessage); err != nil {
					log.Printf("error writing to client %v: %v", peer.RemoteAddr(), err)
				}
			}
		}
	}
}

// remove client from old room and switch to new room
func handleRoomSwitch(conn net.Conn, room *Room, newRoom *Room) {
	delete(room.clientConns, conn)
	for peer := range room.clientConns {
		fmt.Println("old room %v", peer.RemoteAddr())
	}
	newRoom.clientConns[conn] = time.Now()
	for peer := range newRoom.clientConns {
		fmt.Println("new room %v", peer.RemoteAddr())
	}
	fmt.Printf("%s switched to room %s\n", conn.LocalAddr, newRoom.name)
}

// message handling
func (room *Room) handleMessage(conn net.Conn, delChan chan net.Conn, msgChan chan string, srv *Server) {
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
		splitLine := strings.Split(line, " ")

		// if message = /create roomName, create a new room
		if splitLine[0] == "/create" && len(splitLine) > 1 {
			newRoom := srv.NewRoom(splitLine[1])
			log.Printf("Created room %v", splitLine[1])
			handleRoomSwitch(conn, room, newRoom)
		} else if splitLine[0] == "/join" && len(splitLine) > 1 {

			//handleRoomSwitch(conn, room, newRoom)

		} else if splitLine[0] == "/leave" && len(splitLine) > 1 {



		} else if splitLine[0] == "/destroy" && len(splitLine) > 1 {



		} else {
			log.Printf("%s: %v", conn.RemoteAddr(), line)

			// sending message to the room for publishing
			msgChan <- line
		}
	}
}

// creating server and main room, where clients are initially connected
func main() {
	server := NewServer()
	fmt.Println("Server created")
	mainRoom := server.NewRoom("main")
	fmt.Println("Main room created")
	if err := server.ListenAndServe(":9801", mainRoom); err != nil {
		log.Fatal(err)
	}
}
