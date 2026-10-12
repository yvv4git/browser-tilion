// Package cdp provides a helper to connect to a remote Tilion browser over
// the Chrome DevTools Protocol for the examples.
package cdp

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-rod/rod"
)

// Connect connects to the Tilion container. The addr can be either the plain
// HTTP endpoint (e.g. http://localhost:9222) or the full WebSocket URL
// (e.g. ws://localhost:9222/devtools/browser/<id>).
func Connect(addr string) *rod.Browser {
	wsURL := addr
	if !strings.HasPrefix(wsURL, "ws://") && !strings.HasPrefix(wsURL, "wss://") {
		var err error
		wsURL, err = resolveWSURL(addr)
		if err != nil {
			panic(err)
		}
	}
	return rod.New().ControlURL(wsURL).MustConnect()
}

// Open creates a new tab, navigates it to url and waits until the page is no
// longer loading. It does not wait for the window "load" event: on some sites
// (go.dev, google.com, golang.org) a hanging third-party request prevents that
// event from ever firing, which would block forever.
func Open(b *rod.Browser, url string) *rod.Page {
	page := b.MustPage("")
	page.MustNavigate(url)
	WaitReady(page)
	return page
}

// WaitReady waits until the page document leaves the "loading" state, bounded
// by a timeout, and ignores the error so a slow third-party request cannot
// hang the program.
func WaitReady(page *rod.Page) {
	ready := rod.Eval(`() => document.readyState !== "loading"`)
	_ = page.Timeout(20 * time.Second).Wait(ready)
}

func resolveWSURL(addr string) (string, error) {
	endpoint := strings.TrimSuffix(addr, "/") + "/json/version"

	resp, err := http.Get(endpoint)
	if err != nil {
		return "", fmt.Errorf("fetch %s: %w", endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("fetch %s: %s", endpoint, resp.Status)
	}

	var data struct {
		WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("parse %s: %w", endpoint, err)
	}
	if data.WebSocketDebuggerURL == "" {
		return "", errors.New("no webSocketDebuggerUrl in " + endpoint)
	}
	return data.WebSocketDebuggerURL, nil
}
