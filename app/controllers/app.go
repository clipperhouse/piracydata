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
	week := models.CurrentWeek
	return c.Render(week)
}
