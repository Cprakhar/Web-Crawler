package main

import "testing"

// https://blog.boot.dev/path/
// https://blog.boot.dev/path
// http://blog.boot.dev/path/
// http://blog.boot.dev/path

func TestNormalizeURL(t *testing.T) {
	cases := []struct {
		input string
		expected string
	}{
		{
			input: "https://blog.boot.dev/path/",
			expected: "blog.boot.dev/path",
		},
		{
			input: "https://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
		{
			input: "http://blog.boot.dev/path/",
			expected: "blog.boot.dev/path",
		},
		{
			input: "http://blog.boot.dev/path",
			expected: "blog.boot.dev/path",
		},
	}

	for i, c := range cases {
		t.Run("remove scheme", func (t *testing.T) {
			normalizedUrl, err := normalizeURL(c.input)
			if err != nil {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, "remove scheme", err)
				return
			}
			if normalizedUrl != c.expected {
				t.Errorf("Test %v - %s FAIL: expected URL: %v, actual: %v", i, "remove scheme", c.expected, normalizedUrl)
			}
		})
		
	}
}