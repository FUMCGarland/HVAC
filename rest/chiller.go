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

func putChillerStart(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	headers(w, r)

	inid, err := strconv.ParseInt(ps.ByName("id"), 10, 8)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	id := hvac.ChillerID(inid)
	z := id.Get()
	if z == nil {
		err := fmt.Errorf("unknown chiller %d", id)
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	cc := hvac.Command{}
	if err := json.NewDecoder(r.Body).Decode(&cc); err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	if cc.TargetState {
		log.Info("starting chiller manually", "id", id, "cmd", cc, "user", getUser(r))
		if err := id.Start(cc.RunTime, "manual: "+getUser(r)); err != nil {
			log.Error(err.Error())
			http.Error(w, jsonError(err), http.StatusInternalServerError)
			return
		}
	} else {
		log.Info("stopping chiller manually", "id", id, "cmd", cc, "user", getUser(r))
		id.Stop("manual: " + getUser(r))
	}

	fmt.Fprint(w, jsonStatusOK)
}
