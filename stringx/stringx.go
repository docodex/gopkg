package stringx

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Reverse returns the string reversed at the code-point (rune) level.
//
// Note: this does NOT operate on user-perceived characters (grapheme
// clusters). Inputs that include combining marks (e.g. "é" composed as
// "e" + U+0301), emoji ZWJ sequences, skin-tone modifiers, or regional
// indicator flag sequences will be corrupted by this function because
// the parts of each cluster are reversed individually. If you need
// grapheme-safe reversal, use a library that segments by grapheme
// cluster (UAX #29) first.
func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Truncate truncates the string to the specified max rune count.
// If the string is shorter than maxLen, it is returned as-is.
//
// Note: this operates at the code-point (rune) level, not the
// user-perceived character (grapheme cluster) level. Inputs containing
// combining marks (e.g. "é" composed as "e" + U+0301), emoji ZWJ
// sequences, skin-tone modifiers, or regional-indicator flag sequences
// may be cut in the middle of a cluster. For grapheme-safe truncation,
// use a library that segments by grapheme cluster (UAX #29) first.
func Truncate(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	// Ranging over a string yields (byteOffset, rune) pairs, so when we
	// have already accepted maxLen runes the current i is the byte
	// offset of the NEXT rune's first byte - i.e. exactly on a rune
	// boundary. Slicing s[:i] therefore never splits a multi-byte UTF-8
	// sequence.
	count := 0
	for i := range s {
		if count == maxLen {
			return s[:i]
		}
		count++
	}
	return s
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
	first := true
	for _, part := range parts {
		if part == "" {
			continue
		}
		if first {
			result.WriteString(part)
			first = false
		} else {
			r, size := utf8.DecodeRuneInString(part)
			result.WriteRune(unicode.ToUpper(r))
			result.WriteString(part[size:])
		}
	}
	return result.String()
}
