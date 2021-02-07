package database

import (
	"context"
	"database/sql"
)

// SaveFunc is used to commit or rollback changes on execute.
type SaveFunc func(bool)

// Database is a high-level interface wrapping the sql.DB type.
type Database interface {
	Multiple(ctx context.Context, query string, args ...interface{}) (Rows, error)
	Single(ctx context.Context, query string, args ...interface{}) (Row, error)
	Execute(ctx context.Context, query string, args ...interface{}) (int64, error)
	ExecuteTx(ctx context.Context, query string, args ...interface{}) (int64, SaveFunc, error)
}

// Rows is an interface which is implemented by sql.Rows. This makes testing
// of the database layer easier.
type Rows interface {
	Err() error
	Next() bool
	Scan(dest ...interface{}) error
}

// Row is an interface which is implemented by sql.Row. This makes testing
// of the database layer easier.
type Row interface {
	Scan(dest ...interface{}) error
}

type database struct {
	sql *sql.DB
}

// New returns a new instance of Database.
func New(connectionString string) (Database, error) {
	db, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		return nil, err
	}

	return &database{
		sql: db,
	}, nil
}

func (db *database) Multiple(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	stmt, err := db.sql.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.QueryContext(ctx, args...)
}

func (db *database) Single(ctx context.Context, query string, args ...interface{}) (Row, error) {
	stmt, err := db.sql.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	return stmt.QueryRowContext(ctx, args...), nil
}

func (db *database) Execute(ctx context.Context, query string, args ...interface{}) (int64, error) {
	stmt, err := db.sql.PrepareContext(ctx, query)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return -1, err
	}

	rowsAffected, _ := res.RowsAffected()
	return rowsAffected, nil
}

func (db *database) ExecuteTx(ctx context.Context, query string, args ...interface{}) (int64, SaveFunc, error) {
	tx, err := db.sql.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadUncommitted,
	})
	if err != nil {
		return -1, nil, err
	}

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return -1, nil, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return -1, nil, err
	}

	rowsAffected, _ := res.RowsAffected()
	return rowsAffected, save(tx), nil
}

func save(tx *sql.Tx) SaveFunc {
	return func(save bool) {
		if save {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}
}
