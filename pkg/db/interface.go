package db

import (
	"github.com/jmoiron/sqlx"
)

type IDatabase interface {
	Connection() *sqlx.DB
}
