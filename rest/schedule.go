package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	// "strconv"

	"github.com/FUMCGarland/hvac"
	"github.com/FUMCGarland/hvac/log"
	"github.com/julienschmidt/httprouter"
)

func getSchedule(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := hvac.GetConfig()
	headersMW(w, r)

	sched, err := c.GetSchedule()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(sched); err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}

func postSchedule(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c := hvac.GetConfig()
	headersMW(w, r)

	if err := getAuth(r); err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusForbidden)
		return
	}

	e := hvac.ScheduleEntry{}
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
	log.Info("adding schedule entry", "e", e)

	schedule, err := c.GetSchedule()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	if err := schedule.AddEntry(&e); err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	fmt.Fprint(w, jsonStatusOK)
}
