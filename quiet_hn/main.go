package main

import (
	"flag"
	"fmt"
	"github.com/jwambugu/gophercises/quiet_hn/hn"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"runtime"
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
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var client hn.Client
		ids, err := client.TopItems()
		if err != nil {
			http.Error(w, "Failed to load top stories", http.StatusInternalServerError)
			return
		}
		var stories []item
		for _, id := range ids {
			hnItem, err := client.GetItem(id)
			if err != nil {
				continue
			}
			item := parseHNItem(hnItem)
			if isStoryLink(item) {
				stories = append(stories, item)
				if len(stories) >= numStories {
					break
				}
			}
		}
		data := templateData{
			Stories: stories,
			Time:    time.Now().Sub(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	})
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
