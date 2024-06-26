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

func putBlower(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	headers(w, r)

	inid, err := strconv.ParseInt(ps.ByName("id"), 10, 8)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
	id := hvac.BlowerID(inid)
	p := id.Get()
	if p == nil {
		err := fmt.Errorf("unknown blower %d", id)
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	bc := hvac.Command{}
	if err := json.NewDecoder(r.Body).Decode(&bc); err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	if bc.TargetState {
		log.Info("starting blower manually", "id", id, "cmd", bc, "user", getUser(r))
		if err := id.Start(bc.RunTime, "manual: "+getUser(r)); err != nil {
			log.Error(err.Error())
			http.Error(w, jsonError(err), http.StatusNotAcceptable)
			return
		}
	} else {
		log.Info("stopping blower manually", "id", id, "cmd", bc, "user", getUser(r))
		if id.Stop("manual: " + getUser(r)); err != nil {
			log.Error(err.Error())
			http.Error(w, jsonError(err), http.StatusNotAcceptable)
			return
		}
	}

	fmt.Fprint(w, jsonStatusOK)
}

func resetFilter(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	headers(w, r)

	inid, err := strconv.ParseInt(ps.ByName("id"), 10, 8)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
	id := hvac.BlowerID(inid)
	p := id.Get()
	if p == nil {
		err := fmt.Errorf("unknown blower %d", id)
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
	p.FilterTime = 0

	fmt.Fprint(w, jsonStatusOK)
}
