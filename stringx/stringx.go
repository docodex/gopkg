package stringx

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Reverse returns the string reversed (Unicode-aware).
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Truncate truncates the string to the specified max rune count.
// If the string is shorter than maxLen, it is returned as-is.
func Truncate(s string, maxLen int) string {
	if maxLen < 0 {
		return ""
	}
	if utf8.RuneCountInString(s) <= maxLen {
		return s
	}
	runes := []rune(s)
	return string(runes[:maxLen])
}

// PadLeft pads the string on the left to the specified length with the given pad character.
func PadLeft(s string, length int, pad rune) string {
	n := length - utf8.RuneCountInString(s)
	if n <= 0 {
		return s
	}
	return strings.Repeat(string(pad), n) + s
}

// PadRight pads the string on the right on the specified length with the given pad character.
func PadRight(s string, length int, pad rune) string {
	n := length - utf8.RuneCountInString(s)
	if n <= 0 {
		return s
	}
	return s + strings.Repeat(string(pad), n)
}

// CamelToSnake converts a CamelCase or camelCase string to snake_case.
func CamelToSnake(s string) string {
	runes := []rune(s)
	var result []rune
	for i, r := range runes {
		if unicode.IsUpper(r) {
			if i > 0 {
				prev := runes[i-1]
				if unicode.IsLower(prev) || unicode.IsDigit(prev) {
					result = append(result, '_')
				} else if unicode.IsUpper(prev) && i+1 < len(runes) && unicode.IsLower(runes[i+1]) {
					result = append(result, '_')
				}
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

// SnakeToCamel converts a snake_case string to CamelCase (PascalCase).
func SnakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	var result strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}
		r, size := utf8.DecodeRuneInString(part)
		result.WriteRune(unicode.ToUpper(r))
		result.WriteString(part[size:])
	}
	return result.String()
}

// SnakeToCamelLower converts a snake_case string to camelCase (lower camelCase).
func SnakeToCamelLower(s string) string {
	parts := strings.Split(s, "_")
	var result strings.Builder
	for i, part := range parts {
		if part == "" {
			continue
		}
		if i == 0 {
			result.WriteString(part)
		} else {
			r, size := utf8.DecodeRuneInString(part)
			result.WriteRune(unicode.ToUpper(r))
			result.WriteString(part[size:])
		}
	}
	return result.String()
}
