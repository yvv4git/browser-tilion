// Command rutube reproduces the crawler flow on a page (by default
// https://rutube.ru, where the crawler fails) and prints diagnostics that
// show where the Tilion page hangs: readyState timeline, load/stable waits,
// pending network requests, JS exceptions, HTML size and a screenshot.
//
// The container must be running first (make up).
//
//	go run ./rutube -addr http://localhost:9222
//	go run ./rutube -addr http://localhost:9222 -url https://quotes.toscrape.com
package main

import (
	"flag"
	"fmt"
	"log"
	"slices"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/rod/lib/utils"

	"github.com/yvv4git/browser-tilion/examples/internal/cdp"
)

const (
	stableWindow         = 500 * time.Millisecond
	readyStatePollingGap = time.Second
)

type pendingRequest struct {
	url     string
	started time.Time
}

type networkTracker struct {
	mu        sync.Mutex
	sent      int
	responded int
	finished  int
	failed    int
	pending   map[proto.NetworkRequestID]pendingRequest
}

func newNetworkTracker() *networkTracker {
	return &networkTracker{pending: make(map[proto.NetworkRequestID]pendingRequest)}
}

func (tracker *networkTracker) onRequest(event *proto.NetworkRequestWillBeSent) {
	tracker.mu.Lock()
	defer tracker.mu.Unlock()

	tracker.sent++
	tracker.pending[event.RequestID] = pendingRequest{url: event.Request.URL, started: time.Now()}
}

func (tracker *networkTracker) onResponse(_ *proto.NetworkResponseReceived) {
	tracker.mu.Lock()
	defer tracker.mu.Unlock()

	tracker.responded++
}

func (tracker *networkTracker) onFinished(event *proto.NetworkLoadingFinished) {
	tracker.mu.Lock()
	defer tracker.mu.Unlock()

	tracker.finished++
	delete(tracker.pending, event.RequestID)
}

func (tracker *networkTracker) onFailed(event *proto.NetworkLoadingFailed) {
	tracker.mu.Lock()
	defer tracker.mu.Unlock()

	tracker.failed++
	delete(tracker.pending, event.RequestID)
}

func (tracker *networkTracker) onException(event *proto.RuntimeExceptionThrown) {
	details := event.ExceptionDetails
	if details == nil {
		return
	}

	fmt.Printf("page exception: %s (%s)\n", details.Text, details.URL)
}

func (tracker *networkTracker) report() {
	tracker.mu.Lock()
	defer tracker.mu.Unlock()

	fmt.Printf("network: sent=%d responded=%d finished=%d failed=%d pending=%d\n",
		tracker.sent, tracker.responded, tracker.finished, tracker.failed, len(tracker.pending))

	if len(tracker.pending) == 0 {
		return
	}

	pending := make([]pendingRequest, 0, len(tracker.pending))
	for _, request := range tracker.pending {
		pending = append(pending, request)
	}
	slices.SortFunc(pending, func(a, b pendingRequest) int {
		return a.started.Compare(b.started)
	})

	fmt.Printf("still pending: %d\n", len(pending))
	for _, request := range pending {
		age := time.Since(request.started).Round(time.Second)
		fmt.Printf("  [%s ago] %s\n", age, request.url)
	}
}

func main() {
	addr := flag.String("addr", "http://localhost:9222", "CDP endpoint of the Tilion container")
	target := flag.String("url", "https://rutube.ru", "page to load")
	timeout := flag.Duration("timeout", 30*time.Second, "budget for load and stable waits")
	output := flag.String("output", "rutube.png", "screenshot output file")
	flag.Parse()

	browser := cdp.Connect(*addr)

	page := browser.MustPage("")
	defer func() { _ = page.Close() }()

	tracker := newNetworkTracker()
	waitEvents := page.EachEvent(
		tracker.onRequest,
		tracker.onResponse,
		tracker.onFinished,
		tracker.onFailed,
		tracker.onException,
	)
	go waitEvents()

	start := time.Now()
	if err := page.Navigate(*target); err != nil {
		log.Fatalf("navigate %s: %v", *target, err)
	}

	watchReadyState(page, start, *timeout)
	waitLoad(page, start, *timeout)
	waitStable(page, start, *timeout)
	printInfo(page)
	tracker.report()
	saveScreenshot(page, *output)
}

func watchReadyState(page *rod.Page, start time.Time, timeout time.Duration) {
	deadline := time.Now().Add(timeout)
	lastState := ""

	for time.Now().Before(deadline) {
		state := readyState(page)
		if state != lastState {
			fmt.Printf("t=%-9s readyState=%s\n", time.Since(start).Round(time.Second), state)
			lastState = state
		}
		if state != "loading" {
			return
		}

		time.Sleep(readyStatePollingGap)
	}
}

func readyState(page *rod.Page) string {
	result, err := page.Eval(`() => document.readyState`)
	if err != nil {
		return "error: " + err.Error()
	}

	return result.Value.Str()
}

func waitLoad(page *rod.Page, start time.Time, timeout time.Duration) {
	if err := page.Timeout(timeout).WaitLoad(); err != nil {
		fmt.Printf("wait load:   timeout after %s: %v\n", timeout, err)
		return
	}

	fmt.Printf("wait load:   ok in %s\n", time.Since(start).Round(time.Millisecond))
}

func waitStable(page *rod.Page, start time.Time, timeout time.Duration) {
	if err := page.Timeout(timeout).WaitStable(stableWindow); err != nil {
		fmt.Printf("wait stable: timeout after %s: %v\n", timeout, err)
		return
	}

	fmt.Printf("wait stable: ok in %s\n", time.Since(start).Round(time.Millisecond))
}

func printInfo(page *rod.Page) {
	info, err := page.Info()
	if err != nil {
		log.Printf("info: %v", err)
		return
	}
	fmt.Printf("title: %s\n", info.Title)
	fmt.Printf("url:   %s\n", info.URL)

	html, err := page.HTML()
	if err != nil {
		log.Printf("html: %v", err)
		return
	}
	fmt.Printf("html:  %d bytes\n", len(html))
}

func saveScreenshot(page *rod.Page, output string) {
	data, err := page.Screenshot(false, nil)
	if err != nil {
		log.Printf("screenshot: %v", err)
		return
	}
	if err := utils.OutputFile(output, data); err != nil {
		log.Printf("save %s: %v", output, err)
		return
	}
	fmt.Printf("screenshot saved: %s\n", output)
}
