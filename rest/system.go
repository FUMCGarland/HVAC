package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/FUMCGarland/hvac"
	"github.com/FUMCGarland/hvac/log"

	"github.com/julienschmidt/httprouter"
)

func getSystem(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c := hvac.GetConfig()
	headers(w, r)

	if err := json.NewEncoder(w).Encode(c); err != nil {
		log.Error("json.Encode failed in getSystem", err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
}

func putSystemMode(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c := hvac.GetConfig()
	headers(w, r)

	if !contentTypeIs(r, jsonTypeShort) {
		err := fmt.Errorf("invalid request format: use JSON")
		log.Warn(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	sm := hvac.SystemMode{}
	if err := json.NewDecoder(r.Body).Decode(&sm); err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	// enforce ControlMode == off before allowing SystemMode changes
	if c.ControlMode != hvac.ControlOff {
		err := fmt.Errorf("cannot change system mode while system is running, set control mode to off")
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	switch sm.Mode {
	case hvac.SystemModeHeat, hvac.SystemModeCool:
		if err := c.SetSystemMode(sm.Mode); err != nil {
			log.Error(err.Error())
			http.Error(w, jsonError(err), http.StatusNotAcceptable)
		}
	default:
		err := fmt.Errorf("unknown SystemMode %d", sm.Mode)
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, jsonStatusOK)
}

func putControl(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c := hvac.GetConfig()
	headers(w, r)

	if !contentTypeIs(r, jsonTypeShort) {
		err := fmt.Errorf("invalid request format: use JSON")
		log.Warn(err.Error())
		http.Error(w, jsonError(err), http.StatusNotAcceptable)
		return
	}

	scm := hvac.ControlMode{}
	if err := json.NewDecoder(r.Body).Decode(&scm); err != nil {
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}

	log.Info("setting system control mode", "mode", scm.ControlMode)

	switch scm.ControlMode {
	case hvac.ControlManual, hvac.ControlSchedule, hvac.ControlTemp, hvac.ControlOff:
		hvac.StopAll()
		if err := c.SetControlMode(scm.ControlMode); err != nil {
			log.Error("set control mode failed", "error", err.Error())
			http.Error(w, jsonError(err), http.StatusInternalServerError)
			return
		}
	default:
		err := fmt.Errorf("unknown ControlMode %d", scm.ControlMode)
		log.Error(err.Error())
		http.Error(w, jsonError(err), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, jsonStatusOK)
}
