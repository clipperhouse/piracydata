package app

import (
	"github.com/robfig/revel"
	"github.com/robfig/revel/modules/jobs/app/jobs"
	"piracydata/app/models"
	"time"
)

var loadCurrentWeek = jobs.Func(models.LoadCurrentWeek)
var fetch = jobs.Func(models.FetchAll)

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		// revel.SessionFilter,           // Restore and write the session cookie.
		// revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,  // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,        // Resolve the requested language
		revel.InterceptorFilter, // Run interceptors around the action.
		revel.ActionInvoker,     // Invoke the action.
	}

	revel.OnAppStart(func() {
		jobs.Now(loadCurrentWeek)
		jobs.Now(fetch)
		jobs.Every(1*time.Hour, loadCurrentWeek)
		jobs.Every(1*time.Hour, fetch)
	})
}
