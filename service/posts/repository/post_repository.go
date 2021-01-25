package repository

import (
	"context"
	"database/sql"

	"github.com/reecerussell/open-social/service/posts/dto"
	"github.com/reecerussell/open-social/service/posts/model"

	// MSSQL driver
	_ "github.com/denisenkom/go-mssqldb"
)

// PostRepository is a high level interface used to manipulate post data.
type PostRepository interface {
	Create(ctx context.Context, p *model.Post) error
	GetFeed(ctx context.Context, userReferenceID string) ([]*dto.FeedItem, error)
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

func (r *postRepository) GetFeed(ctx context.Context, userReferenceID string) ([]*dto.FeedItem, error) {
	db, err := sql.Open("sqlserver", r.url)
	if err != nil {
		return nil, err
	}

	const query = `;WITH [Feed] AS (
		SELECT 
			[P].[ReferenceId], 
			[P].[Caption], 
			[P].[Posted],
			[U].[Username],
			[dbo].GetPostLikes([P].[Id]) AS [Likes],
			[dbo].HasUserLikedPost([P].[Id], [U].[Id]) AS [HasUserLiked]
		FROM [Posts] AS [P]
		INNER JOIN [Users] AS [U] ON [U].[Id] = [P].[UserId]
		WHERE [U].[ReferenceId] = @userReference
		UNION
		SELECT
			[P].[ReferenceId],
			[P].[Caption],
			[P].[Posted],
			[U].[Username],
			[dbo].GetPostLikes([P].[Id]) AS [Likes],
			[dbo].HasUserLikedPost([P].[Id], [U].[Id]) AS [HasUserLiked]
		FROM [UserFollowers] AS [UF]
		INNER JOIN [Posts] AS [P] ON [P].[UserId] = [UF].[FollowerId]
		INNER JOIN [Users] AS [U] ON [U].[Id] = [P].[UserId]
		WHERE [U].[ReferenceId] = @userReference)
		
		SELECT CAST([ReferenceId] AS CHAR(36)), [Caption], [Posted], [Username], [Likes], [HasUserLiked] FROM [Feed]
		ORDER BY [Posted] DESC`

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, sql.Named("userReference", userReferenceID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feed []*dto.FeedItem

	for rows.Next() {
		var item dto.FeedItem
		err := rows.Scan(&item.ID, &item.Caption, &item.Posted, &item.Username, &item.Likes, &item.HasUserLiked)
		if err != nil {
			return nil, err
		}

		feed = append(feed, &item)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return feed, nil
}
