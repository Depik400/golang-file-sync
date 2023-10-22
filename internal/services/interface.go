package services

import "golang-file-sync/pkg/watcher"

type IWatchService interface {
	Run()
	Stop()
	GetChannel() chan watcher.Event
}

type ISyncService interface {
	Run()
}
