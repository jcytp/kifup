package env

import (
	"os"
)

func Get(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func SecretKey() []byte {
	return []byte(Get("SECRET_KEY", ""))
}
