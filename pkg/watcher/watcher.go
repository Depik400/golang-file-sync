package watcher

import (
	"github.com/fsnotify/fsnotify"
	"log"
)

type Watcher struct {
	Handlers    *[]HandlerF
	Errors      *[]ErrorHandlerF
	FileWatcher *fsnotify.Watcher
	Close       chan bool
}

func (w *Watcher) AddHandler(handler HandlerF) {
	*w.Handlers = append(*w.Handlers, handler)
}

func (w *Watcher) AddErrorHandler(handler ErrorHandlerF) {
	*w.Errors = append(*w.Errors, handler)
}

func (w *Watcher) AddPathes(pathes []string) {
	for _, s := range pathes {
		w.FileWatcher.Add(s)
	}

}

func (w *Watcher) initWatcher() {
	watcher, err := fsnotify.NewWatcher()
	w.FileWatcher = watcher
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case event, ok := <-w.FileWatcher.Events:
				if !ok {
					return
				}
				for _, funcs := range *w.Handlers {
					go (funcs)(event)
				}
			case err, ok := <-w.FileWatcher.Errors:
				if !ok {
					return
				}
				for _, funcs := range *w.Errors {
					go (funcs)(err)
				}
			case <-w.Close:
				w.FileWatcher.Close()
				return
			}
		}
	}()

	if err != nil {
		log.Fatal(err)
	}
}

func NewWatcher(close chan bool) *Watcher {
	watcher := &Watcher{Close: close, Errors: &[]ErrorHandlerF{}, Handlers: &[]HandlerF{}}
	watcher.initWatcher()
	return watcher
}
