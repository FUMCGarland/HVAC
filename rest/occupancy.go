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

func getOccupancySchedule(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := hvac.GetConfig()
	headers(w, r)

	o, err := c.GetOccupancySchedule()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(o); err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}

func postOccupancyRecurring(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c := hvac.GetConfig()
	headers(w, r)

	e := hvac.OccupancyRecurringEntry{}
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
	log.Info("adding occupancy entry", "e", e)

	o, err := c.GetOccupancySchedule()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	if err := o.AddRecurringEntry(&e); err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	fmt.Fprint(w, jsonStatusOK)
}

func deleteOccupancyRecurring(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := hvac.GetConfig()
	headers(w, r)

	inid, err := strconv.ParseInt(ps.ByName("id"), 10, 8)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
	log.Info("removing occupancy entry", "id", inid)

	o, err := c.GetOccupancySchedule()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	o.RemoveRecurringEntry(hvac.OccupancyRecurringID(inid))

	fmt.Fprint(w, jsonStatusOK)
}

func putOccupancyRecurring(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := hvac.GetConfig()
	headers(w, r)

	inid, err := strconv.ParseInt(ps.ByName("id"), 10, 8)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	e := hvac.OccupancyRecurringEntry{}
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
	if hvac.OccupancyRecurringID(inid) != e.ID {
		err := fmt.Errorf("url ID does not match incoming json")
		log.Error(err.Error(), "id", inid, "e", e)
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	log.Info("updating occupancy entry", "id", inid, "e", e)
	o, err := c.GetOccupancySchedule()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
	if err := o.EditRecurringEntry(&e); err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, jsonStatusOK)
}

func postOccupancyOneTime(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c := hvac.GetConfig()
	headers(w, r)

	e := hvac.OccupancyOneTimeEntry{}
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
	log.Info("adding occupancy entry", "e", e)

	o, err := c.GetOccupancySchedule()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	if err := o.AddOneTimeEntry(&e); err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	fmt.Fprint(w, jsonStatusOK)
}

func deleteOccupancyOneTime(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := hvac.GetConfig()
	headers(w, r)

	inid, err := strconv.ParseInt(ps.ByName("id"), 10, 8)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
	log.Info("removing occupancy entry", "id", inid)

	o, err := c.GetOccupancySchedule()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	o.RemoveOneTimeEntry(hvac.OccupancyOneTimeID(inid))

	fmt.Fprint(w, jsonStatusOK)
}

func putOccupancyOneTime(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := hvac.GetConfig()
	headers(w, r)

	inid, err := strconv.ParseInt(ps.ByName("id"), 10, 8)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	e := hvac.OccupancyOneTimeEntry{}
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
	if hvac.OccupancyOneTimeID(inid) != e.ID {
		err := fmt.Errorf("url ID does not match incoming json")
		log.Error(err.Error(), "id", inid, "e", e)
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	log.Info("updating occupancy entry", "id", inid, "e", e)
	o, err := c.GetOccupancySchedule()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
	if err := o.EditOneTimeEntry(&e); err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, jsonStatusOK)
}
