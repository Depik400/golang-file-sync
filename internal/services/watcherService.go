package services

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"golang-file-sync/internal/models"
	"golang-file-sync/internal/repository"
	"golang-file-sync/pkg/logger"
	"golang-file-sync/pkg/watcher"
)

type WatcherService struct {
	watcher      watcher.IWatcher
	logger       logger.ILogger
	repository   repository.IWatcherRepository
	eventChannel chan watcher.Event
}

func (w *WatcherService) windowsHandlers() watcher.HandlerF {
	return func(event watcher.Event) {
		if event.Has(fsnotify.Create) {
			go w.logger.Info(fmt.Sprintf("File %s was created", event.Name))
			go w.repository.Insert(&models.WatcherModel{
				DirectoryName: event.DirectoryAlias,
				FileName:      event.Filename,
				ActionKey:     "create",
			})
		}
		if event.Has(fsnotify.Remove) {
			go w.logger.Info(fmt.Sprintf("File %s was removed", event.Name))
			go w.repository.Insert(&models.WatcherModel{
				DirectoryName: event.DirectoryAlias,
				FileName:      event.Filename,
				ActionKey:     "remove",
			})
		}
		if event.Has(fsnotify.Write) {
			go w.logger.Info(fmt.Sprintf("File %s was updated", event.Name))
			go w.repository.Insert(&models.WatcherModel{
				DirectoryName: event.DirectoryAlias,
				FileName:      event.Filename,
				ActionKey:     "update",
			})
			w.eventChannel <- event
		}
		if event.Has(fsnotify.Rename) {
			go w.logger.Info(fmt.Sprintf("File %s was Rename", event.Name))
			go w.repository.Insert(&models.WatcherModel{
				DirectoryName: event.DirectoryAlias,
				FileName:      event.Filename,
				ActionKey:     "remove",
			})
		}
		if event.Has(fsnotify.Chmod) {
			go w.logger.Info(fmt.Sprintf("File %s was Chmod", event.Name))
		}
	}
}

func (w *WatcherService) GetChannel() chan watcher.Event {
	return w.eventChannel
}

func (w *WatcherService) Run() {
	w.watcher.AddHandler(w.windowsHandlers())
}

func (w *WatcherService) Stop() {
}

// TODO ProvideDependencyContainer
func NewWatcherService(w watcher.IWatcher, l logger.ILogger, r repository.IWatcherRepository) *WatcherService {
	return &WatcherService{
		watcher:      w,
		logger:       l,
		repository:   r,
		eventChannel: make(chan watcher.Event),
	}
}
