package internal

import "os"

// GetEnvOrElse looks for some environment variable
// if not found, a default value is returned
func GetEnvOrElse(name string, defauktValue string) string {
	value, exists := os.LookupEnv(name)
	if exists {
		return value
	}
	return defauktValue
}