package command

import (
	"fmt"
	"net"
	"tcp-server/message"
)

// DefaultMessageCommand handles any unrecognized input from the client.
// It treats the message as a generic chat message and broadcasts it.
type DefaultMessageCommand struct{}

// Execute handles non-command messages:
// It broadcasts the message to the server and responds with an acknowledgment.
func (d *DefaultMessageCommand) Execute(api ServerAPI, conn net.Conn, addr string, msgBuffer []byte) bool {
	api.Broadcast(message.Message{
		From:    addr,
		Payload: msgBuffer,
	})
	reply := fmt.Sprintf("[Server]: Thank you for your message, %s!\r\n", addr)
	conn.Write([]byte(reply))
	return false
}
