package command

import (
	"net"
	"tcp-server/message"
)

// ServerAPI defines the limited interface that command handlers are allowed to use from the Server.
type ServerAPI interface {
	GetClientList() []string
	Broadcast(msg message.Message)
}

// CommandHandler is the interface that each command (e.g. /list, /quit) must implement.
type CommandHandler interface {
	Execute(api ServerAPI, conn net.Conn, addr string, msgBuffer []byte) bool
}
