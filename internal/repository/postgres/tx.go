package postgres

import (
	"context"
	"database/sql"

	"github.com/wb-go/wbf/dbpg"
)

type Transactor struct {
	db *dbpg.DB
}

func NewTransactor(db *dbpg.DB) *Transactor {
	return &Transactor{db: db}
}

func (t *Transactor) WithinTx(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error {
	return t.db.WithTx(ctx, func(tx *sql.Tx) error {
		return fn(ctx, tx)
	})
}
