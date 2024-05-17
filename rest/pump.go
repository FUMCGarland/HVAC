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
	headers(w, r)

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

	pc := hvac.Command{}
	if err := json.NewDecoder(r.Body).Decode(&pc); err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	if pc.TargetState {
		log.Info("starting pump manually", "id", id, "cmd", pc)
		if err := id.Start(pc.RunTime, "manual"); err != nil {
			log.Error(err.Error())
			http.Error(w, jsonError(err), http.StatusInternalServerError)
			return
		}
	} else {
		log.Info("stopping pump manually", "id", id, "cmd", pc)
		id.Stop("manual")
	}

	fmt.Fprint(w, jsonStatusOK)
}
