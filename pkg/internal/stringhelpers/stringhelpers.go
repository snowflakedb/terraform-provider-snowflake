package stringhelpers

import "strings"

// SnakeCaseToCamel converts a snake_case string to a camelCase string.
func SnakeCaseToCamel(snake string) string {
	var suffix string
	if strings.HasSuffix(snake, "_") {
		suffix = "_"
		snake = strings.TrimSuffix(snake, "_")
	}
	snake = strings.ToLower(snake)
	parts := strings.Split(snake, "_")
	for idx, p := range parts {
		if p == "" {
			continue
		}
		parts[idx] = strings.ToUpper(p[:1]) + p[1:]
	}
	return strings.Join(parts, "") + suffix
}
