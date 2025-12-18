package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

func main() {
	// 1. Command Line Arguments
	urlPtr := flag.String("url", "", "The URL to scrape (e.g., https://example.com)")
	bravePtr := flag.Bool("brave", false, "Use Brave browser (default location)")
	execPathPtr := flag.String("exec-path", "", "Path to browser executable")
	flag.Parse()

	if *urlPtr == "" {
		fmt.Println("Usage: go run scrapper.go -url=<URL> [-brave] [-exec-path=<path>]")
		os.Exit(1)
	}

	targetURL := *urlPtr
	log.Printf("Starting scrape for: %s\n", targetURL)

	// 2. Setup Chromedp
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

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Set a timeout for the entire operation
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Variables to hold results
	var htmlContent string
	var buf []byte
	var links []string

	// 3. Scraping Tasks
	log.Println("Connecting and navigating...")
	err := chromedp.Run(ctx,
		chromedp.Navigate(targetURL),
		// Wait for body to be visible to ensure page loaded
		chromedp.WaitVisible("body", chromedp.ByQuery),
		// Capture HTML
		chromedp.OuterHTML("html", &htmlContent),
		// Capture Screenshot (Full page)
		chromedp.FullScreenshot(&buf, 90),
		// Bonus: Extract URLs
		chromedp.Evaluate(`Array.from(document.querySelectorAll('a')).map(a => a.href)`, &links),
	)

	// Error Handling for Connection/Navigation
	if err != nil {
		log.Printf("Error: Failed to scrape %s. Details: %v\n", targetURL, err)
		// Check for specific timeout or context cancellation
		if err == context.DeadlineExceeded {
			log.Println("Reason: Operation timed out. The site might be slow or unreachable.")
		}
		os.Exit(1)
	}

	// 4. Save Data
	// Save HTML
	err = ioutil.WriteFile("site_data.html", []byte(htmlContent), 0644)
	if err != nil {
		log.Printf("Error writing HTML to file: %v\n", err)
	} else {
		log.Println("Success: HTML content saved to 'site_data.html'")
	}

	// Save Screenshot
	if len(buf) > 0 {
		err = ioutil.WriteFile("screenshot.png", buf, 0644)
		if err != nil {
			log.Printf("Error writing screenshot to file: %v\n", err)
		} else {
			log.Println("Success: Screenshot saved to 'screenshot.png'")
		}
	} else {
		log.Println("Warning: Screenshot buffer was empty.")
	}

	// Bonus: Print URLs
	if len(links) > 0 {
		fmt.Println("\n--- Extracted URLs ---")
		for i, link := range links {
			if link != "" {
				fmt.Printf("[%d] %s\n", i+1, link)
			}
		}
		fmt.Printf("Total URLs found: %d\n", len(links))
	} else {
		log.Println("No URLs found on the page.")
	}

	log.Println("Done.")
}
