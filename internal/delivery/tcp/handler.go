package tcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"golang-file-sync/internal/delivery"
	"golang-file-sync/pkg/logger"
	"log"
	"net"
	"time"
)

type TcpDelivery struct {
	Port     string
	Comrades []string
	Sockets  map[string]net.Conn
	Logger   logger.ILogger
	messages chan delivery.IMessage
}

func (w *TcpDelivery) Run() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", w.Port))
	if err != nil {
		return
	}

	go func() {
		for {
			conn, err := listen.Accept()

			if err != nil {
				log.Fatal(err)
			}
			go func() {
				for {
					messageString, err := bufio.NewReader(conn).ReadString('\n')
					if err != nil {
						w.Logger.Error(err.Error())
						break
					}
					message := w.prepareMessage(messageString)
					message.Connection = conn
					//w.Logger.Info(fmt.Sprintf("Local address:%s\nRemote address: %s", conn.LocalAddr().String(), conn.RemoteAddr().String()))
					if isHandled := w.messageHandler(message); !isHandled {
						w.messages <- *message
					}
					if message.Cmd == "file" {
						break
					}
				}
			}()

		}
	}()
	w.initConnections(true)
	go w.Ping()
}

func (w *TcpDelivery) messageHandler(message *delivery.IMessage) bool {
	isHandled := false
	switch message.Cmd {
	case "new_comrade":
		isHandled = true
		w.Logger.Info("new server is up")
		w.reinitConnections()
		break
	case "ping":
		isHandled = true
		break
	}

	return isHandled
}

func (w *TcpDelivery) PushMessageAll(message *delivery.IMessage) {
	for _, socket := range w.Sockets {
		obj, err := json.Marshal(*message)
		if err != nil {
			w.Logger.Error(err.Error())
			continue
		}
		obj = append(obj, []byte("\n")...)
		socket.Write(obj)
	}
}

func (w *TcpDelivery) PushFileToAll(message *delivery.IMessage) *[]net.Conn {
	fileSockets := []net.Conn{}
	for _, comrade := range w.Comrades {
		if conn, err := net.Dial("tcp", comrade); err == nil {
			fileSockets = append(fileSockets, conn)
		}
	}
	msgString, _ := json.Marshal(*message)
	for _, socket := range fileSockets {
		socket.Write(append(msgString, []byte("\n")...))
	}
	return &fileSockets
}

func (w *TcpDelivery) prepareMessage(str string) *delivery.IMessage {
	message := delivery.IMessage{}
	err := json.Unmarshal([]byte(str), &message)
	if err != nil {
		w.Logger.Error(err.Error())
		return nil
	}
	return &message
}

func (w *TcpDelivery) Ping() {
	for {
		for key, conn := range w.Sockets {
			_, err := (conn).Write([]byte("{\"cmd\":\"ping\"}\n"))
			if err != nil {
				w.Logger.Error(((conn).RemoteAddr().String()) + " is dead")
				delete(w.Sockets, key)
			} else {
				//   w.Logger.Info(((conn).RemoteAddr().String()) + " is alive")
			}
		}
		time.Sleep(time.Second)
	}
}

func (w *TcpDelivery) initConnections(sendMessage bool) {
	for _, comrade := range w.Comrades {
		conn, err := net.Dial("tcp", comrade)
		if err != nil {
			w.Logger.Error(err.Error())
			continue
		}
		if sendMessage {
			marshal, _ := json.Marshal(map[string]string{"cmd": "new_comrade"})
			if err != nil {
				return
			}
			marshal = append(marshal, []byte("\n")...)
			conn.Write(marshal)
		}
		w.Sockets[comrade] = conn
	}
}

func (w *TcpDelivery) reinitConnections() {
	for _, comrade := range w.Comrades {
		if _, ok := w.Sockets[comrade]; ok {
			continue
		}
		conn, err := net.Dial("tcp", comrade)
		if err != nil {
			w.Logger.Error(err.Error())
			continue
		}
		w.Logger.Info(fmt.Sprintf("reconnect to %s", comrade))
		w.Sockets[comrade] = conn
	}
}

func (w *TcpDelivery) GetMessageChannel() chan delivery.IMessage {

	return w.messages
}

func (w *TcpDelivery) WriteMessage(message []byte, conn net.Conn) {
	_, err := conn.Write(message)
	if err != nil {
		w.Logger.Error(err.Error())
	}
}

func NewTcpDelivery(port string, comrades []string, logger logger.ILogger) *TcpDelivery {
	return &TcpDelivery{
		Port:     port,
		Comrades: comrades,
		messages: make(chan delivery.IMessage),
		Sockets:  make(map[string]net.Conn),
		Logger:   logger,
	}
}

func RemoveIndex[T any](s []T, index int) []T {
	ret := make([]T, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}
