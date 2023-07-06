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
	baseURL      string
	outPutDir    string
	waitPage     time.Duration
	visitedPages = make(map[string]bool)
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

	visitPage(browser, u.String(), "index", visitedPages)

}

func visitPage(browser *rod.Browser, pageURL, path string, visitedPages map[string]bool) {

	if ok := visitedPages[pageURL]; ok {
		return
	}

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
		page = browser.MustPage(pageURL).MustWaitLoad().MustWaitIdle()
	}()

	if page == nil {
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

		visitPage(browser, baseURL+*href, strings.Trim(u.Path, "/"), visitedPages)
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
