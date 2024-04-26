package rest

import (
	"net/http"
	"os"

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
	dir := http.Dir(c.HTTPStaticDir)
	m.ServeFiles("/static/*filepath", dir)
	m.NotFound = http.FileServer(dir)

	// Add handlers for all the endpoints
	m.GET("/api/v1/system", authMW(getSystem, AuthLevelView))           // all devices in one shot
	m.PUT("/api/v1/system/mode", authMW(putSystemMode, AuthLevelAdmin)) // heating or cooling
	m.PUT("/api/v1/system/control", authMW(putControl, AuthLevelAdmin)) // manual, schedule, or temp-sensor

	// manual control
	m.PUT("/api/v1/pump/:id/target", authMW(putPump, AuthLevelControl)) // set target state
	m.PUT("/api/v1/blower/:id/target", authMW(putBlower, AuthLevelControl))

	// manual system scheduling
	m.GET("/api/v1/schedule", authMW(getSchedule, AuthLevelView))           // get entire schedule
	m.POST("/api/v1/schedule", authMW(postSchedule, AuthLevelControl))      // add a new entry
	m.PUT("/api/v1/sched/:id", authMW(TODO, AuthLevelControl))              // update an entry
	m.DELETE("/api/v1/sched/:id", authMW(deleteSchedule, AuthLevelControl)) // delete an entry

	// temp/occupancy based scheduling (phase 2, requires sensors)
	m.PUT("/api/v1/zone/:id/targets", authMW(putZone, AuthLevelControl)) // set target temp range for zone
	m.POST("/api/v1/room/:id/schedule", TODO)                            // add an occupancy-expected entry
	m.PUT("/api/v1/room/:id/schedule/:sched", TODO)                      // update an occupancy-expected entry
	m.DELETE("/api/v1/room/:id/schedule/:sched", TODO)                   // remove an occupancy-expected entry

	return m
}

func headers(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		log.Info("request", "url", r.URL, "method", r.Method)
	}

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
