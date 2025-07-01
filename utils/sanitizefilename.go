package utils

import "strings"

func SanitizeFilename(filename string) string {
	safeFilename := strings.ReplaceAll(filename, " ", "_")
	var b strings.Builder
	for _, r := range safeFilename {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '.' || r == '_' || r == '-' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
