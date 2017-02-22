package main

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/newrelic/go-agent"
)

func startPirriWebApp() {
	routes := map[string]func(http.ResponseWriter, *http.Request){

		// charts and reporting
		"/stats/4": statsStationActivity,
		"/stats/2": statsActivityByDayOfWeek,

		// weather
		// TODO write a better algorithm for weather handling

		// station
		"/station/run": stationRunWeb,
		"/station/all": stationAllWeb,
		"/station":     stationGetWeb,

		// schedule
		"/schedule/all":    stationScheduleAllWeb,
		"/schedule/edit":   stationScheduleEditWeb,
		"/schedule/delete": stationScheduleDeleteWeb,

		// history
		"/history": historyAllWeb,

		// static
		"/static/": (func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, r.URL.Path[1:])
		}),

		// root
		"/": (func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "templates/index.html")
		}),
	}

	if SETTINGS.UseNewRelic {
		SETTINGS.NewRelicLicense = loadNewRelicKey(SETTINGS.NewRelicLicensePath)
		config := newrelic.NewConfig("PirriGo v"+VERSION, SETTINGS.NewRelicLicense)
		NRAPPMON, ERR := newrelic.NewApplication(config)
		fmt.Println("Using New Relic Monitoring Agent")
		if NRAPPMON == nil || ERR != nil {
			fmt.Println("Unable to load New Relic Agent using given configuration.")
		} else {
			for k, v := range routes {
				http.HandleFunc(newrelic.WrapHandleFunc(NRAPPMON, k, v))
			}
		}

	} else {
		for k, v := range routes {
			fmt.Println("Not using New Relic for", k)
			http.HandleFunc(k, v)
		}
	}

	// Host server
	panic(http.ListenAndServe(":"+SETTINGS.HttpPort, nil))
}

func logTraffic() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}