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

func NewServer() *Server {
	srv := &Server{
		clientConns: map[net.Conn]time.Time{},
	}
	return srv
}

/*
 * server maintains information about chat rooms
 * new struct for each chat room, handles connection map
 * struct provides a channel which is used to update client connection data
 * thus reducing drastically how often the map is accessed
*/
type Server struct {
	clientConns map[net.Conn]time.Time
}

func (srv *Server) ListenAndServe(address string) error {
	listener, err := net.Listen("tcp4", address)
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		srv.clientConns[conn] = time.Now()
		go srv.handleClientConn(conn)
	}
}

func (srv *Server) handleClientConn(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			delete(srv.clientConns, conn)
			return
		}
		if err != nil {
			log.Printf("client error: %v", err)
			return
		}
		line = strings.TrimSpace(line)

		log.Printf("%s: %v", conn.RemoteAddr(), line)

		for peer := range srv.clientConns {
			if _, err := fmt.Fprintf(peer, "%s\n", line); err != nil {
				log.Printf("error writing to client %v: %v", conn.RemoteAddr(), err)
			}
		}
	}
}

func main() {
	server := NewServer()
	if err := server.ListenAndServe(":9801"); err != nil {
		log.Fatal(err)
	}
}
