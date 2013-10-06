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
	persist(CurrentWeek.Movies)
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
		fmt.Println(err)
	}
	loc, _ := time.LoadLocation("")
	currentWeek.Date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, loc)

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
		movies[i] = Movie{Title: title, Imdb: imdb, Rank: i + 1, Week: currentWeek.Date}
	})

	currentWeek.Movies = movies
	CurrentWeek = currentWeek
}

func persist(movies []Movie) {
	dbmap := GetDbMap()
	for _, movie := range movies {
		var existing []Movie
		_, err := dbmap.Select(&existing, "select * from movies where week = :week and title = :title", map[string]interface{}{
			"week":  movie.Week,
			"title": movie.Title,
		})
		if err != nil {
			fmt.Println(err)
		}
		if len(existing) > 0 {
			movie.Id = existing[0].Id
			fmt.Println("Updating movie " + movie.Title)
			_, err = dbmap.Update(&movie)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Inserting movie " + movie.Title)
			err = dbmap.Insert(&movie)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("New Id: %d\n", movie.Id)
		}

		for _, service := range movie.Services {
			var existing []Service
			_, err := dbmap.Select(&existing, "select * from services where movie_id = :movie_id and name = :name", map[string]interface{}{
				"movie_id": movie.Id,
				"name":     service.Name,
			})
			if err != nil {
				fmt.Println(err)
			}
			if len(existing) > 0 {
				service.Id = existing[0].Id
				service.MovieId = existing[0].MovieId
				fmt.Println("Updating service " + service.Name)
				_, err = dbmap.Update(&service)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				fmt.Println("Inserting service " + service.Name)
				service.MovieId = movie.Id
				err = dbmap.Insert(&service)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
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
	movie.ServicesMap = make(map[string]bool)

	services.Each(func(i int, s *goquery.Selection) {
		if class, exists := s.Attr("class"); exists {
			name := strings.Split(class, " ")[0]
			available := s.HasClass("available")
			service := Service{Name: name, Available: available}
			movie.Services = append(movie.Services, service)
			movie.ServicesMap[name] = available
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
