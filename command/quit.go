package command

import (
	"fmt"
	"net"
)

// QuitCommand implements the /quit command, which disconnects the client from the server.
type QuitCommand struct{}

// Execute handles the /quit command:
// It sends a goodbye message to the client and signals the server to close the connection.
func (q *QuitCommand) Execute(api ServerAPI, conn net.Conn, addr string, msgBuffer []byte) bool {
	fmt.Printf("Client %s requested quit\n", addr)
	reply := fmt.Sprintf("[Server]: Goodbye, %s!", addr)
	conn.Write([]byte(reply))
	return true
}
