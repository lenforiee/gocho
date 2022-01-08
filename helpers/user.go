package helpers

import (
	"strings"
)

func UsernameSafe(s string) string {
	return strings.Replace(strings.ToLower(s), " ", "_", -1)
}
