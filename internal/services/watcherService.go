package services

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"golang-file-sync/pkg/logger"
	"golang-file-sync/pkg/watcher"
	"time"
)

type WatcherService struct {
	watcher watcher.IWatcher
	logger  logger.ILogger
	created map[string]interface{}
	updated map[string]interface{}
	removed map[string]interface{}
}

func (w *WatcherService) windowsTransaction() {
	toCreate := "\n"
	toDelete := "\n"
	toUpdate := "\n"
	for file := range w.created {
		toCreate += fmt.Sprintf("\t\t - File '%s'\n", file)
	}
	for file := range w.updated {
		toUpdate += fmt.Sprintf("\t\t - File '%s'\n", file)
	}
	for file := range w.removed {
		toDelete += fmt.Sprintf("\t\t - File '%s'", file)
	}
	go w.logger.Info(fmt.Sprintf("\n\t\tcreated: %s\t\tupdated: %s\t\tdeleted: %s", toCreate, toUpdate, toDelete))
}

func (w *WatcherService) windowsHandlers() watcher.HandlerF {
	hasTimer := false
	return func(event fsnotify.Event) {
		if !hasTimer {
			go func() {
				<-time.After(200 * time.Millisecond)
				hasTimer = false
				w.windowsTransaction()
			}()
			hasTimer = true
		}
		if event.Has(fsnotify.Create) {
			w.created[event.Name] = nil
			//go w.logger.Info(fmt.Sprintf("File %s was created", event.Name))
		}
		if event.Has(fsnotify.Write) {
			w.updated[event.Name] = nil
			//go w.logger.Info(fmt.Sprintf("File %s was updated", event.Name))
		}
		if event.Has(fsnotify.Remove) {
			go w.logger.Info(fmt.Sprintf("File %s was removed", event.Name))
			_, isCreated := w.created[event.Name]
			_, isUpdated := w.updated[event.Name]
			if isUpdated {
				delete(w.updated, event.Name)
			}
			if isCreated {
				delete(w.created, event.Name)
			} else {
				w.removed[event.Name] = nil
			}

		}
		if event.Has(fsnotify.Chmod) {
			w.updated[event.Name] = nil
		}
		if event.Has(fsnotify.Rename) {
			w.removed[event.Name] = nil
		}
	}
}

func (w *WatcherService) Run() {
	w.watcher.AddHandler(w.windowsHandlers())
}

func NewWatcherService(w watcher.IWatcher, l logger.ILogger) *WatcherService {
	return &WatcherService{
		watcher: w,
		logger:  l,
		removed: make(map[string]interface{}),
		created: make(map[string]interface{}),
		updated: make(map[string]interface{}),
	}
}
