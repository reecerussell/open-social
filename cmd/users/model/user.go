package model

import (
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strings"

	hashpkg "github.com/reecerussell/adaptive-password-hasher"

	"github.com/reecerussell/open-social/cmd/users/dao"
	"github.com/reecerussell/open-social/cmd/users/password"
)

const (
	minUsernameLength = 3
	maxUsernameLength = 20
	usernameRegex     = "^[a-z0-9-_.]+$"
)

// Common user errors
var (
	ErrInvalidPassword = errors.New("password is invalid")
)

// User is a domain model representing a user.
type User struct {
	id           int
	referenceID  string
	username     string
	passwordHash string
}

// NewUser constructs a new user domain model.
// An error is returned if the given data is invalid.
func NewUser(username, password string, val password.Validator, hasher hashpkg.Hasher) (*User, error) {
	user := new(User)

	err := user.UpdateUsername(username)
	if err != nil {
		return nil, err
	}

	err = user.SetPassword(password, val, hasher)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// NewUserFromDao returns a new instance of User, populated with
// data from the data access object. This should only be used
// by the repository.
func NewUserFromDao(user *dao.User) *User {
	return &User{
		id:           user.ID,
		referenceID:  user.ReferenceID,
		username:     user.Username,
		passwordHash: user.PasswordHash,
	}
}

// ReferenceID returns the user's reference id.
func (u *User) ReferenceID() string {
	return u.referenceID
}

// Username returns the user's username.
func (u *User) Username() string {
	return u.username
}

// Dao returns a data access object for the user.
func (u *User) Dao() *dao.User {
	return &dao.User{
		ID:           u.id,
		ReferenceID:  u.referenceID,
		Username:     u.username,
		PasswordHash: u.passwordHash,
	}
}

// UpdateUsername updates the user's username.
func (u *User) UpdateUsername(username string) error {
	if username == "" {
		return errors.New("username is a required field")
	}

	username = strings.ToLower(username)

	l := len(username)
	if l < minUsernameLength {
		return fmt.Errorf("username must be greater than %d characters long", minUsernameLength)
	}

	if l > maxUsernameLength {
		return fmt.Errorf("username cannot be greater than %d characters long", maxUsernameLength)
	}

	re := regexp.MustCompile(usernameRegex)
	if !re.MatchString(username) {
		return errors.New("username must only contain alphanumerics, hyphens, underscores and periods")
	}

	u.username = username

	return nil
}

// SetPassword validates the user's password, then sets the new password hash.
func (u *User) SetPassword(pwd string, val password.Validator, hasher hashpkg.Hasher) error {
	err := val.Validate(pwd)
	if err != nil {
		return err
	}

	bytes := hasher.Hash([]byte(pwd))
	u.passwordHash = base64.StdEncoding.EncodeToString(bytes)

	return nil
}

// VerifyPassword uses the hasher to verify the given password
// against the user's password hash.
func (u *User) VerifyPassword(pwd string, hasher hashpkg.Hasher) error {
	hashBytes, _ := base64.StdEncoding.DecodeString(u.passwordHash)
	ok := hasher.Verify([]byte(pwd), hashBytes)
	if !ok {
		return ErrInvalidPassword
	}

	return nil
}

// SetID sets the id of the user. This should only
// be used in the repository when creating a user.
func (u *User) SetID(id int) {
	u.id = id
}

// SetReferenceID sets the reference id of the user. This should only
// be used in the repository when creating a user.
func (u *User) SetReferenceID(referenceID string) {
	u.referenceID = referenceID
}
