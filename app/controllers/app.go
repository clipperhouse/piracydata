package controllers

import (
	"github.com/robfig/revel"
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
	home = models.Home{CurrentWeek: models.CurrentWeek}
	return
}
