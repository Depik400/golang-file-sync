package models

type WatcherModel struct {
	Id            int    `db:"id"`
	DirectoryName string `db:"directory_name"`
	FileName      string `db:"file_name"`
	ActionKey     string `db:"action_key"`
}
