package relational

import (
	"context"
	"database/sql"
)

type ConnectionManager interface {
	Create(ctx context.Context) (*sql.DB, error)
	CreateInDatabase(ctx context.Context, database string) (*sql.DB, error)
}

type BulkInserter interface {
	Insert(ctx context.Context, db *sql.DB, table *Table, rows []map[string]interface{}) error
}
