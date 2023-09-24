package container

import (
	"golang-file-sync/internal/config"
	"golang-file-sync/internal/services"
	"golang-file-sync/pkg/logger"
	"golang-file-sync/pkg/watcher"
	"sync"
)

type Container struct {
	Logger         logger.ILogger
	WatcherService services.IWatchService
}

var _container *Container

var _group = &sync.WaitGroup{}

var _close = make(chan bool, 1)

func InitContainer(config *config.Config) {
	_container = &Container{}
	initLogger(config)
	initWatcherService(config)
	_group.Add(1)
}

func initLogger(config *config.Config) {
	if config.Logger.Type == "file" {
		_container.Logger = logger.NewFileLogger(&logger.FileLoggerConfig{
			Path:     config.Logger.Path,
			Rotating: config.Logger.Rotating,
			Name:     config.Server.Name,
		})
	} else if config.Logger.Type == "console" {
		_container.Logger = logger.NewConsoleLogger(config.Server.Name)
	}
}

func initWatcherService(config *config.Config) {
	newWatcher := watcher.NewWatcher(_close)
	newWatcher.AddPathes(config.Watcher.Directories)
	_container.WatcherService = services.NewWatcherService(newWatcher, _container.Logger)
}

func GetContainer() *Container {
	return _container
}
