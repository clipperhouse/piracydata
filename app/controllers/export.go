package controllers

import (
	"github.com/robfig/revel"
	"piracydata/app/models"
)

type Export struct {
	*revel.Controller
}

func (c Export) Csv() revel.Result {
	models.AwaitData()
	week := models.CurrentWeek
	c.Response.ContentType = "text/csv"
	return c.Render(week)
}
