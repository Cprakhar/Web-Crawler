package main

import (
	"fmt"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
)

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

func main() {
	if len(os.Args[1:]) < 3 {
		fmt.Println("no website provided")
		os.Exit(1)
	}
	if len(os.Args[1:]) > 3 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	fmt.Printf("starting crawl of: %s\n", os.Args[1])
	pages := make(map[string]int, 0)
	rawbaseURL, err := url.Parse(os.Args[1])
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}
	maxConcurrency, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("Invalid argument for maxConcurrency: %s, error: %v", os.Args[2], err)
		os.Exit(1)
	}
	maxPages, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Printf("Invalid argument for maxPages: %s, error: %v", os.Args[3], err)
		os.Exit(1)
	}
	cfg := config{
		pages:              pages,
		baseURL:            rawbaseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
	}
	cfg.wg.Add(1)
	go cfg.crawlPage(os.Args[1])
	cfg.wg.Wait()
	printReport(cfg.pages, os.Args[1])
}

type Page struct {
	page string
	count int
}

func printReport(pages map[string]int, baseURL string) {

	fmt.Println("=============================")
	fmt.Printf("  REPORT for %s\n", baseURL)
	fmt.Println("=============================")
	pageSlice := make([]Page, 0)
	for page, cnt := range pages {
		pageSlice = append(pageSlice, Page{page: page, count: cnt})
	}	
	slices.SortFunc(pageSlice, func(a Page, b Page) int {
		if a.count == b.count {
			return strings.Compare(a.page, b.page)
		}
		if a.count > b.count {
			return -1
		}
		return 1
	})
	for _, page := range pageSlice {
		fmt.Printf("Found %d internal links to %s\n", page.count, page.page)
	}
}