package model

import (
	"encoding/base64"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	hashermock "github.com/reecerussell/adaptive-password-hasher/mock"
	"github.com/stretchr/testify/assert"

	"github.com/reecerussell/open-social/cmd/users/dao"
	"github.com/reecerussell/open-social/cmd/users/mock"
)

func TestNewUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const testUsername = "John-doe"
	const testPassword = "Password123"

	mockValidator := mock.NewMockValidator(ctrl)
	mockValidator.EXPECT().Validate(testPassword).Return(nil)

	mockHasher := hashermock.NewMockHasher(ctrl)
	mockHasher.EXPECT().Hash([]byte(testPassword)).Return([]byte(testPassword))

	user, err := NewUser(testUsername, testPassword, mockValidator, mockHasher)
	assert.NoError(t, err)
	assert.Equal(t, "john-doe", user.username)
	assert.Equal(t, "UGFzc3dvcmQxMjM=", user.passwordHash)
}

func TestNewUser_GivenInvalidData_ReturnsError(t *testing.T) {
	t.Run("Given Empty Username", func(t *testing.T) {
		const testUsername = ""
		const testPassword = "Password123"

		user, err := NewUser(testUsername, testPassword, nil, nil)
		assert.Nil(t, user)
		assert.Equal(t, "username is a required field", err.Error())
	})

	t.Run("Given Invalid Password", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		const testUsername = "john-doe"
		const testPassword = "Password123"

		expectedErr := errors.New("an error occured")

		mockValidator := mock.NewMockValidator(ctrl)
		mockValidator.EXPECT().Validate(testPassword).Return(expectedErr)

		user, err := NewUser(testUsername, testPassword, mockValidator, nil)
		assert.Nil(t, user)
		assert.Equal(t, expectedErr, err)
	})
}

func TestUser_UpdateUsername(t *testing.T) {
	var user User
	err := user.UpdateUsername("John_Doe")
	assert.NoError(t, err)
	assert.Equal(t, "john_doe", user.username)
}

func TestUser_UpdateUsername_ReturnsError(t *testing.T) {
	t.Run("Given Empty Username", func(t *testing.T) {
		var user User
		err := user.UpdateUsername("")
		assert.Equal(t, "username is a required field", err.Error())
	})

	t.Run("Given Short Username", func(t *testing.T) {
		var user User
		err := user.UpdateUsername("ab")
		exp := fmt.Sprintf("username must be greater than %d characters long", minUsernameLength)
		assert.Equal(t, exp, err.Error())
	})

	t.Run("Given Short Username", func(t *testing.T) {
		var user User
		err := user.UpdateUsername("a-really-long-username")
		exp := fmt.Sprintf("username cannot be greater than %d characters long", maxUsernameLength)
		assert.Equal(t, exp, err.Error())
	})

	t.Run("Given Invalid Username", func(t *testing.T) {
		invalidUsernames := []string{
			"John Doe",
			"JaneDoe?",
		}
		exp := "username must only contain alphanumerics, hyphens, underscores and periods"

		for _, u := range invalidUsernames {
			var user User
			err := user.UpdateUsername(u)
			assert.Equal(t, exp, err.Error())
		}
	})
}

func TestUser_SetPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const testPassword = "Password123"
	const expectedHash = "UGFzc3dvcmQxMjM="

	mockValidator := mock.NewMockValidator(ctrl)
	mockValidator.EXPECT().Validate(testPassword).Return(nil)

	mockHasher := hashermock.NewMockHasher(ctrl)
	mockHasher.EXPECT().Hash([]byte(testPassword)).
		DoAndReturn(func(pwd []byte) []byte {
			bytes, _ := base64.StdEncoding.DecodeString(expectedHash)
			return bytes
		})

	var user User
	err := user.SetPassword(testPassword, mockValidator, mockHasher)
	assert.NoError(t, err)
	assert.Equal(t, expectedHash, user.passwordHash)
}

func TestUser_SetPassword_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const testPassword = "Password123"
	expectedErr := errors.New("an error occured")

	mockValidator := mock.NewMockValidator(ctrl)
	mockValidator.EXPECT().Validate(testPassword).Return(expectedErr)

	var user User
	err := user.SetPassword(testPassword, mockValidator, nil)
	assert.Equal(t, expectedErr, err)
}

func TestUser_VerifyPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const testPassword = "Password123"
	const testPasswordHash = "UGFzc3dvcmQxMjM="

	mockHasher := hashermock.NewMockHasher(ctrl)
	mockHasher.EXPECT().Verify([]byte(testPassword), gomock.Any()).Return(true)

	user := &User{passwordHash: testPasswordHash}
	err := user.VerifyPassword(testPassword, mockHasher)
	assert.NoError(t, err)
}

func TestUser_VerifyPassword_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	const testPassword = "Password123"
	const testPasswordHash = "UGFzc3dvcmQxMjM="

	mockHasher := hashermock.NewMockHasher(ctrl)
	mockHasher.EXPECT().Verify([]byte(testPassword), gomock.Any()).Return(false)

	user := &User{passwordHash: testPasswordHash}
	err := user.VerifyPassword(testPassword, mockHasher)
	assert.Equal(t, ErrInvalidPassword, err)
}

func TestNewUserFromDao(t *testing.T) {
	const (
		testID           = 23
		testReferenceID  = "kho3hewq"
		testUsername     = "testing"
		testPasswordHash = "UGFzc3dvcmQxMjM="
		testIsFollowing  = true
	)

	d := &dao.User{
		ID:           testID,
		ReferenceID:  testReferenceID,
		Username:     testUsername,
		PasswordHash: testPasswordHash,
		IsFollowing:  testIsFollowing,
	}

	user := NewUserFromDao(d)
	assert.Equal(t, testID, user.id)
	assert.Equal(t, testReferenceID, user.referenceID)
	assert.Equal(t, testUsername, user.username)
	assert.Equal(t, testPasswordHash, user.passwordHash)
	assert.Equal(t, testIsFollowing, user.isFollowing)
}

func TestUser_Dao(t *testing.T) {
	const (
		testID           = 23
		testReferenceID  = "kho3hewq"
		testUsername     = "testing"
		testPasswordHash = "UGFzc3dvcmQxMjM="
	)

	user := &User{
		id:           testID,
		referenceID:  testReferenceID,
		username:     testUsername,
		passwordHash: testPasswordHash,
	}

	d := user.Dao()
	assert.Equal(t, testID, d.ID)
	assert.Equal(t, testReferenceID, d.ReferenceID)
	assert.Equal(t, testUsername, d.Username)
	assert.Equal(t, testPasswordHash, d.PasswordHash)
}

func TestUser_Username(t *testing.T) {
	const testUsername = "testing"

	user := &User{username: testUsername}
	username := user.Username()

	assert.Equal(t, testUsername, username)
}

func TestUser_ReferenceID(t *testing.T) {
	const testReferenceID = "ioau97023"

	user := &User{referenceID: testReferenceID}
	refID := user.ReferenceID()

	assert.Equal(t, testReferenceID, refID)
}

func TestUser_ID(t *testing.T) {
	const testID = 1232

	user := &User{id: testID}
	id := user.ID()

	assert.Equal(t, testID, id)
}

func TestUser_SetID(t *testing.T) {
	const testID = 1

	var user User
	user.SetID(testID)

	assert.Equal(t, testID, user.id)
}

func TestUser_SetReferenceID(t *testing.T) {
	const testReferenceID = "239432"

	var user User
	user.SetReferenceID(testReferenceID)

	assert.Equal(t, testReferenceID, user.referenceID)
}

func TestUser_CanFollow(t *testing.T) {
	user := User{isFollowing: false}
	err := user.CanFollow()
	assert.NoError(t, err)
}

func TestUser_CanFollow_ReturnsError(t *testing.T) {
	t.Run("User Already Follows", func(t *testing.T) {
		user := User{isFollowing: true}
		err := user.CanFollow()
		assert.Equal(t, "user is already following this user", err.Error())
	})
}

func TestUser_CanUnfollow(t *testing.T) {
	user := User{isFollowing: true}
	err := user.CanUnfollow()
	assert.NoError(t, err)
}

func TestUser_CanUnfollow_ReturnsError(t *testing.T) {
	t.Run("User Already Follows", func(t *testing.T) {
		user := User{isFollowing: false}
		err := user.CanUnfollow()
		assert.Equal(t, "user is not following this user", err.Error())
	})
}
