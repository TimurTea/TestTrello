package config

import (
	"os"
)

func GetOrDefault(varName, defaultValue string) string {
	val := os.Getenv(varName)
	if len(val) == 0 {
		return defaultValue
	}
	return val
}
