package utilx

import "strings"

// Helper privado para buscar case-insensitive en un slice
func ContainsFold(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, item) {
			return true
		}
	}
	return false
}
