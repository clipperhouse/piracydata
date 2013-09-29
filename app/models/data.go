package models

import (
	"encoding/xml"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/mjibson/goread/goapp/rss"
	"html"
	"net/http"
	"regexp"
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
	url := "http://www.canistream.it/external/imdb/" + movie.Imdb
	doc, err := goquery.NewDocument(url)
	if err != nil {
		panic(err)
	}

	movie.Stream = doc.Find("#streaming > ul > li.available").Length() > 0
	movie.Rent = doc.Find("#rental > ul > li.available").Length() > 0
	movie.Buy = doc.Find("#purchase > ul > li.available").Length() > 0 || doc.Find("#dvd > ul > li.available").Length() > 0
	done <- true
}
