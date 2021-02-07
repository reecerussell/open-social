package provider

import (
	"context"
	"database/sql"
	"errors"

	"github.com/reecerussell/open-social/cmd/posts/dto"
	"github.com/reecerussell/open-social/database"
)

// Common errors
var (
	ErrPostNotFound = errors.New("post not found")
)

// PostProvider is used to read post data.
type PostProvider interface {
	Get(ctx context.Context, postReferenceID, userReferenceID string) (*dto.Post, error)
}

type postProvider struct {
	db database.Database
}

// NewPostProvider returns a new instance of PostProvider.
func NewPostProvider(db database.Database) PostProvider {
	return &postProvider{db: db}
}

func (p *postProvider) Get(ctx context.Context, postReferenceID, userReferenceID string) (*dto.Post, error) {
	const query = `;WITH [Likes] AS (
			SELECT [U].[ReferenceId] AS [UserReferenceId] FROM [PostLikes] AS [L]
				INNER JOIN [Posts] AS [P] ON [P].[Id] = [L].[PostId]
				INNER JOIN [Users] AS [U] ON [U].[Id] = [L].[UserId]
			WHERE [P].[ReferenceId] = @postReferenceId
		)
		
		SELECT
			[P].[ReferenceId] AS [Id],
			[M].[ReferenceId] AS [MediaId],
			[P].[Posted],
			[U].[Username],
			[P].[Caption],
			(SELECT COUNT(*) FROM [PostLikes] WHERE [PostId] = [P].[Id]) AS [LikeCount],
			CASE (SELECT COUNT(*) 
					FROM [Likes] WHERE [UserReferenceId] = @userReferenceId) 
				WHEN 1 THEN CAST(1 AS BIT) 
				ELSE CAST(0 AS BIT) 
			END AS [HasUserLiked]
		FROM [Posts] AS [P]
		INNER JOIN [Users] AS [U] ON [U].[Id] = [P].[UserId]
		LEFT JOIN [Media] AS [M] ON [M].[Id] = [P].[MediaId]
		WHERE [P].[ReferenceId] = @postReferenceId 
			AND [U].[ReferenceId] = @userReferenceId;`

	row, err := p.db.Single(ctx, query, sql.Named("postReferenceId", postReferenceID), sql.Named("userReferenceId", userReferenceID))
	if err != nil {
		return nil, err
	}

	var post dto.Post
	err = row.Scan(
		&post.ID,
		&post.MediaID,
		&post.Posted,
		&post.Username,
		&post.Caption,
		&post.Likes,
		&post.HasLiked,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			err = ErrPostNotFound
		}

		return nil, err
	}

	return &post, nil
}
