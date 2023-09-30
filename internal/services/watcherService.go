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
	watcher    watcher.IWatcher
	logger     logger.ILogger
	repository repository.IWatcherRepository
}

func (w *WatcherService) windowsHandlers() watcher.HandlerF {
	return func(event fsnotify.Event) {
		if event.Has(fsnotify.Create) {
			go w.logger.Info(fmt.Sprintf("File %s was created", event.Name))
			go w.repository.Insert(&models.WatcherModel{DirectoryName: "test", FileName: "test", ActionKey: "create"})
		}
		if event.Has(fsnotify.Remove) {
			go w.logger.Info(fmt.Sprintf("File %s was removed", event.Name))
		}
		if event.Has(fsnotify.Write) {
			go w.logger.Info(fmt.Sprintf("File %s was updated", event.Name))
		}
		if event.Has(fsnotify.Rename) {
			go w.logger.Info(fmt.Sprintf("File %s was Rename", event.Name))
		}
		if event.Has(fsnotify.Chmod) {
			go w.logger.Info(fmt.Sprintf("File %s was Chmod", event.Name))
		}
	}
}

func (w *WatcherService) Run() {
	w.watcher.AddHandler(w.windowsHandlers())
}

// TODO ProvideDependencyContainer
func NewWatcherService(w watcher.IWatcher, l logger.ILogger, r repository.IWatcherRepository) *WatcherService {
	return &WatcherService{
		watcher:    w,
		logger:     l,
		repository: r,
	}
}
