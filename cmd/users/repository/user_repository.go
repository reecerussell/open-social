package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/reecerussell/open-social/cmd/users/dao"
	"github.com/reecerussell/open-social/cmd/users/model"

	// MSSQL driver
	_ "github.com/denisenkom/go-mssqldb"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository interface {
	Create(ctx context.Context, u *model.User) error
	DoesUsernameExist(ctx context.Context, username string, excludeRefID *string) (bool, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)

	// GetUserByReference gets a user's model by referenceID,
	// for the user with the reference userReferenceID.
	GetUserByReference(ctx context.Context, referenceID, userReferenceID string) (*model.User, error)
	GetIDByReference(ctx context.Context, referenceID string) (*int, error)
}

type userRepository struct {
	url string
}

func NewUserRepository(url string) UserRepository {
	return &userRepository{url: url}
}

func (r *userRepository) Create(ctx context.Context, u *model.User) error {
	db, err := sql.Open("sqlserver", r.url)
	if err != nil {
		return err
	}

	const query = `INSERT INTO [Users] ([ReferenceId],[Username],[PasswordHash])
					VALUES (NEWID(), @username, @passwordHash)
				SELECT [Id], CAST([ReferenceId] AS CHAR(36)) FROM [Users] WHERE [Id] = SCOPE_IDENTITY()`

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	user := u.Dao()
	row := stmt.QueryRowContext(ctx,
		sql.Named("username", user.Username),
		sql.Named("passwordHash", user.PasswordHash))

	// Read the user's ids
	err = row.Scan(&user.ID, &user.ReferenceID)
	if err != nil {
		return err
	}

	// Set the user's ids
	u.SetID(user.ID)
	u.SetReferenceID(user.ReferenceID)

	return nil
}

func (r *userRepository) DoesUsernameExist(ctx context.Context, username string, excludeRefID *string) (bool, error) {
	db, err := sql.Open("sqlserver", r.url)
	if err != nil {
		return false, err
	}

	query := "SELECT COUNT(*) FROM [Users] WHERE [Username] = @username"
	args := []interface{}{sql.Named("username", username)}

	if excludeRefID != nil {
		query += " [ReferenceId] != @referenceId"
		args = append(args, sql.Named("referenceId", *excludeRefID))
	}

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	var count int64
	err = stmt.QueryRowContext(ctx, args...).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	db, err := sql.Open("sqlserver", r.url)
	if err != nil {
		return nil, err
	}

	const query = `SELECT [Id], CAST([ReferenceId] AS CHAR(36)), [Username], [PasswordHash]
					FROM [Users] WHERE [Username] = @username`

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var user dao.User
	err = stmt.QueryRowContext(ctx, sql.Named("username", username)).Scan(
		&user.ID,
		&user.ReferenceID,
		&user.Username,
		&user.PasswordHash,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return model.NewUserFromDao(&user), nil
}

func (r *userRepository) GetIDByReference(ctx context.Context, referenceID string) (*int, error) {
	db, err := sql.Open("sqlserver", r.url)
	if err != nil {
		return nil, err
	}

	const query = `SELECT [Id] FROM [Users] WHERE [ReferenceId] = @referenceId;`

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRowContext(ctx, sql.Named("referenceID", referenceID)).Scan(
		&id,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &id, nil
}

func (r *userRepository) GetUserByReference(ctx context.Context, referenceID, userReferenceID string) (*model.User, error) {
	db, err := sql.Open("sqlserver", r.url)
	if err != nil {
		return nil, err
	}

	const query = `SELECT
		[Id],
		CAST([ReferenceId] AS CHAR(36)),
		[Username],
		[PasswordHash],
		CASE (SELECT COUNT([UserId]) FROM [UserFollowers] AS [UF]
				INNER JOIN [Users] AS [U] ON [U].[Id] = [UF].[FollowerId]
				WHERE [UF].[UserId] = 8 AND [U].[ReferenceId] = @userReferenceId)
			WHEN 1 THEN CAST(1 AS BIT)
			ELSE CAST(0 AS BIT)
		END AS [IsFollowing]
		FROM [Users]
		WHERE [ReferenceId] = @referenceId;`

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var user dao.User
	err = stmt.QueryRowContext(ctx,
		sql.Named("referenceId", referenceID),
		sql.Named("userReferenceId", userReferenceID)).Scan(
		&user.ID,
		&user.ReferenceID,
		&user.Username,
		&user.PasswordHash,
		&user.IsFollowing,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return model.NewUserFromDao(&user), nil
}
