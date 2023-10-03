package delivery

type IDelivery interface {
	Run()
	GetMessageChannel() chan IMessage
	WriteMessage([]byte)
}

type IMessage string
