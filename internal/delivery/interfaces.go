package delivery

import "net"

const (
	CommandPing = "ping"
)

type IDelivery interface {
	Run()
	GetMessageChannel() chan IMessage
	WriteMessage([]byte, net.Conn)
}

type IMessage struct {
	Message    string
	Connection net.Conn
	Cmd        string
	Metadata   struct {
		LocationName string
		FileName     string
		FileSize     string
	}
}
