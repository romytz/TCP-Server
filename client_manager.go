package main

import "net"

// Client represents a connected client, storing the connection and their display name.
// For now, Name is set to the client's address, but can later be customized by the user.
type Client struct {
	Conn net.Conn
	Name string
}

// registerClient adds a new client to the server's clients map.
// The client's name is initially set to their network address.
func (s *Server) registerClient(addr string, conn net.Conn) {
	s.mu.Lock()
	s.clients[addr] = &Client{
		Conn: conn,
		Name: addr, // Default name (can be updated via /name command later)
	}
	s.mu.Unlock()
}

// unregisterClient removes a client from the server's clients map they quit.
func (s *Server) unregisterClient(addr string) {
	s.mu.Lock()
	delete(s.clients, addr)
	s.mu.Unlock()
}
