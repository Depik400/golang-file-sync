package watcher

import "github.com/fsnotify/fsnotify"

type Event struct {
	*fsnotify.Event
	Directory      string
	Filename       string
	DirectoryAlias string
}
