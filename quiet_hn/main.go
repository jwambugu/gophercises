package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/jwambugu/gophercises/quiet_hn/hn"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host string
}

type templateData struct {
	Stories []item
	Time    time.Duration
}

func getAbsolutePath() string {
	_, b, _, _ := runtime.Caller(0)

	return filepath.Join(filepath.Dir(b))
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{
		Item: hnItem,
	}

	parsedURL, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(parsedURL.Hostname(), "www.")
	}

	return ret
}

func getTopStories(numStories int) ([]item, error) {
	var client hn.Client

	ids, err := client.TopItems()
	if err != nil {
		return nil, errors.New("failed to load top stories")
	}

	type result struct {
		index int
		item  item
		err   error
	}

	resultsCh := make(chan result)

	for i := 0; i < numStories; i++ {
		go func(index, id int) {
			hnItem, err := client.GetItem(id)
			if err != nil {
				resultsCh <- result{
					index: index,
					err:   err,
				}
			}

			resultsCh <- result{
				index: index,
				item:  parseHNItem(hnItem),
			}
		}(i, ids[i])
	}

	var results []result

	for i := 0; i < numStories; i++ {
		results = append(results, <-resultsCh)
	}

	// Sort the items by index asc
	sort.Slice(results, func(i, j int) bool {
		return results[i].index < results[j].index
	})

	var stories []item

	for _, r := range results {
		if r.err != nil {
			continue
		}

		if isStoryLink(r.item) {
			stories = append(stories, r.item)
		}
	}

	return stories, nil
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		stories, err := getTopStories(numStories)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := templateData{
			Stories: stories,
			Time:    time.Now().Sub(start),
		}

		if err := tpl.Execute(w, data); err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	// parse flags
	var port, numStories int

	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	indexPage := fmt.Sprintf("%s/index.html", getAbsolutePath())
	tpl := template.Must(template.ParseFiles(indexPage))

	http.HandleFunc("/", handler(numStories, tpl))

	// Start the server
	log.Printf("Server running on port :%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
