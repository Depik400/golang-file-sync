package delivery

import (
	"github.com/fsnotify/fsnotify"
	"net"
)

const (
	CommandPing       = "ping"
	CommandNewComrade = "new_comrade"
	CommandCreateFile = "create"
)

type IDelivery interface {
	Run()
	PushMessageAll(message *IMessage)
	PushFileToAll(message *IMessage) *[]net.Conn
	GetMessageChannel() chan IMessage
	WriteMessage([]byte, net.Conn)
}

type IMessage struct {
	Message    string
	Connection net.Conn
	Cmd        string
	Metadata   IMetaData
}

type IMetaData struct {
	LocationName string
	FileName     string
	FileSize     string
	Action       fsnotify.Op
}
