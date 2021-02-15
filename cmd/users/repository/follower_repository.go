package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/reecerussell/open-social/database"
)

// Common errors
var (
	ErrFollowerNotFound = errors.New("follower not found")
)

// FollowerRepository is used to manipulate and perform write operations
// on the user follower records.
type FollowerRepository interface {
	Create(ctx context.Context, userID int, followerReferenceID string) error
	Delete(ctx context.Context, userID int, followerReferenceID string) error
}

type followerRepository struct {
	db database.Database
}

// NewFollowerRepository returns a new instance of FollowerRepository.
func NewFollowerRepository(db database.Database) FollowerRepository {
	return &followerRepository{db: db}
}

func (r *followerRepository) Create(ctx context.Context, userID int, followerReferenceID string) error {
	const query = `INSERT INTO [UserFollowers] ([UserId], [FollowerId])
		SELECT @userId, [Id] FROM [Users] WHERE [ReferenceId] = @followerReferenceId;`

	rowsAffected, err := r.db.Execute(ctx, query, sql.Named("userId", userID), sql.Named("followerReferenceId", followerReferenceID))
	if err != nil {
		return err
	}

	if rowsAffected < 1 {
		return ErrFollowerNotFound
	}

	return nil
}

func (r *followerRepository) Delete(ctx context.Context, userID int, followerReferenceID string) error {
	const query = `DELETE [UF] FROM [UserFollowers] AS [UF] 
		INNER JOIN [Users] AS [F] ON [F].[Id] = [UF].[FollowerId]
		WHERE [UF].[UserId] = @userId AND [F].[ReferenceId] = @followerReferenceId;`

	_, err := r.db.Execute(ctx, query, sql.Named("userId", userID), sql.Named("followerReferenceId", followerReferenceID))
	if err != nil {
		return err
	}

	return nil
}
