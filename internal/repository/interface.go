package repository

import "golang-file-sync/internal/models"

type IWatcherRepository interface {
	Insert(model *models.WatcherModel)
	Find()
	All()
}
