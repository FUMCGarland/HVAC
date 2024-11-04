package rest

import (
	"fmt"
	"net/http"

	"github.com/FUMCGarland/hvac/log"
	"github.com/julienschmidt/httprouter"
)

func getRunlog(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	headers(w, r)

	w.Header().Set("Content-Type", "text/plain")

	d, err := log.ReadBuf()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err = fmt.Fprint(w, d); err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
