package main

import (
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	urls := make([]string, 0)
	reader := strings.NewReader(htmlBody)
	doc, err := html.Parse(reader)
	if err != nil {
		return []string{}, err
	}
	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					parsedURL, err := url.Parse(attr.Val)
					if err != nil {
						return urls, err
					}
					if parsedURL.Host == "" {
						urls = append(urls, rawBaseURL + attr.Val)
						continue
					}
					urls = append(urls, attr.Val)
				}
			}
		}
	}
	return urls, nil
}