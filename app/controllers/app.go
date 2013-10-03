package controllers

import (
	"github.com/robfig/revel"
	"piracydata/app/models"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	models.AwaitData()
	c.Response.Out.Header().Set("Cache-Control", "public, max-age=600")
	week := models.CurrentWeek
	return c.Render(week)
}

func (c App) Csv() revel.Result {
	models.AwaitData()
	week := models.CurrentWeek
	c.Response.ContentType = "text/csv"
	return c.Render(week)
}
