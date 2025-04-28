package main

import (
	"fmt"
	"log"
)

func (s *Server) handleMessages() {
	for msg := range s.msgch {
		fmt.Printf("Received message from connection (%s): %s", msg.from, string(msg.payload))
	}
}

func main() {
	server := NewServer(":3000")

	go server.handleMessages()

	err := server.Start()
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
