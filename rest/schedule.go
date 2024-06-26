package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/FUMCGarland/hvac"
	"github.com/FUMCGarland/hvac/log"
	"github.com/julienschmidt/httprouter"
)

func getSchedule(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := hvac.GetConfig()
	headers(w, r)

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
	headers(w, r)

	e := hvac.ScheduleEntry{}
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
	log.Debug("adding schedule entry", "e", e)

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

func deleteSchedule(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := hvac.GetConfig()
	headers(w, r)

	inid, err := strconv.ParseInt(ps.ByName("id"), 10, 8)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
	log.Debug("removing schedule entry", "id", inid)

	schedule, err := c.GetSchedule()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	schedule.RemoveEntry(uint8(inid))

	fmt.Fprint(w, jsonStatusOK)
}

func putSchedule(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := hvac.GetConfig()
	headers(w, r)

	inid, err := strconv.ParseInt(ps.ByName("id"), 10, 8)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	e := hvac.ScheduleEntry{}
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
	if uint8(inid) != e.ID {
		err := fmt.Errorf("url ID does not match incoming json")
		log.Error(err.Error(), "id", inid, "e", e)
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	log.Debug("updating schedule entry", "id", inid, "e", e)
	schedule, err := c.GetSchedule()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
	if err := schedule.EditEntry(&e); err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, jsonStatusOK)
}
