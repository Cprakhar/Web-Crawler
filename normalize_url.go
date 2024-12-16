package main

import (
	"net/url"
	"strings"
)

func normalizeURL(urlString string) (string, error) {
	parsedUrl, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}

	normalizedURL := parsedUrl.Host + parsedUrl.Path
	normalizedURL = strings.ToLower(normalizedURL)
	normalizedURL = strings.TrimSuffix(normalizedURL, "/")
	return normalizedURL, nil
}

/*
https://blog.boot.dev/path/
https://blog.boot.dev/path
http://blog.boot.dev/path/
http://blog.boot.dev/path
*/