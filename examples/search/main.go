// Command search runs a Wikipedia search in the Tilion container: fills the
// search input, clicks search, reads the article title and saves the first
// article image.
//
// The container must be running first (make up).
//
//	go run ./search -addr http://localhost:9222 -q Earth
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/go-rod/rod/lib/utils"

	"github.com/yvv4git/browser-tilion/examples/internal/cdp"
)

func main() {
	addr := flag.String("addr", "http://localhost:9222", "CDP endpoint of the Tilion container")
	query := flag.String("q", "Earth", "search query")
	output := flag.String("output", "article.png", "article image output file")
	flag.Parse()

	browser := cdp.Connect(*addr)

	page := cdp.Open(browser, "https://www.wikipedia.org")
	defer func() { _ = page.Close() }()

	searchInput := page.MustElement("#searchInput")
	searchInput.MustInput(*query)
	page.MustElement("#search-form > fieldset > button").MustClick()
	cdp.WaitReady(page)

	fmt.Printf("article: %s\n", page.MustElement("#firstHeading").MustText())

	img, err := page.Element("#mw-content-text img")
	if err != nil {
		log.Printf("image: %v", err)
		return
	}
	data, err := img.Resource()
	if err != nil {
		log.Printf("resource: %v", err)
		return
	}
	if err := utils.OutputFile(*output, data); err != nil {
		log.Printf("save: %v", err)
		return
	}
	log.Printf("article image saved: %s", *output)
}