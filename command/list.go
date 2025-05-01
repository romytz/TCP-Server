package command

import (
	"fmt"
	"net"
	"strings"
)

// ListCommand implements the /list command, which sends the client
// a list of all currently connected users.
type ListCommand struct{}

// Execute handles the /list command:
// It fetches the list of clients from the server and writes it back to the requesting client.
func (l *ListCommand) Execute(api ServerAPI, conn net.Conn, addr string, msgBuffer []byte) bool {
	fmt.Printf("Client %s requested list\n", addr)
	clients := api.GetClientList()
	msg := fmt.Sprintf("[Server]: %d clients connected:\r\n%s\r\n", len(clients), strings.Join(clients, "\r\n"))
	conn.Write([]byte(msg))
	return false
}
