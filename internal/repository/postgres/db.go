package postgres

import (
	"context"
	"time"

	"github.com/wb-go/wbf/dbpg"
)

func Connect(ctx context.Context, dsn string) (*dbpg.DB, error) {
	opts := &dbpg.Options{
		MaxOpenConns:    10,
		MaxIdleConns:    1,
		ConnMaxLifetime: 5 * time.Minute,
	}

	db, err := dbpg.New(dsn, nil, opts)
	if err != nil {
		return nil, err
	}

	ctxPing, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := db.Master.PingContext(ctxPing); err != nil {
		_ = db.Master.Close()
		for _, slave := range db.Slaves {
			_ = slave.Close()
		}
		return nil, err
	}

	return db, nil
}
