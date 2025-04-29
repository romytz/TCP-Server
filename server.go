package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

type Server struct {
	listenAddr string
	ln         net.Listener
	quitch     chan struct{}
	msgch      chan Message
	clients    map[string]net.Conn
	mu         sync.Mutex
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		msgch:      make(chan Message, 10),
		clients:    make(map[string]net.Conn),
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
			fmt.Println("Accept error: ", err)
			continue
		}

		fmt.Println("New connection to the server: ", conn.RemoteAddr())

		go s.readLoop(conn)
	}
}

func (s *Server) registerClient(addr string, conn net.Conn) {
	s.mu.Lock()
	s.clients[addr] = conn
	s.mu.Unlock()
}

func (s *Server) unregisterClient(addr string) {
	s.mu.Lock()
	delete(s.clients, addr)
	s.mu.Unlock()
}

func (s *Server) handleMessage(conn net.Conn, msgBuffer []byte, addr string) bool {
	line := strings.TrimSpace(string(msgBuffer))
	if line == "/quit" {
		fmt.Printf("Client requested quit: %s\n", addr)
		reply := fmt.Sprintf("[Server]: Goodbye, %s!\r\n", conn.RemoteAddr().String())
		conn.Write([]byte(reply))
		return true
	}
	s.msgch <- Message{
		from:    conn.RemoteAddr().String(),
		payload: msgBuffer,
	}
	reply := fmt.Sprintf("\r\n[Server]: Thank you for your message, %s!\r\n", conn.RemoteAddr().String())
	conn.Write([]byte(reply))

	return false
}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close() // Make sure to close the connection when done

	addr := conn.RemoteAddr().String()
	s.registerClient(addr, conn)
	fmt.Printf("Client connected: %s\n", addr)

	buff := make([]byte, 2048)
	var msgBuffer []byte
	for {
		n, err := conn.Read(buff)
		if err != nil {
			fmt.Println("Read error: ", err)
			continue
		}
		msgBuffer = append(msgBuffer, buff[:n]...)

		if buff[n-1] == '\n' {
			if s.handleMessage(conn, msgBuffer, addr) {
				break
			}
			msgBuffer = nil
		}
	}

	s.unregisterClient(addr)
	fmt.Printf("Client disconnected: %s\n", addr)

}
