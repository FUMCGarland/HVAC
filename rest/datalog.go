package rest

import (
	"io"
	"net/http"
	"os"

	"github.com/FUMCGarland/hvac"
	"github.com/FUMCGarland/hvac/log"
	"github.com/julienschmidt/httprouter"
)

func getDatalog(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	headers(w, r)

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=datalog.csv")

	file, err := os.Open(hvac.GetConfig().DataLogFile)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if _, err = io.Copy(w, file); err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
