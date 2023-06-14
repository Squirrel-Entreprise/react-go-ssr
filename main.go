package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
)

var baseURL string

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <url>")
		os.Exit(1)
	}

	visitedPages := make(map[string]bool)

	website := os.Args[1]
	u, err := url.Parse(website)
	if err != nil {
		fmt.Printf("failed to parse URL: %v\n", err)
		os.Exit(1)
	}

	baseURL = website

	browser := rod.New().MustConnect()
	defer browser.MustClose()

	var wg sync.WaitGroup
	wg.Add(1)
	visitPage(&wg, browser, u.String(), "index", visitedPages)

	wg.Wait()
}

func visitPage(wg *sync.WaitGroup, browser *rod.Browser, pageURL, path string, visitedPages map[string]bool) {
	defer wg.Done()

	// Vérifier si la page a déjà été visitée
	if visitedPages[pageURL] {
		return
	}

	visitedPages[pageURL] = true

	log.Println("Visiting", pageURL)

	page := browser.MustPage(pageURL).MustWaitLoad()

	page = page.MustWaitIdle()

	time.Sleep(10 * time.Second)

	// Save the page content
	content := page.MustHTML()

	saveFile("outhtml/"+path, content)

	// Find all the internal links on the page
	links := page.MustElements("a")
	for _, link := range links {
		href, err := link.Attribute("href")
		if err != nil {
			fmt.Printf("failed to get href attribute: %v\n", err)
			continue
		}

		// Ignore if href is nil or an external link
		if href == nil || !strings.HasPrefix(*href, "/") {
			continue
		}

		u, err := url.Parse(*href)
		if err != nil {
			continue
		}

		wg.Add(1)
		go visitPage(wg, browser, baseURL+*href, strings.Trim(u.Path, "/"), visitedPages)
	}
}

func saveFile(path, content string) {
	path += ".html"
	os.MkdirAll(getDir(path), os.ModePerm)
	file, err := os.Create(path)
	if err != nil {
		fmt.Printf("failed to create file: %v\n", err)
		return
	}
	defer file.Close()

	file.WriteString(content)
}

func getDir(filePath string) string {
	if strings.Contains(filePath, "/") {
		return filePath[:strings.LastIndex(filePath, "/")]
	} else {
		return ""
	}
}
