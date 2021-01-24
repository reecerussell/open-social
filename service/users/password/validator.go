package password

import (
	"errors"
	"fmt"
)

// Validator is an interface used to validate passwords.
type Validator interface {
	Validate(password string) error
}

type validator struct {
	opt *Options
}

// New returns a new instance of Validator, with the given options.
func New(opt *Options) Validator {
	return &validator{opt: opt}
}

func (v *validator) Validate(password string) error {
	l := len(password)
	if l < 1 {
		return errors.New("password is required")
	}

	if l < v.opt.RequiredLength {
		return fmt.Errorf("password must be at least %d characters long", v.opt.RequiredLength)
	}

	var (
		hasNonAlphanumeric bool
		hasDigit           bool
		hasLower           bool
		hasUpper           bool
		uniqueChars        []byte
	)

	for i := 0; i < l; i++ {
		c := byte(password[i])

		if !hasNonAlphanumeric && !isLetterOrDigit(c) {
			hasNonAlphanumeric = true
		}

		if !hasDigit && isDigit(c) {
			hasDigit = true
		}

		if !hasLower && isLower(c) {
			hasLower = true
		}

		if !hasUpper && isUpper(c) {
			hasUpper = true
		}

		d := true

		for _, dc := range uniqueChars {
			if dc == c {
				d = false
			}
		}

		if d {
			uniqueChars = append(uniqueChars, c)
		}
	}

	if v.opt.RequireNonAlphanumeric && !hasNonAlphanumeric {
		return errors.New("password requires an non-alphanumeric character")
	}

	if v.opt.RequireDigit && !hasDigit {
		return errors.New("password requires a digit")
	}

	if v.opt.RequireLowercase && !hasLower {
		return errors.New("password requires a lowercase letter")
	}

	if v.opt.RequireUppercase && !hasUpper {
		return errors.New("password requires an uppercase letter")
	}

	if v.opt.RequiredUniqueChars >= 1 && len(uniqueChars) < v.opt.RequiredUniqueChars {
		return fmt.Errorf("password requires at least %d unique characters", v.opt.RequiredUniqueChars)
	}

	return nil
}

// isDigit returns a flag indicating whether the supplied character
// is a digit - true if the character is a digit, otherwise false.
func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// isLower returns a flag indicating whether the supplied character is
// a lower case ASCII letter - true if the character is a lower case
// ASCII letter, otherwise false.
func isLower(c byte) bool {
	return c >= 'a' && c <= 'z'
}

// isUpper returns a flag indicating whether the supplied character is
// an upper case ASCII letter - true if the character is an upper case ASCII
// letter, otherwise false.
func isUpper(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

// isLetterOrDigit returns a flag indicating whether the supplied character is
// an ASCII letter or digit - true if the character is an ASCII letter or digit,
// otherwise false.
func isLetterOrDigit(c byte) bool {
	return isUpper(c) || isLower(c) || isDigit(c)
}
