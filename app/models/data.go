package models

import (
	gohtml "code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"encoding/xml"
	"fmt"
	"github.com/mjibson/goread/goapp/rss"
	"html"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

var CurrentWeek Week

var lock = sync.Mutex{}
var dataComplete = false
var onComplete = make(chan bool, 1)

func FetchAll() {
	fmt.Println("Starting FetchAll")
	getCurrentMovies()
	getCurrentAvailability()
	onComplete <- true
}

func AwaitData() {
	lock.Lock() // avoid pile-on
	for !dataComplete {
		dataComplete = <-onComplete
	}
	lock.Unlock()
}

var url string = "http://torrentfreak.com/category/dvdrip/feed/"
var titleRegex = regexp.MustCompile("<a href=\"http://www.(rottentomatoes.com|pnop.com|filmtied.com)/(.+)\">(.+)</a>")
var imdbRegex = regexp.MustCompile("<a href=\"http://www.imdb.com/title/([t\\d]+)[\\?/]+\">[\\d\\.\\?]+</a>")

const layout = "Jan 2, 2006"

func getCurrentMovies() {
	resp, _ := http.Get(url)
	decoder := xml.NewDecoder(resp.Body)
	feed := rss.Rss{}
	decoder.Decode(&feed)

	firstItem := feed.Items[0]
	currentWeek := Week{}
	date, err := time.Parse(time.RFC1123Z, firstItem.PubDate)
	if err != nil {
		currentWeek.Name = fmt.Sprintf("Week ending %s", firstItem.PubDate)
	} else {
		currentWeek.Name = fmt.Sprintf("Week ending %s", date.Format(layout))
	}

	content := html.UnescapeString(firstItem.Content)
	titles := titleRegex.FindAllStringSubmatch(content, -1)
	imdbs := imdbRegex.FindAllStringSubmatch(content, -1)

	movies := make([]Movie, len(titles))

	for j, _ := range titles {
		movies[j] = Movie{Title: titles[j][3], Imdb: imdbs[j][1], Rank: j + 1, Week: currentWeek.Name}
	}

	currentWeek.Movies = movies
	CurrentWeek = currentWeek
}

func getCurrentAvailability() {
	done := make(chan bool, 1)
	for m := range CurrentWeek.Movies {
		go getAvailability(&CurrentWeek.Movies[m], done)
	}
	for _ = range CurrentWeek.Movies {
		<-done
	}
}

func getAvailability(movie *Movie, done chan bool) {
	doc := scrape(movie)
	process(movie, doc)
	done <- true
}

func scrape(movie *Movie) *gohtml.Node {
	response, err := http.Get("http://www.canistream.it/external/imdb/" + movie.Imdb)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	doc, err := gohtml.Parse(response.Body)
	if err != nil {
		panic(err)
	}
	return doc
}

func process(movie *Movie, n *gohtml.Node) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.DataAtom == atom.Div && c.Attr[0].Key == "id" {
			if c.Attr[0].Val == "streaming" {
				movie.Stream = includesAvailableClass(c)
				continue
			}
			if c.Attr[0].Val == "rental" {
				movie.Rent = includesAvailableClass(c)
				continue
			}
			if c.Attr[0].Val == "purchase" || c.Attr[0].Val == "dvd" {
				movie.Buy = includesAvailableClass(c)
				continue
			}
		}
		process(movie, c)
		movie.Any = movie.Stream || movie.Rent || movie.Buy
	}
}

func includesAvailableClass(n *gohtml.Node) bool {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.DataAtom == atom.Ul {
			for l := c.FirstChild; l != nil; l = l.NextSibling {
				if l.DataAtom == atom.Li && l.Attr[0].Key == "class" && l.Attr[0].Val != "none-avail" {
					data := strings.Split(l.Attr[0].Val, " ")
					if data[1] == "available" {
						return true
					}
				}
			}
		}
	}
	return false
}
