package watcher

type HandlerF func(event Event)
type ErrorHandlerF func(event error)

type IWatcher interface {
	AddHandler(HandlerF)
	AddErrorHandler(ErrorHandlerF)
	AddPathes(*map[string]string)
}
