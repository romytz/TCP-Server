package main

import (
	"fmt"
	"log"
)

// handleMessages listens on the server's message channel (msgch)
// and logs each received message.
func (s *Server) handleMessages() {
	for msg := range s.msgch {
		fmt.Printf("Received message from connection (%s): %s\n", msg.From, string(msg.Payload))
	}
}

// main initializes and starts the TCP server, along with a message handler
// and a goroutine to allow graceful shutdown via ENTER key.
func main() {
	server := NewServer(":3000")

	// Start a background goroutine to handle incoming messages
	go server.handleMessages()

	// Start another goroutine to allow stopping the server with ENTER key
	go func() {
		fmt.Println("Press ENTER to stop the server")
		fmt.Scanln()
		close(server.quitch)
	}()

	err := server.Start()
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
