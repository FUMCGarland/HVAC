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

func putZoneTemps(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	c := hvac.GetConfig()
	headers(w, r)

	id, err := strconv.ParseInt(ps.ByName("id"), 10, 8)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
	z := (hvac.ZoneID(id)).Get()
	if z == nil {
		err := fmt.Errorf("unknown zone %d", id)
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	zt := hvac.ZoneTargets{}
	if err := json.NewDecoder(r.Body).Decode(&zt); err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	if err := z.SetTargets(c, &zt); err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	fmt.Fprint(w, jsonStatusOK)
}

func putZoneStart(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	headers(w, r)

	inid, err := strconv.ParseInt(ps.ByName("id"), 10, 8)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	id := hvac.ZoneID(inid)
	z := id.Get()
	if z == nil {
		err := fmt.Errorf("unknown zone %d", id)
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	zc := hvac.Command{}
	if err := json.NewDecoder(r.Body).Decode(&zc); err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	if zc.TargetState {
		log.Info("starting zone manually", "id", id, "cmd", zc)
		if err := id.Start(zc.RunTime, "manual"); err != nil {
			log.Error(err.Error())
			http.Error(w, jsonError(err), http.StatusInternalServerError)
			return
		}
	} else {
		log.Info("stopping zone manually", "id", id, "cmd", zc)
		id.Stop("manual")
	}

	fmt.Fprint(w, jsonStatusOK)
}
