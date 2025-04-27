package main

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	listenAddr string
	ln         net.Listener
	quitch     chan struct{}
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr) // Start listening for incoming TCP connections
	if err != nil {
		return err
	}
	defer ln.Close() // Ensure listener is closed when Start exits
	s.ln = ln
	go s.acceptLoop()
	<-s.quitch // Block here until the server is told to quit
	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("accept error: ", err)
			continue
		}

		fmt.Println("new connection to the server: ", conn.RemoteAddr())

		go s.readLoop(conn)
	}
}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close() // Make sure to close the connection when done
	buff := make([]byte, 2048)
	for {
		n, err := conn.Read(buff)
		if err != nil {
			fmt.Println("read error: ", err)
			continue
		}
		msg := buff[:n]        // Get only the portion of the buffer that was read
		fmt.Print(string(msg)) // Convert the received bytes into a string and print the message
	}
}

func main() {
	server := NewServer(":3000")
	err := server.Start()
	if err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
