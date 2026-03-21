package env

import "strings"

// sensitivePatterns are substrings that indicate a key holds sensitive data.
var sensitivePatterns = []string{
	"PASSWORD", "SECRET", "KEY", "TOKEN", "API_KEY", "PRIVATE", "CREDENTIAL",
}

// IsSensitive returns true if the key likely holds sensitive data.
func IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range sensitivePatterns {
		if strings.Contains(upper, p) {
			return true
		}
	}
	return false
}

// MaskValue returns "****" for sensitive keys, original value otherwise.
func MaskValue(key, value string) string {
	if IsSensitive(key) {
		return "****"
	}
	return value
}
