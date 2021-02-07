package repository

import (
	"context"
	"database/sql"

	"github.com/reecerussell/open-social/database"
)

// LikeRepository is a high level interface used to manipulate persisted post like data.
type LikeRepository interface {
	Create(ctx context.Context, postID int, userReferenceID string) error
	Delete(ctx context.Context, postID int, userReferenceID string) error
}

type likeRepository struct {
	db database.Database
}

// NewLikeRepository returns a new instance of LikeRepository.
func NewLikeRepository(db database.Database) LikeRepository {
	return &likeRepository{db: db}
}

func (r *likeRepository) Create(ctx context.Context, postID int, userReferenceID string) error {
	const query = `INSERT INTO [PostLikes] ([PostId],[UserId])
					SELECT @postId, [Id] FROM [Users]
					WHERE [ReferenceId] = @userReferenceId;`

	_, err := r.db.Execute(ctx, query, sql.Named("postId", postID), sql.Named("userReferenceId", userReferenceID))
	if err != nil {
		return err
	}

	return nil
}

func (r *likeRepository) Delete(ctx context.Context, postID int, userReferenceID string) error {
	const query = `DELETE [L] FROM [PostLikes] AS [L]
						INNER JOIN [Users] AS [U] ON [U].[Id] = [L].[UserId]
					WHERE [L].[PostId] = @postId AND [U].[ReferenceId] = @userReferenceId;`

	_, err := r.db.Execute(ctx, query, sql.Named("postId", postID), sql.Named("userReferenceId", userReferenceID))
	if err != nil {
		return err
	}

	return nil
}
