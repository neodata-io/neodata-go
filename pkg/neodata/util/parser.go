package util

import "strconv"

// parseQueryParam converts a query parameter to an integer or returns a default value.
func ParseQueryParam(param string, defaultValue int) int {
	if value, err := strconv.Atoi(param); err == nil && value > 0 {
		return value
	}
	return defaultValue
}
