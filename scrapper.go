package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	// 1. Command Line Arguments
	urlPtr := flag.String("url", "", "Comma-separated list of URLs to scrape (e.g., https://example.com,https://google.com)")
	bravePtr := flag.Bool("brave", false, "Use Brave browser (default location)")
	execPathPtr := flag.String("exec-path", "", "Path to browser executable")
	flag.Parse()

	if *urlPtr == "" {
		fmt.Println("Usage: go run scrapper.go -url=<URL1,URL2> [-brave] [-exec-path=<path>]")
		os.Exit(1)
	}

	// Create output directories
	dirs := []string{"html", "screenshots", "url"}
	for _, d := range dirs {
		if _, err := os.Stat(d); os.IsNotExist(err) {
			err := os.Mkdir(d, 0755)
			if err != nil {
				log.Fatalf("Error creating directory '%s': %v", d, err)
			}
		}
	}

	targets := strings.Split(*urlPtr, ",")
	log.Printf("Starting scrape for %d targets: %v\n", len(targets), targets)

	// 2. Setup Chromedp Allocator options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	// Configure Executable Path
	if *execPathPtr != "" {
		opts = append(opts, chromedp.ExecPath(*execPathPtr))
		log.Printf("Using custom browser path: %s\n", *execPathPtr)
	} else if *bravePtr {
		// Common Brave paths on Windows
		bravePath := "C:\\Program Files\\BraveSoftware\\Brave-Browser\\Application\\brave.exe"
		if _, err := os.Stat(bravePath); os.IsNotExist(err) {
			// Try x86 path
			bravePath = "C:\\Program Files (x86)\\BraveSoftware\\Brave-Browser\\Application\\brave.exe"
		}
		opts = append(opts, chromedp.ExecPath(bravePath))
		log.Printf("Using Brave browser at: %s\n", bravePath)
	}

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelAlloc()

	var wg sync.WaitGroup

	for _, targetURL := range targets {
		targetURL = strings.TrimSpace(targetURL)
		if targetURL == "" {
			continue
		}
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			scrapeURL(allocCtx, u)
		}(targetURL)
	}

	wg.Wait()
	log.Println("All scrapes completed.")
}

func scrapeURL(allocCtx context.Context, targetURL string) {
	log.Printf("[%s] parsing URL...", targetURL)
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		log.Printf("Error parsing URL %s: %v", targetURL, err)
		return
	}
	host := parsedURL.Hostname()
	host = strings.ReplaceAll(host, ":", "_")

	// Create a new context for this tab
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Set a timeout for this operation
	ctx, cancel = context.WithTimeout(ctx, 45*time.Second)
	defer cancel()

	// Variables to hold results
	var htmlContent string
	var buf []byte
	var links []string

	log.Printf("[%s] Navigating...", host)
	err = chromedp.Run(ctx,
		chromedp.Navigate(targetURL),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.OuterHTML("html", &htmlContent),
		chromedp.FullScreenshot(&buf, 90),
		chromedp.Evaluate(`Array.from(document.querySelectorAll('a')).map(a => a.href)`, &links),
	)

	if err != nil {
		log.Printf("[%s] Error: Failed to scrape. Details: %v", host, err)
		return
	}

	// Save HTML
	htmlFilename := filepath.Join("html", fmt.Sprintf("%s_site_data.html", host))
	err = ioutil.WriteFile(htmlFilename, []byte(htmlContent), 0644)
	if err != nil {
		log.Printf("[%s] Error writing HTML: %v", host, err)
	} else {
		log.Printf("[%s] HTML content saved to '%s'", host, htmlFilename)
	}

	// Save Screenshot
	if len(buf) > 0 {
		screenshotFilename := filepath.Join("screenshots", fmt.Sprintf("%s_screenshot.png", host))
		err = ioutil.WriteFile(screenshotFilename, buf, 0644)
		if err != nil {
			log.Printf("[%s] Error writing screenshot: %v", host, err)
		} else {
			log.Printf("[%s] Screenshot saved to '%s'", host, screenshotFilename)
		}
	}

	// Save URLs
	if len(links) > 0 {
		urlFilename := filepath.Join("url", fmt.Sprintf("%s_urls.txt", host))
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Scraped URL: %s\n", targetURL))
		sb.WriteString(fmt.Sprintf("Timestamp: %s\n", time.Now().Format(time.RFC3339)))
		sb.WriteString("--- Extracted URLs ---\n")
		for _, link := range links {
			if link != "" {
				sb.WriteString(link + "\n")
			}
		}
		err = ioutil.WriteFile(urlFilename, []byte(sb.String()), 0644)
		if err != nil {
			log.Printf("[%s] Error writing extracted URLs: %v", host, err)
		} else {
			log.Printf("[%s] URLs saved to '%s' (%d urls)", host, urlFilename, len(links))
		}
	} else {
		log.Printf("[%s] No URLs found.", host)
	}
}
