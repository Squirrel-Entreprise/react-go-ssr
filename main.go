package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

var (
	baseURL      string
	outPutDir    string
	waitPage     time.Duration
	thread       int
	visitedPages = make(map[string]bool)
	mu           = &sync.Mutex{}
)

func init() {
	flag.StringVar(&baseURL, "h", "", "base URL e.g. http://localhost:3000")
	flag.StringVar(&outPutDir, "o", "outhtml", "output directory")
	flag.DurationVar(&waitPage, "w", 2*time.Second, "wait on page")
	flag.IntVar(&thread, "c", 2, "concurrent thread")

	flag.Parse()
}

func main() {

	sem := make(chan bool, thread)

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

	var wg sync.WaitGroup
	wg.Add(1)
	visitPage(sem, &wg, browser, u.String(), "index", visitedPages)

	wg.Wait()
}

func visitPage(sem chan bool, wg *sync.WaitGroup, browser *rod.Browser, pageURL, path string, visitedPages map[string]bool) {
	defer wg.Done()
	defer func() { <-sem }()

	mu.Lock()
	visited := visitedPages[pageURL]
	if !visited {
		visitedPages[pageURL] = true
	}
	mu.Unlock()

	if visited {
		return
	}

	sem <- true

	visitedPages[pageURL] = true

	log.Println("Visiting", pageURL)

	var page *rod.Page
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
				log.Printf("failed to visit page %v: %v", pageURL, r)
			}
		}()
		page = browser.MustPage(pageURL)
	}()

	if page == nil {
		return
	}

	if err := page.WaitLoad(); err != nil {
		log.Printf("failed to load page %v: %v", pageURL, err)
		return
	}

	if err := page.WaitIdle(30 * time.Second); err != nil {
		log.Printf("failed to wait for page idle %v: %v", pageURL, err)
		return
	}

	time.Sleep(waitPage)

	content := page.MustHTML()

	saveFile(outPutDir+"/"+path, content)

	links := page.MustElements("a")
	for _, link := range links {
		href, err := link.Attribute("href")
		if err != nil {
			fmt.Printf("failed to get href attribute: %v\n", err)
			continue
		}

		if href == nil || !strings.HasPrefix(*href, "/") {
			continue
		}

		u, err := url.Parse(*href)
		if err != nil {
			continue
		}

		wg.Add(1)
		go func() {
			visitPage(sem, wg, browser, baseURL+*href, strings.Trim(u.Path, "/"), visitedPages)
			<-sem
		}()
	}
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
