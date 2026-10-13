// Command connect connects to the Tilion container over CDP and lists the
// browser version and the currently open tabs.
//
// The container must be running first (make up).
//
//	go run ./connect -addr http://localhost:9222
package main

import (
	"flag"
	"fmt"

	"github.com/yvv4git/browser-tilion/examples/internal/cdp"
)

func main() {
	addr := flag.String("addr", "http://localhost:9222", "CDP endpoint of the Tilion container")
	flag.Parse()

	browser := cdp.Connect(*addr)

	version, err := browser.Version()
	if err != nil {
		panic(err)
	}
	fmt.Printf("browser: %s\n", version)

	pages := browser.MustPages()
	fmt.Printf("tabs: %d\n", len(pages))
	for i, p := range pages {
		info := p.MustInfo()
		fmt.Printf("  [%d] %s - %s\n", i+1, info.URL, info.Title)
	}
}
