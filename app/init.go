package app

import (
	"github.com/darkhelmet/webutil"
	"github.com/robfig/revel"
	"github.com/robfig/revel/modules/jobs/app/jobs"
	"piracydata/app/models"
	"time"
)

var loadAllWeeks = jobs.Func(models.LoadAllWeeks)
var fetch = jobs.Func(models.FetchAll)
var Version string = time.Now().Format("200601021504")

func init() {
	revel.TemplateFuncs["only"] = func(a int) string {
		if a > 49 || a == 0 {
			return ""
		} else {
			return "only "
		}
	}
	
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
		enableGzip()
		jobs.Now(loadAllWeeks)
		jobs.Now(fetch)
		jobs.Every(1*time.Hour, loadAllWeeks)
		jobs.Every(1*time.Hour, fetch)
	})
}

func enableGzip() {
	handler := revel.Server.Handler
	handler = webutil.GzipHandler{handler}
	revel.Server.Handler = handler
}
