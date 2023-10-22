package watcher

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"path/filepath"
)

type Watcher struct {
	Handlers    *[]HandlerF
	Errors      *[]ErrorHandlerF
	FileWatcher *fsnotify.Watcher
	Directories *map[string]string
	Pause       chan bool
}

func (w *Watcher) AddHandler(handler HandlerF) {
	*w.Handlers = append(*w.Handlers, handler)
}

func (w *Watcher) AddErrorHandler(handler ErrorHandlerF) {
	*w.Errors = append(*w.Errors, handler)
}

func (w *Watcher) AddPathes(m *map[string]string) {
	for k, s := range *m {
		err := w.FileWatcher.Add(s)
		(*w.Directories)[k] = s
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}

func (w *Watcher) resolveEvent(event *fsnotify.Event) *Event {
	dir, file := filepath.Split(event.Name)
	var findedKey string
	for key, directory := range *w.Directories {
		if directory == dir || directory == dir[:len(dir)-1] {
			findedKey = key
			dir = directory
		}
	}
	return &Event{
		Event:          event,
		Directory:      dir,
		Filename:       file,
		DirectoryAlias: findedKey,
	}
}

func (w *Watcher) initWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	w.FileWatcher = watcher

	go func() {
		for {
			select {
			case event, ok := <-w.FileWatcher.Events:
				if !ok {
					return
				}
				for _, funcs := range *w.Handlers {
					resolvedEvent := w.resolveEvent(&event)
					go (funcs)(*resolvedEvent)
				}
			case err, ok := <-w.FileWatcher.Errors:
				if !ok {
					return
				}
				for _, funcs := range *w.Errors {
					go (funcs)(err)
				}
			}
		}
	}()

	if err != nil {
		log.Fatal(err)
	}
}

func NewWatcher() *Watcher {
	watcher := &Watcher{
		Pause:       make(chan bool),
		Handlers:    &[]HandlerF{},
		Errors:      &[]ErrorHandlerF{},
		Directories: &map[string]string{},
	}
	watcher.initWatcher()
	return watcher
}
