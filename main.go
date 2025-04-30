package main

import (
	"fmt"
	"log"
)

func (s *Server) handleMessages() {
	for msg := range s.msgch {
		fmt.Printf("Received message from connection (%s): %s\n", msg.from, string(msg.payload))
	}
}

func main() {
	server := NewServer(":3000")

	go server.handleMessages()

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
