package container

import (
	"golang-file-sync/internal/config"
	"golang-file-sync/internal/delivery"
	"golang-file-sync/internal/delivery/tcp"
	"golang-file-sync/internal/repository"
	"golang-file-sync/internal/services"
	"golang-file-sync/pkg/db"
	"golang-file-sync/pkg/logger"
	"golang-file-sync/pkg/watcher"
)

type Container struct {
	Logger            logger.ILogger
	Database          db.IDatabase
	Watcher           watcher.IWatcher
	Delivery          delivery.IDelivery
	WatcherService    services.IWatchService
	SyncService       services.ISyncService
	WatcherRepository repository.IWatcherRepository
}

var _container *Container

func InitContainer(config *config.Config) {
	_container = &Container{}
	initLogger(config)
	initDatabase(config)
	initWatcher(config)
	initDelivery(config)
	initRepositories(config)
	initServices(config)
}

func initDatabase(config *config.Config) {
	_container.Database = db.NewSqlLiteDatabase()
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

func initWatcher(config *config.Config) {
	_container.Watcher = watcher.NewWatcher()
	_container.Watcher.AddPathes(&config.Watcher.Directories)
}

func initDelivery(config *config.Config) {
	_container.Delivery = tcp.NewTcpDelivery(config.Server.Port, config.Server.Comrades, _container.Logger)
}

func initRepositories(config *config.Config) {
	_container.WatcherRepository = repository.NewWatcherRepository(_container.Database)
}

func initServices(config *config.Config) {
	_container.WatcherService = services.NewWatcherService(
		_container.Watcher,
		_container.Logger,
		_container.WatcherRepository,
	)
	_container.SyncService = services.NewSyncService(
		_container.Delivery,
		_container.Logger,
		_container.WatcherService,
		config.Watcher.Directories,
	)
}

func GetContainer() *Container {
	return _container
}
