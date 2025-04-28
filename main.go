package main

import (
	"fmt"
	"log"
	"net"
)

type Message struct {
	from    string
	payload []byte
}

type Server struct {
	listenAddr string
	ln         net.Listener
	quitch     chan struct{}
	msgch      chan Message
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		msgch:      make(chan Message, 10),
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
	close(s.msgch)
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
	var msgBuffer []byte
	for {
		n, err := conn.Read(buff)
		if err != nil {
			fmt.Println("read error: ", err)
			continue
		}
		// msg := buff[:n]        // Get only the portion of the buffer that was read
		// fmt.Print(string(msg)) // Convert the received bytes into a string and print the message
		// conn.Write(buff[:n])
		msgBuffer = append(msgBuffer, buff[:n]...)

		if buff[n-1] == '\n' {
			s.msgch <- Message{
				from:    conn.RemoteAddr().String(),
				payload: msgBuffer,
			}
			reply := fmt.Sprintf("\r\n[Server]: Thank you for your message, %s!\r\n", conn.RemoteAddr().String())
			conn.Write([]byte(reply))
			msgBuffer = nil
		}
	}
}

func (s *Server) handleMessages() {
	for msg := range s.msgch {
		fmt.Printf("received message from connection (%s): %s", msg.from, string(msg.payload))
	}
}

func main() {
	server := NewServer(":3000")

	go server.handleMessages()

	err := server.Start()
	if err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
