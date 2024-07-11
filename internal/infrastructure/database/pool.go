package database

import (
	"database/sql"
)

type Pool interface {
	Builder() *sql.DB
}
