package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fiffu/mailtl/app/infra"
	"github.com/fiffu/mailtl/app/storage/sqlite"
	"go.uber.org/fx"
)

const name = "storage"

type Storage interface {
	OnStart(context.Context) error
	OnStop(context.Context) error
	DB() *sql.DB
}

func WithTransaction(ctx context.Context, storage Storage, callback func(txn *sql.Tx) error) error {
	txn, err := storage.DB().BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	callbackErr := callback(txn)
	switch callbackErr {
	case nil:
		return txn.Commit()
	default:
		rollbackErr := txn.Rollback()

		if rollbackErr != nil {
			return fmt.Errorf("%T during rollback: %w (root cause: %T<%v>)", rollbackErr, rollbackErr, callbackErr, callbackErr)
		} else {
			return rollbackErr
		}
	}
}

func NewStorage(lc fx.Lifecycle, cfg infra.RootConfig, logger infra.RootLogger) (Storage, error) {
	log := infra.NewLogger(logger, name)
	client := sqlite.NewClient(cfg.SQLiteDSN)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Infof(ctx, "Connecting to database: %s", cfg.SQLiteDSN)
			if err := client.OnStart(ctx); err != nil {
				return err
			}
			if err := Migrate(ctx, client); err != nil {
				return err
			}
			log.Infof(ctx, "Migrations done")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Infof(ctx, "Closing database connection")
			return client.OnStop(ctx)
		},
	})

	return client, nil
}
