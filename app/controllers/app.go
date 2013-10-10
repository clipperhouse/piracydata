package controllers

import (
	"github.com/robfig/revel"
	"piracydata/app"
	"piracydata/app/models"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	c.Response.Out.Header().Set("Cache-Control", "public, max-age=600")
	home := getModel()
	return c.Render(home)
}

func (c App) Csv() revel.Result {
	c.Response.Out.Header().Set("Content-Disposition", "attachment; filename=piracydata.csv")
	c.Response.ContentType = "text/csv"
	home := getModel()
	return c.Render(home)
}

func getModel() (home models.Home) {
	home = models.Home{CurrentWeek: models.CurrentWeek, AllWeeks: models.Weeks, Stats: calculateStats(), AppVersion: app.Version}
	return
}

func calculateStats() models.Stats {
	var digital, rentStream, streaming, n int
	for _, w := range models.Weeks {
		for _, m := range w.Movies {
			if m.Streaming > 0 {
				streaming++
			}
			if m.Streaming > 0 || m.Rental > 0 {
				rentStream++
			}
			if m.Streaming > 0 || m.Rental > 0 || m.Purchase > 0 {
				digital++
			}
			n++
		}
	}
	nWeeks := len(models.Weeks)
	digital = int(100*digital/n)
	rentStream = int(100*rentStream/n)
	streaming = int(100*streaming/n)
	return models.Stats{Digital: digital, RentStream: rentStream, Streaming: streaming, NWeeks: nWeeks}
}