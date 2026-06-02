package observability

import "os"

// GetEnv mengambil nilai environment variable dengan fallback ke defaultValue.
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
