package stringx_test

import (
	"testing"

	"github.com/docodex/gopkg/stringx"
)

func TestReverse(t *testing.T) {
	reversed := stringx.Reverse("abc")
	if reversed != "cba" {
		t.Fatalf("expected %s, got %s", "cba", reversed)
	}
}

func TestTruncate(t *testing.T) {
	truncated := stringx.Truncate("abcde", 3)
	if truncated != "abc" {
		t.Fatalf("expected %s, got %s", "abc", truncated)
	}
}

func TestPadLeft(t *testing.T) {
	padded := stringx.PadLeft("abc", 6, 'd')
	if padded != "dddabc" {
		t.Fatalf("expected %s, got %s", "dddabc", padded)
	}
}

func TestPadRight(t *testing.T) {
	padded := stringx.PadRight("abc", 6, 'd')
	if padded != "abcddd" {
		t.Fatalf("expected %s, got %s", "abcddd", padded)
	}
}

func TestCamelToSnake(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"camelToSnake", "camel_to_snake"},
		{"CamelToSnake", "camel_to_snake"},
		{"HTMLParser", "html_parser"},
		{"getHTTPResponse", "get_http_response"},
		{"myURLParser", "my_url_parser"},
		{"ATest", "a_test"},
		{"simpleTest", "simple_test"},
		{"", ""},
	}
	for _, tc := range testCases {
		snake := stringx.CamelToSnake(tc.input)
		if snake != tc.expected {
			t.Fatalf("CamelToSnake(%q): expected %q, got %q", tc.input, tc.expected, snake)
		}
	}
}

func TestSnakeToCamel(t *testing.T) {
	camel := stringx.SnakeToCamel("snake_to_camel")
	if camel != "SnakeToCamel" {
		t.Fatalf("expected %s, got %s", "SnakeToCamel", camel)
	}
}

func TestSnakeToCamelLower(t *testing.T) {
	camel := stringx.SnakeToCamelLower("snake_to_camel")
	if camel != "snakeToCamel" {
		t.Fatalf("expected %s, got %s", "snakeToCamel", camel)
	}
}
