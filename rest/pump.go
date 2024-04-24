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

func putPump(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	headersMW(w, r)

	if err := getAuth(r); err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusForbidden)
		return
	}

	inid, err := strconv.ParseInt(ps.ByName("id"), 10, 8)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	id := hvac.PumpID(inid)
	p := id.Get()
	if p == nil {
		err := fmt.Errorf("unknown pump %d", id)
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	pc := hvac.PumpCommand{}
	if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	if pc.TargetState {
		if err := id.Start(pc.RunTime, "manual"); err != nil {
			log.Error(err.Error())
			http.Error(w, jsonError(err), http.StatusNotAcceptable)
			return
		}
	} else {
		if err := id.Stop("manual"); err != nil {
			log.Error(err.Error())
			http.Error(w, jsonError(err), http.StatusNotAcceptable)
			return
		}
	}

	fmt.Fprint(w, jsonStatusOK)
}
