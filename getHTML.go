package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func getHTML(rawURL string) (string, error) {
	r, err := http.Get(rawURL)
	if err != nil {
		return ``, err
	}
	defer r.Body.Close()
	if r.StatusCode >= 400 {
		return ``, errors.New("couldn't respond")
	}
	if !strings.Contains(r.Header.Get("Content-Type"), "text/html") {
		return ``, errors.New("unsupported content type")
	}
	html, err := io.ReadAll(r.Body)
	if err != nil {
		return ``, err
	}
	return string(html), nil
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer func(){
		<-cfg.concurrencyControl
		cfg.wg.Done()
		}()
		
	if cfg.maxPagesVisited() {
		return
	}
	parsedCurrentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		return
	}
	if cfg.baseURL.Host != parsedCurrentURL.Host {
		return
	}
	normalizedCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		return
	}
	isFirst := cfg.addPageVisit(normalizedCurrentURL)
	if !isFirst {
		return
	}
	fmt.Printf("crawling: %s\n", rawCurrentURL)
	html, err := getHTML(rawCurrentURL)
	if err != nil {
		return
	}
	urls, err := getURLsFromHTML(html, cfg.baseURL.String())
	if err != nil {
		return
	}
	for _, url := range urls {
		cfg.wg.Add(1)
		go cfg.crawlPage(url)
	}
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	if _, ok := cfg.pages[normalizedURL]; ok {
		cfg.pages[normalizedURL]++
		return false
	}
	cfg.pages[normalizedURL] = 1
	return true
}

func (cfg *config) maxPagesVisited() bool {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	return len(cfg.pages) > cfg.maxPages
}