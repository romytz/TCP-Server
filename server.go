package main

import (
	"bufio"
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
	wg         sync.WaitGroup
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

	fmt.Println("Waiting for all clients to finish...")
	s.wg.Wait()
	fmt.Println("All clients finished.")
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
		s.wg.Add(1)
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

func (s *Server) handleQuit(conn net.Conn, addr string) bool {
	fmt.Printf("Client %s requested quit\n", addr)
	reply := fmt.Sprintf("[Server]: Goodbye, %s!", addr)
	conn.Write([]byte(reply))
	return true
}

func (s *Server) handleList(conn net.Conn, addr string) bool {
	fmt.Printf("Client %s requested list\n", addr)
	s.mu.Lock()
	var clientsList []string
	for clientAddr := range s.clients {
		clientsList = append(clientsList, clientAddr)
	}
	s.mu.Unlock()
	list := fmt.Sprintf("[Server]: %d clients connected:\r\n%s\r\n", len(clientsList), strings.Join(clientsList, "\r\n"))
	conn.Write([]byte(list))
	return false
}

func (s *Server) handleRegularMessage(conn net.Conn, msgBuffer []byte, addr string) bool {
	s.msgch <- Message{
		from:    addr,
		payload: msgBuffer,
	}
	reply := fmt.Sprintf("[Server]: Thank you for your message, %s!\r\n", addr)
	conn.Write([]byte(reply))
	return false
}

func cleanInput(input string) string {
	var result []rune
	for _, r := range input {
		if r == '\b' {
			if len(result) > 0 {
				result = result[:len(result)-1]
			}
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

func (s *Server) handleMessage(conn net.Conn, msgBuffer []byte, addr string) bool {
	// line := strings.TrimSpace(string(msgBuffer))
	line := cleanInput(string(msgBuffer))
	line = strings.TrimSpace(line)
	switch line {
	case "/quit":
		return s.handleQuit(conn, addr)
	case "/list":
		return s.handleList(conn, addr)
	default:
		return s.handleRegularMessage(conn, msgBuffer, addr)
	}
}

func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	defer s.wg.Done()

	addr := conn.RemoteAddr().String()
	s.registerClient(addr, conn)
	fmt.Printf("Client connected: %s\n", addr)

	greeting := fmt.Sprintf("[Server]: Welcome %s!\r\n", addr)
	conn.Write([]byte(greeting))

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text() // line without '\n'
		done := s.handleMessage(conn, []byte(line), addr)
		if done {
			break
		}
	}
	s.unregisterClient(addr)
	fmt.Printf("Client disconnected: %s\n", addr)

}
