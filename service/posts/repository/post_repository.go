package repository

import (
	"context"
	"database/sql"

	"github.com/reecerussell/open-social/service/posts/model"

	// MSSQL driver
	_ "github.com/denisenkom/go-mssqldb"
)

// PostRepository is a high level interface used to manipulate post data.
type PostRepository interface {
	Create(ctx context.Context, p *model.Post) error
}

type postRepository struct {
	url string
}

// NewPostRepository returns a new instance of PostRepository.
func NewPostRepository(url string) PostRepository {
	return &postRepository{url: url}
}

func (r *postRepository) Create(ctx context.Context, p *model.Post) error {
	db, err := sql.Open("sqlserver", r.url)
	if err != nil {
		return err
	}

	const query = `INSERT INTO [Posts] ([ReferenceId],[UserId],[Posted],[Caption])
					VALUES (NEWID(), @userId, @posted, @caption)
				SELECT [Id], CAST([ReferenceId] AS CHAR(36)) FROM [Posts] WHERE [Id] = SCOPE_IDENTITY()`

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	post := p.Dao()
	row := stmt.QueryRowContext(ctx,
		sql.Named("userId", post.UserID),
		sql.Named("posted", post.Posted),
		sql.Named("caption", post.Caption))

	// Read the post's ids
	err = row.Scan(&post.ID, &post.ReferenceID)
	if err != nil {
		return err
	}

	// Set the post's ids
	p.SetID(post.ID)
	p.SetReferenceID(post.ReferenceID)

	return nil
}
