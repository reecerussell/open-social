package password

// Options contains a set of critrea that a
// password has to hit, to be considered valid.
type Options struct {
	// RequiredLength is the minimum length a password has to be.
	RequiredLength int `json:"requiredLength"`

	// RequireUppercase is a flag which demands at least one
	// character to be uppercase.
	RequireUppercase bool `json:"requireUppercase"`

	// RequireLowercase is a flag which demands at least one
	// character to be lowercase.
	RequireLowercase bool `json:"requireLowercase"`

	// RequireNonAlphanumeric is a flag which demands at least one
	// character to be non alphanumeric.
	RequireNonAlphanumeric bool `json:"requireNonAlphanumeric"`

	// RequireDigit is a flag which demands at least one
	// character to be a digit.
	RequireDigit bool `json:"requireDigit"`

	// RequiredUniqueChars is an integer value, that determines
	// how many unique characters are required in a password.
	RequiredUniqueChars int `json:"requiredUniqueChars"`
}

// HashOptions contains the configuration used to hash and
// verify passwords.
type HashOptions struct {
	// IterationCount is the number of times the password will be hashed.
	IterationCount int `json:"iterationCount"`

	// SaltSize is the bit count of the salt.
	SaltSize int `json:"saltSize"`

	// KeySize is the bit count of the hashed password.
	KeySize int `json:"keySize"`

	// HashKey represents the hash algorithm used to hash the
	// passwords. This can either be set to 1 or 2.
	HashKey int `json:"hashKey"`
}
