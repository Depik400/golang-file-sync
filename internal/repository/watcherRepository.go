package repository

import (
	"golang-file-sync/internal/models"
	"golang-file-sync/pkg/db"
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
		return
	}
}

func (r *WatcherRepository) Find() {

}

func (r *WatcherRepository) All() {

}

func (r *WatcherRepository) GetListOfLastActions() *[]models.WatcherModel {
	sql := `
select *
from system_actions
where created_at = (select max(created_at)
                    from system_actions s
                    where s.directory_name = system_actions.directory_name
                      and s.file_name = system_actions.file_name
                    group by s.directory_name, s.file_name)
`
	dest := &[]models.WatcherModel{}
	err := r.database.Connection().Select(dest, sql)
	if err != nil {
		return nil
	}
	return dest
}

func NewWatcherRepository(database db.IDatabase) *WatcherRepository {
	return &WatcherRepository{
		database,
	}
}
