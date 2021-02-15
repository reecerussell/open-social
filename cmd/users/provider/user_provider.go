package provider

import (
	"context"
	"database/sql"
	"errors"

	"github.com/reecerussell/open-social/cmd/users/dto"
	"github.com/reecerussell/open-social/database"
)

// Common errors.
var (
	ErrProfileNotFound = errors.New("profile not found")
)

// UserProvider is used to query user data, for read-only operations.
type UserProvider interface {
	GetProfile(ctx context.Context, username, userReferenceID string) (*dto.Profile, error)
	GetInfo(ctx context.Context, userReferenceID string) (*dto.Info, error)
}

type userProvider struct {
	db database.Database
}

// NewUserProvider returns a new instance of UserProvider.
func NewUserProvider(db database.Database) UserProvider {
	return &userProvider{db: db}
}

func (p *userProvider) GetProfile(ctx context.Context, username, userReferenceID string) (*dto.Profile, error) {
	const query = `;WITH [Followers] AS (
		SELECT [FollowerId] FROM [UserFollowers] AS [UF]
		INNER JOIN [Users] AS [U] ON [U].[Id] = [UF].[UserId]
		WHERE [U].[Username] = @username
	)
	
	SELECT
		CAST([U].[ReferenceId] AS CHAR(36)) AS [UserId],
		[U].[Username],
		CAST([M].[ReferenceId] AS CHAR(36)) AS [MediaId],
		[U].[Bio],
		(SELECT COUNT(*) FROM [Followers]) AS [FollowerCount],
		CASE (SELECT COUNT(*) FROM [Followers]
				WHERE [FollowerId] = [CU].[Id])
			WHEN 1 THEN CAST(1 AS BIT)
			ELSE CAST(0 AS BIT)
		END AS [IsFollowing],
		CASE [U].[Id] WHEN [CU].[Id] THEN CAST(1 AS BIT) ELSE CAST(0 AS BIT) END AS [IsOwner],
		(SELECT COUNT([Id]) FROM [Posts] WHERE [UserId] = [U].[Id]) AS [PostCount]
	FROM [Users] AS [U]
	LEFT JOIN [Media] AS [M] ON [M].[Id] = [U].[MediaId]
	INNER JOIN [Users] [CU] ON [CU].[ReferenceId] = @userReferenceId
	WHERE [U].[Username] = @username;`

	row, err := p.db.Single(ctx, query, sql.Named("username", username), sql.Named("userReferenceId", userReferenceID))
	if err != nil {
		return nil, err
	}

	var profile dto.Profile
	err = row.Scan(
		&profile.UserID,
		&profile.Username,
		&profile.MediaID,
		&profile.Bio,
		&profile.FollowerCount,
		&profile.IsFollowing,
		&profile.IsOwner,
		&profile.PostCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrProfileNotFound
		}

		return nil, err
	}

	return &profile, nil
}

func (p *userProvider) GetInfo(ctx context.Context, userReferenceID string) (*dto.Info, error) {
	const query = `SELECT 
		CAST([U].[ReferenceId] AS CHAR(36)) AS [Id],
		[U].[Username] AS [Username],
		CAST([M].[ReferenceId] AS CHAR(36)) AS [MediaId],
		(SELECT COUNT([UserId]) FROM [UserFollowers] WHERE [UserId] = [U].[Id]) AS [FollowerCount]
	FROM [Users] AS [U]
	LEFT JOIN [Media] AS [M] ON [M].[Id] = [U].[MediaId]
	WHERE [U].[ReferenceId] = @userReferenceId;`

	row, err := p.db.Single(ctx, query, sql.Named("userReferenceId", userReferenceID))
	if err != nil {
		return nil, err
	}

	var info dto.Info
	err = row.Scan(
		&info.ID,
		&info.Username,
		&info.MediaID,
		&info.FollowerCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrProfileNotFound
		}

		return nil, err
	}

	return &info, nil
}
