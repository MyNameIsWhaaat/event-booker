package repository

import (
	"context"
	"database/sql"
)

type Transactor interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error
}
