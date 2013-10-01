package models

import (
	"encoding/xml"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/mjibson/goread/goapp/rss"
	"html"
	"log"
	"net/http"
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
	summarizeWeek(&CurrentWeek)
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
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		log.Println(err)
	}
	rows := doc.Find("tbody > tr")
	movies := make([]Movie, rows.Length())

	rows.Each(func(i int, s *goquery.Selection) {
		title := s.Find("td").Eq(2).Find("a").Text()
		imdbUrl, _ := s.Find("a[href^=\"http://www.imdb.com/title\"]").First().Attr("href")
		imdb := strings.Split(imdbUrl, "/")[4]
		movies[i] = Movie{Title: title, Imdb: imdb, Rank: i + 1, Week: currentWeek.Name, Services: make(map[string]bool)}
	})

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

	movie.Streaming = doc.Find("#streaming > ul > li.available").Length()
	movie.Rental = doc.Find("#rental > ul > li.available").Length()
	movie.Purchase = doc.Find("#purchase > ul > li.available").Length()
	movie.DVD = doc.Find("#dvd > ul > li.available").Length()
	movie.All = movie.Streaming + movie.Rental + movie.Purchase // We're not counting DVDs here

	services := doc.Find("#streaming, #rental, #rental, #dvd").Find("ul > li").Not(".none-avail")

	services.Each(func(i int, s *goquery.Selection) {
		if class, exists := s.Attr("class"); exists {
			name := strings.Split(class, " ")[0]
			movie.Services[name] = s.HasClass("available")
		}
	})

	done <- true
}

func summarizeWeek(week *Week) {
	var streaming, rental, purchase, dvd, all int
	for m := range week.Movies {
		if week.Movies[m].Streaming > 0 {
			streaming += 1
		}
		if week.Movies[m].Rental > 0 {
			rental += 1
		}
		if week.Movies[m].Purchase > 0 {
			purchase += 1
		}
		if week.Movies[m].DVD > 0 {
			dvd += 1
		}
		if week.Movies[m].All > 0 {
			all += 1
		}
	}
	week.Streaming = streaming
	week.Rental = rental
	week.Purchase = purchase
	week.DVD = dvd
	week.All = all
}
