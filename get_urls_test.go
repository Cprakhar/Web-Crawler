package main

import (
	"reflect"
	"testing"
)

func TestGetURLsFromHTML(t *testing.T) {
	cases := []struct {
		name string
		inputURL string
		inputBody string
		expected []string
	}{
		{
			name:     "absolute and relative URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `
		<html>
			<body>
				<a href="/path/one">
					<span>Boot.dev</span>
				</a>
				<a href="https://other.com/path/one">
					<span>Boot.dev</span>
				</a>
			</body>
		</html>
		`,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},
	}

	for i, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			urls, err := getURLsFromHTML(c.inputBody, c.inputURL)
			if err != nil {
				t.Errorf("Test %d - '%s' FAIL: unexpected error: %v", i, c.name, err)
				return
			}
			if !reflect.DeepEqual(urls, c.expected) {
				t.Errorf("Test %d - '%s' FAIL: expected: %v, actual: %v", i, c.name, c.expected, urls)
			}
		})
	}
}