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
	"sync"
	"time"
)

var (
	cache           []item
	cacheExpiration time.Time
	cacheMutex      sync.Mutex
)

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host string
}

type storyCache struct {
	numStories int
	cache      []item
	useA       bool
	expiration time.Time
	duration   time.Duration
	mutext     sync.Mutex
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

func getStories(ids []int) []item {
	type result struct {
		index int
		item  item
		err   error
	}

	resultsCh := make(chan result)
	defer close(resultsCh)

	for i := 0; i < len(ids); i++ {
		go func(index, id int) {
			var client hn.Client

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

	for i := 0; i < len(ids); i++ {
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

	return stories
}

func getTopStories(numStories int) ([]item, error) {
	var client hn.Client

	ids, err := client.TopItems()
	if err != nil {
		return nil, errors.New("failed to load top stories")
	}

	var stories []item
	currentPosition := 0

	for len(stories) < numStories {
		remainingStories := (numStories - len(stories)) * 5 / 4

		stories = append(stories, getStories(ids[currentPosition:currentPosition+remainingStories])...)
		currentPosition += remainingStories
	}

	return stories[:numStories], nil
}

func (sc *storyCache) stories() ([]item, error) {
	sc.mutext.Lock()
	defer sc.mutext.Unlock()

	if time.Now().Sub(sc.expiration) < 0 {
		return sc.cache, nil
	}

	stories, err := getTopStories(sc.numStories)
	if err != nil {
		return nil, err
	}

	sc.expiration = time.Now().Add(5 * time.Minute)
	sc.cache = stories

	return sc.cache, nil
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	sc := &storyCache{
		numStories: numStories,
		duration:   5 * time.Minute,
	}

	go func() {
		ticker := time.NewTicker(4 * time.Minute)
		for {
			temp := &storyCache{
				numStories: numStories,
				duration:   sc.duration * 2,
			}

			_, _ = temp.stories()

			sc.mutext.Lock()
			sc.cache = temp.cache
			sc.expiration = temp.expiration
			sc.mutext.Unlock()

			//fmt.Printf("[*] cache has %d stories == %+v\n", len(sc.cache), sc.duration)
			<-ticker.C
		}
	}()

	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		stories, err := sc.stories()
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
