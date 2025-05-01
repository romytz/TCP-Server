package main

import (
	"fmt"
	"net"
	"sync"
	"tcp-server/command"
	"tcp-server/message"
)

// Server manages TCP client connections, incoming messages, and command execution.
type Server struct {
	listenAddr string                            // Address the server listens on (e.g., ":8080")
	ln         net.Listener                      // TCP listener
	quitch     chan struct{}                     // Channel used to signal server shutdown
	msgch      chan message.Message              // Channel for broadcasting messages
	clients    map[string]*Client                // Active clients mapped by address
	mu         sync.Mutex                        // Mutex to protect shared state (clients map)
	wg         sync.WaitGroup                    // WaitGroup to track active client goroutines
	commands   map[string]command.CommandHandler // Registered command handlers
}

// NewServer creates and initializes a new Server instance.
func NewServer(listenAddr string) *Server {
	s := &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		msgch:      make(chan message.Message, 10),
		clients:    make(map[string]*Client),
		commands:   make(map[string]command.CommandHandler),
	}
	s.initCommands()
	return s
}

// initCommands registers supported command handlers for the server.
func (s *Server) initCommands() {
	s.commands = map[string]command.CommandHandler{
		"/list":       &command.ListCommand{},
		"/quit":       &command.QuitCommand{},
		"__default__": &command.DefaultMessageCommand{},
	}
}

// Start launches the server, listens for incoming TCP connections,
// and waits for a shutdown signal.
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()

	s.ln = ln
	go s.acceptLoop() // Start accepting client connections in a separate goroutine

	<-s.quitch // Block until shutdown signal is received

	fmt.Println("Waiting for all clients to finish...")
	s.wg.Wait()
	fmt.Println("All clients finished.")
	close(s.msgch)

	return nil
}

// GetClientList returns the names of all currently connected clients.
func (s *Server) GetClientList() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	var clientsList []string
	for _, client := range s.clients {
		clientsList = append(clientsList, client.Name)
	}
	return clientsList
}

// Broadcast enqueues a message to be processed by the server's message handler.
func (s *Server) Broadcast(msg message.Message) {
	s.msgch <- msg
}
