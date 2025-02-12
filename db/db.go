package db

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/pressly/goose"
	_ "github.com/qtrix/simulator-wavespread/migrations"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithField("module", "db")

type DB struct {
	config Config
	pool   *pgxpool.Pool
}

func New(ctx context.Context, config Config) (*DB, error) {
	db := &DB{
		config: config,
	}

	cfg, err := db.pgxPoolConfig(config)
	if err != nil {
		return nil, err
	}

	db.pool, err = pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		err = errors.Wrap(err, "unable to connect to db")
		return nil, err
	}

	// if db.config.Reset {
	// 	err = db.ResetDB(ctx)
	// 	if err != nil {
	// 		return nil, errors.Wrap(err, "automigrate")
	// 	}
	// }

	if config.AutoMigrate {
		log.Info("attempting automatic execution of migrations")
		conn, err := sql.Open("postgres", config.ConnectionString)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		err = conn.Ping()
		if err != nil {
			log.Error(err)
			return nil, err
		}

		goose.SetTableName("simulator_version")
		err = goose.Up(conn, "/")
		if err != nil && err != goose.ErrNoNextVersion {
			log.Fatal(err)
		}
		log.Info("database version is up to date")

		err = conn.Close()
		if err != nil {
			return nil, errors.Wrap(err, "closing auto migrations connection")
		}
	}

	return db, nil
}

func (db *DB) Exec(ctx context.Context, sql string, args ...interface{}) error {
	_, err := db.pool.Exec(ctx, sql, args...)
	return err
}

func (db *DB) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return db.pool.Query(ctx, sql, args...)
}

func (db *DB) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return db.pool.QueryRow(ctx, sql, args...)
}

func (db *DB) Begin(ctx context.Context) (pgx.Tx, error) {
	return db.pool.Begin(ctx)
}

func (db *DB) BeginTx(ctx context.Context, options pgx.TxOptions) (pgx.Tx, error) {
	return db.pool.BeginTx(ctx, options)
}

func (db *DB) CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return db.pool.CopyFrom(ctx, tableName, columnNames, rowSrc)
}

func (db *DB) Close() {
	db.pool.Close()
}

func (db *DB) ResetDB(ctx context.Context) error {
	log.Warn("RESETTING DATABASE. THIS WILL RESULT IN ALL DATA BEING LOST")

	drop := `
		DROP TABLE IF EXISTS
			price_data,
			store,
			simulator_version
			CASCADE
	;
	`

	err := db.Exec(ctx, drop)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) pgxPoolConfig(config Config) (*pgxpool.Config, error) {
	pgxCfg, err := pgxpool.ParseConfig(config.ConnectionString)
	if err != nil {
		err = errors.Wrap(err, "unable to parse db config")
		return nil, err
	}

	return pgxCfg, nil
}
