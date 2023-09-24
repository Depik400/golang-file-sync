package watcher

import "github.com/fsnotify/fsnotify"

type HandlerF func(event fsnotify.Event)
type ErrorHandlerF func(event error)

type IWatcher interface {
	AddHandler(HandlerF)
	AddErrorHandler(ErrorHandlerF)
	AddPathes([]string)
}
