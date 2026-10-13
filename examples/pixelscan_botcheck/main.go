// Command pixelscan_botcheck opens the PixelScan bot-check page through a local
// HTTP proxy, waits for a configurable time, saves a screenshot into
// examples/.screenshots/ with a timestamped filename, then closes the page and
// the browser.
//
//	go run ./pixelscan_botcheck -proxy http://localhost:8138
//	go run ./pixelscan_botcheck -proxy http://localhost:8138 -wait 15s
package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

func main() {
	proxyURL := flag.String("proxy", "http://localhost:8138", "HTTP proxy URL")
	pageURL := flag.String("url", "https://pixelscan.net/bot-check", "page to load")
	wait := flag.Duration("wait", 10*time.Second, "time to wait before screenshot")
	output := flag.String("output", "", "output file (default examples/.screenshots/pixelscan_botcheck_<timestamp>.png)")
	flag.Parse()

	log.Printf("stage: launching browser via proxy %s", *proxyURL)

	l := launcher.New().
		Proxy(*proxyURL).
		Headless(true)
	controlURL := l.MustLaunch()

	log.Printf("stage: browser launched, control URL %s", controlURL)

	browser := rod.New().ControlURL(controlURL).MustConnect()
	defer func() {
		log.Printf("stage: closing browser")
		browser.MustClose()
	}()

	log.Printf("stage: opening page %s", *pageURL)
	page := browser.MustPage(*pageURL)
	defer func() {
		log.Printf("stage: closing page")
		_ = page.Close()
	}()

	log.Printf("stage: waiting for page load")
	page.MustWaitLoad()

	log.Printf("stage: waiting %s", *wait)
	time.Sleep(*wait)

	out := *output
	if out == "" {
		dir := screenshotsDir()
		if err := os.MkdirAll(dir, 0o755); err != nil {
			log.Fatalf("mkdir %s: %v", dir, err)
		}
		out = filepath.Join(dir, "pixelscan_botcheck_"+time.Now().Format("20060102_150405")+".png")
	}

	log.Printf("stage: taking screenshot -> %s", out)
	page.MustScreenshot(out)
	log.Printf("stage: saved %s (%d bytes)", out, fileSize(out))
}

func screenshotsDir() string {
	_, thisFile, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filepath.Dir(thisFile)), ".screenshots")
}

func fileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}
