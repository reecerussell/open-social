package password

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	val := New(&Options{
		RequiredLength:         8,
		RequireUppercase:       true,
		RequireLowercase:       true,
		RequireNonAlphanumeric: true,
		RequireDigit:           true,
		RequiredUniqueChars:    5,
	})

	err := val.Validate("Password_123")
	assert.NoError(t, err)
}

func TestValidate_GivenInvalidPassword_ReturnsError(t *testing.T) {
	t.Run("Requires Non-Empty Password", func(t *testing.T) {
		val := New(&Options{})
		err := val.Validate("")
		assert.Equal(t, "password is required", err.Error())
	})

	t.Run("Requires 8 Characters", func(t *testing.T) {
		opt := &Options{RequiredLength: 8}
		val := New(opt)
		err := val.Validate("pass123")
		exp := fmt.Sprintf("password must be at least %d characters long", opt.RequiredLength)
		assert.Equal(t, exp, err.Error())
	})

	t.Run("Requires Non Alphanumeric", func(t *testing.T) {
		opt := &Options{RequireNonAlphanumeric: true}
		val := New(opt)
		err := val.Validate("password123")
		assert.Equal(t, "password requires an non-alphanumeric character", err.Error())
	})

	t.Run("Requires Digit", func(t *testing.T) {
		opt := &Options{RequireDigit: true}
		val := New(opt)
		err := val.Validate("password")
		assert.Equal(t, "password requires a digit", err.Error())
	})

	t.Run("Requires Lowercase", func(t *testing.T) {
		opt := &Options{RequireLowercase: true}
		val := New(opt)
		err := val.Validate("PASSWORD")
		assert.Equal(t, "password requires a lowercase letter", err.Error())
	})

	t.Run("Requires Lowercase", func(t *testing.T) {
		opt := &Options{RequireUppercase: true}
		val := New(opt)
		err := val.Validate("password")
		assert.Equal(t, "password requires an uppercase letter", err.Error())
	})

	t.Run("Requires 4 Unique Characters", func(t *testing.T) {
		opt := &Options{RequiredUniqueChars: 4}
		val := New(opt)
		err := val.Validate("aaaa")
		exp := fmt.Sprintf("password requires at least %d unique characters", opt.RequiredUniqueChars)
		assert.Equal(t, exp, err.Error())
	})
}
