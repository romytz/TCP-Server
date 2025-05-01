package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
)

// acceptLoop continuously accepts incoming TCP connections
// and spawns a read loop goroutine for each client.
func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			// If the listener was closed, exit the loop without printing an error
			if errors.Is(err, net.ErrClosed) {
				return
			}
			fmt.Println("Accept error:", err)
			continue
		}

		fmt.Println("New connection to the server:", conn.RemoteAddr())
		s.wg.Add(1) // Track this client so we can wait for all to finish on shutdown
		go s.readLoop(conn)
	}
}

// readLoop handles input from a single client connection.
// It registers the client, greets them, reads messages, and unregisters them upon disconnection.
func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close() // Ensure the connection is closed when this function exits
	defer s.wg.Done()  // Decrement the WaitGroup counter on exit

	addr := conn.RemoteAddr().String()
	s.registerClient(addr, conn)
	fmt.Printf("Client connected: %s\n", addr)

	greeting := fmt.Sprintf("[Server]: Welcome %s!\r\n", addr)
	conn.Write([]byte(greeting))

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		done := s.handleMessage(conn, []byte(line), addr)
		if done {
			break // Stop reading if command indicated to quit (e.g. /quit)
		}
	}

	s.unregisterClient(addr)
	fmt.Printf("Client disconnected: %s\n", addr)
}

// handleMessage determines the appropriate command to execute for a given input.
// It falls back to a default message handler if no command matches.
func (s *Server) handleMessage(conn net.Conn, msgBuffer []byte, addr string) bool {
	line := cleanInput(string(msgBuffer)) // Clean line to simulate client typing (handle backspaces)
	handler, ok := s.commands[line]
	if !ok {
		handler = s.commands["__default__"]
	}
	return handler.Execute(s, conn, addr, msgBuffer)

}

// cleanInput simulates processing of real-time input where backspace characters ('\b') occur.
// It removes the character before each backspace to mimic actual typing correction.
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
	line := strings.TrimSpace(string(result))
	return line
}
