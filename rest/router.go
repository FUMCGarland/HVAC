package rest

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/FUMCGarland/hvac"
	"github.com/FUMCGarland/hvac/log"

	"github.com/julienschmidt/httprouter"
)

func getServeMux(c *hvac.Config) *httprouter.Router {
	m := httprouter.New()
	m.HandleOPTIONS = true
	m.GlobalOPTIONS = http.HandlerFunc(headers)

	if _, err := os.Stat(c.HTTPStaticDir); err != nil {
		panic(err.Error())
	}
	m.ServeFiles("/static/*filepath", http.Dir(c.HTTPStaticDir))
	appDir := fmt.Sprintf("%s/_app", c.HTTPStaticDir)
	m.ServeFiles("/_app/*filepath", http.Dir(appDir))

	m.NotFound = http.HandlerFunc(notFound)

	// Add handlers for all the endpoints
	m.GET("/api/v1/system", authMW(getSystem, AuthLevelView))           // all devices in one shot
	m.PUT("/api/v1/system/mode", authMW(putSystemMode, AuthLevelAdmin)) // heating or cooling
	m.PUT("/api/v1/system/control", authMW(putControl, AuthLevelAdmin)) // manual, schedule, or temp-sensor

	// manual control
	m.PUT("/api/v1/pump/:id/target", authMW(putPump, AuthLevelControl))       // manually start/stop a pump
	m.PUT("/api/v1/blower/:id/target", authMW(putBlower, AuthLevelControl))   // manually start/stop a blower
	m.PUT("/api/v1/blower/:id/filter", authMW(resetFilter, AuthLevelControl)) // reset the filter time
	m.PUT("/api/v1/zone/:id/target", authMW(putZoneStart, AuthLevelControl))  // manually start/stop an entire zone

	// manual system scheduling
	m.GET("/api/v1/schedule", authMW(getSchedule, AuthLevelView))           // get entire schedule
	m.POST("/api/v1/schedule", authMW(postSchedule, AuthLevelControl))      // add a new entry
	m.PUT("/api/v1/sched/:id", authMW(putSchedule, AuthLevelControl))       // update an entry
	m.DELETE("/api/v1/sched/:id", authMW(deleteSchedule, AuthLevelControl)) // delete an entry

	// temp/occupancy based scheduling: requires sensors)
	m.PUT("/api/v1/zone/:id/temps", authMW(putZoneTemps, AuthLevelControl)) // set target temp range for zone

	m.GET("/api/v1/occupancy/:id", authMW(getOccupancySchedule, AuthLevelView))                     // get entire schedule
	m.POST("/api/v1/occupancy/recurring", authMW(postOccupancyRecurring, AuthLevelControl))         // add a new entry
	m.PUT("/api/v1/occupancy/recurring/:id", authMW(putRecurring, AuthLevelControl))                // update an occupancy-expected entry
	m.DELETE("/api/v1/occupancy/recurring/:id", authMW(deleteOccupancyRecurring, AuthLevelControl)) // update an occupancy-expected entry

	m.POST("/api/v1/occupancy/onetime", authMW(postOccupancyOneTime, AuthLevelControl))         // add a new entry
	m.PUT("/api/v1/occupancy/onetime/:id", authMW(putOneTime, AuthLevelControl))                // update an occupancy-expected entry
	m.DELETE("/api/v1/occupancy/onetime/:id", authMW(deleteOccupancyOneTime, AuthLevelControl)) // update an occupancy-expected entry

	return m
}

func headers(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin != "" {
		w.Header().Add("Access-Control-Allow-Origin", origin)
	} else {
		w.Header().Add("Access-Control-Allow-Origin", "*")
	}
	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS, HEAD, DELETE, PATCH")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Accept, If-Modified-Since, If-Match, If-None-Match, Authorization")

	w.Header().Add("Content-Type", jsonType)
}

func TODO(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.Error(w, "Forbidden", http.StatusForbidden)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	if r.URL.String() == "" || r.URL.String() == "/" {
		http.Redirect(w, r, "/static/index.html", http.StatusMovedPermanently)
		return
	}

	if strings.HasPrefix(r.URL.String(), "/static") {
		err := fmt.Errorf("404, not found")
		http.Error(w, jsonError(err), http.StatusNotFound)
		return
	}

	newLoc := fmt.Sprintf("/static%s", r.URL)
	log.Debug("not found, redirecting", "request", r.URL, "new", newLoc, "method", r.Method)
	http.Redirect(w, r, newLoc, http.StatusMovedPermanently)
}
