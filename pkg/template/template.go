package template

import (
	"strings"
)

// Replace replaces string with value that matched given keys suffix with $
func Replace(val string, keyValues map[string]string) string {
	var result = ""

	for k, v := range keyValues {
		result = strings.ReplaceAll(val, "$"+k, v)
	}

	return result
}
