package tcp

import (
	"bufio"
	"fmt"
	"golang-file-sync/internal/delivery"
	"golang-file-sync/pkg/logger"
	"net"
)

type TcpDelivery struct {
	Port       string
	Comrades   []string
	Connection net.Conn
	Logger     logger.ILogger
	messages   chan delivery.IMessage
}

func (w *TcpDelivery) Run() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", w.Port))
	if err != nil {
		return
	}
	conn, err := listen.Accept()
	w.Connection = conn
	go func() {
		for {
			message, _ := bufio.NewReader(w.Connection).ReadString('\n')
			w.messages <- delivery.IMessage(message)
		}
	}()
}

func (w *TcpDelivery) GetMessageChannel() chan delivery.IMessage {

	return w.messages
}

func (w *TcpDelivery) WriteMessage(message []byte) {
	_, err := w.Connection.Write(message)
	if err != nil {
		w.Logger.Error(err.Error())
	}
}

func NewTcpDelivery(port string, comrades []string, logger logger.ILogger) *TcpDelivery {
	return &TcpDelivery{
		Port:     port,
		Comrades: comrades,
		messages: make(chan delivery.IMessage),
		Logger:   logger,
	}
}
