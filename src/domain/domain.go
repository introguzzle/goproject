package domain

import (
	"database/sql"
)

type Entity interface {
	FlushRows(*sql.Rows) error
	Table() string
	FlushRow(*sql.Row) error
	SetId(int64)
	GetId() int64
	Mappings() map[string]string
}
