package util

import "os"

// ReadEnv is equivalent to os.Getenv, but with a default value.
func ReadEnv(name, defaultValue string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}

	return defaultValue
}
