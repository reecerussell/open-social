package repository

import (
	"context"
	"database/sql"

	"github.com/reecerussell/open-social/service/media/model"

	// MSSQL driver
	_ "github.com/denisenkom/go-mssqldb"
)

// MediaRepository is used to interface with the media data store.
type MediaRepository interface {
	Create(ctx context.Context, m *model.Media) (func(bool), error)
}

type mediaRepository struct {
	url string
}

// NewMediaRepository returns a new instance of MediaRepository.
func NewMediaRepository(url string) MediaRepository {
	return &mediaRepository{url: url}
}

func (r *mediaRepository) Create(ctx context.Context, m *model.Media) (func(bool), error) {
	db, err := sql.Open("sqlserver", r.url)
	if err != nil {
		return nil, err
	}

	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadUncommitted})
	if err != nil {
		return nil, err
	}

	const query = `INSERT INTO [Media] ([ReferenceId],[ContentType])
					VALUES (NEWID(), @contentType)
				SELECT [Id], CAST([ReferenceId] AS CHAR(36)) FROM [Media] WHERE [Id] = SCOPE_IDENTITY()`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	media := m.Dao()
	row := stmt.QueryRowContext(ctx,
		sql.Named("contentType", media.ContentType))

	// Read the media's ids
	err = row.Scan(&media.ID, &media.ReferenceID)
	if err != nil {
		return nil, err
	}

	// Set the media's ids
	m.SetID(media.ID)
	m.SetReferenceID(media.ReferenceID)

	return save(tx), nil
}

func save(tx *sql.Tx) func(bool) {
	return func(save bool) {
		var err error
		if save {
			err = tx.Commit()
		} else {
			err = tx.Rollback()
		}

		if err != nil {
			panic(err)
		}
	}
}
