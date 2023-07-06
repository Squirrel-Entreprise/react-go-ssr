package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

var (
	baseURL   string
	outPutDir string
	waitPage  time.Duration
)

func init() {
	flag.StringVar(&baseURL, "h", "", "base URL e.g. http://localhost:3000")
	flag.StringVar(&outPutDir, "o", "outhtml", "output directory")
	flag.DurationVar(&waitPage, "w", 2*time.Second, "wait on page")

	flag.Parse()
}

func main() {

	u, err := url.Parse(baseURL)
	if err != nil {
		fmt.Printf("failed to parse URL: %v\n", err)
		os.Exit(1)
	}

	path, ok := launcher.LookPath()
	if !ok {
		fmt.Printf("failed to find chrome\n")
		os.Exit(1)
	}

	l, err := launcher.New().Bin(path).Launch()
	if err != nil {
		fmt.Printf("failed to launch chrome: %v\n", err)
		os.Exit(1)
	}

	browser := rod.New().ControlURL(l).MustConnect()
	defer browser.MustClose()

	// Use a map to keep track of visited pages.
	visitedPages := make(map[string]bool)

	// Start with a single page in the queue.
	queue := []string{u.String()}

	for len(queue) > 0 {
		// Take the next page from the queue.
		pageURL := queue[0]
		queue = queue[1:]

		// Skip this page if we've visited it already.
		if visitedPages[pageURL] {
			continue
		}

		// Mark this page as visited.
		visitedPages[pageURL] = true

		// Visit the page and get the links.
		newLinks, err := visitPage(browser, pageURL)
		if err != nil {
			log.Println(err)
			continue
		}

		// Add the new links to the queue.
		queue = append(queue, newLinks...)
	}
}

func visitPage(browser *rod.Browser, pageURL string) ([]string, error) {
	log.Println("Visiting", pageURL)

	var page *rod.Page
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
				log.Printf("failed to visit page %v: %v", pageURL, r)
			}
		}()
		page = browser.MustPage(pageURL).MustWaitLoad().MustWaitIdle()
	}()

	if page == nil {
		return nil, fmt.Errorf("failed to load page %v", pageURL)
	}

	defer page.MustClose()

	time.Sleep(waitPage)

	content := page.MustHTML()

	// Create a valid file path from the URL
	u, err := url.Parse(pageURL)
	if err != nil {
		return nil, err
	}
	path := strings.Trim(u.Path, "/")
	if path == "" {
		path = "index"
	}

	saveFile(outPutDir+"/"+path, content)

	links := page.MustElements("a")
	newLinks := []string{}
	for _, link := range links {
		href, err := link.Attribute("href")
		if err != nil {
			fmt.Printf("failed to get href attribute: %v\n", err)
			continue
		}

		if href == nil || !strings.HasPrefix(*href, "/") {
			continue
		}

		newLinks = append(newLinks, baseURL+*href)
	}

	return newLinks, nil
}

func saveFile(path, content string) {
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
