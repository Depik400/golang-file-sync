package repository

import (
	"golang-file-sync/internal/models"
	"golang-file-sync/pkg/db"
	"log"
)

type WatcherRepository struct {
	database db.IDatabase
}

func (r *WatcherRepository) Insert(model *models.WatcherModel) {
	_, err := r.database.Connection().NamedExec(
		"Insert into system_actions (directory_name, file_name, action_key) values (:directory_name, :file_name,:action_key)",
		model,
	)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func (r *WatcherRepository) Find() {

}

func (r *WatcherRepository) All() {

}

func NewWatcherRepository(database db.IDatabase) *WatcherRepository {
	return &WatcherRepository{
		database,
	}
}
