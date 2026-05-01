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
	testCases := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"abcde", 3, "abc"},
		{"abcde", 5, "abcde"},
		{"abcde", 10, "abcde"},
		{"abcde", 0, ""},
		{"abcde", -1, ""},
		{"", 3, ""},
		{"你好世界", 2, "你好"},
		{"你好世界", 4, "你好世界"},
		{"你好世界", 10, "你好世界"},
		{"a你b好c", 3, "a你b"},
	}
	for _, tc := range testCases {
		got := stringx.Truncate(tc.input, tc.maxLen)
		if got != tc.expected {
			t.Fatalf("Truncate(%q, %d): expected %q, got %q", tc.input, tc.maxLen, tc.expected, got)
		}
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
	testCases := []struct {
		input    string
		expected string
	}{
		{"snake_to_camel", "snakeToCamel"},
		{"_test_case", "testCase"},
		{"__test", "test"},
		{"_test", "test"},
		{"___a_b", "aB"},
		{"test", "test"},
		{"", ""},
		{"_", ""},
		{"__", ""},
	}
	for _, tc := range testCases {
		camel := stringx.SnakeToCamelLower(tc.input)
		if camel != tc.expected {
			t.Fatalf("expected %s, got %s", tc.expected, camel)
		}
	}
}
