package rest

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"

	"github.com/FUMCGarland/hvac/log"
)

var ad []AuthData

type AuthData struct {
	Username string
	PwHash   string
	Level    authLevel
}

type authLevel uint8

const (
	AuthLevelView    authLevel = iota // view can see running system state
	AuthLevelControl                  // control can manually stop/start devices and adjust schedule entries
	AuthLevelAdmin                    // admin can stop the system and change heat/cool
)

const BcryptRounds = 4

func LoadAuth(path string) ([]AuthData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	in := make([]AuthData, 0)
	if err = json.Unmarshal(data, &in); err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return in, nil
}

func authMW(h httprouter.Handle, level authLevel) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if len(ad) == 0 {
			log.Error("no http auth data")
			http.Error(w, "No Auth Data", http.StatusInternalServerError)
			return
		}

		username, password, ok := r.BasicAuth()
		if !ok {
			log.Error("BasicAuth !ok")
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		for k := range ad {
			if ad[k].Username == username {
				err := bcrypt.CompareHashAndPassword([]byte(ad[k].PwHash), []byte(password))
				if err != nil {
					log.Error("login failed", "err", err)
					w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
				if ad[k].Level < level {
					log.Error("login level too low", "wanted", level, "got", ad[k].Level)
					w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
					http.Error(w, "Unauthorized", http.StatusForbidden)
					return
				}
				break
			}
		}
		h(w, r, ps)
	}
}
