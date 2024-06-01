package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"

	"github.com/FUMCGarland/hvac/log"

	"github.com/lestrrat-go/jwx/v2/jws"
	"github.com/lestrrat-go/jwx/v2/jwt"
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

// used in the cli to add/update users, should probably just be moved there
const BcryptRounds = 14

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

func authMW(h httprouter.Handle, requiredlevel authLevel) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if len(ad) == 0 {
			err := fmt.Errorf("no auth data")
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token, err := jwt.ParseRequest(r,
			jwt.WithKeySet(sk, jws.WithInferAlgorithmFromKey(true), jws.WithUseDefault(true)),
			jwt.WithValidate(true),
			jwt.WithAudience(sessionName),
			jwt.WithAcceptableSkew(20*time.Second),
		)
		if err != nil {
			log.Error("token parse/validate failed", "error", err.Error())
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		username := string(token.Subject())
		claim, ok := token.Get("level")
		if !ok {
			err := fmt.Errorf("no level claim in token")
			log.Error(err.Error(), "user", username, "claim", claim, "type", fmt.Sprintf("%T", claim))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		checklevel, ok := claim.(float64) // why does this come across as float64?
		if !ok {
			err := fmt.Errorf("authlevel type assertion failed")
			log.Error(err.Error(), "user", username, "claim", claim, "type", fmt.Sprintf("%T", claim))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if authLevel(checklevel) < requiredlevel {
			err := fmt.Errorf("acess level too low")
			log.Warn(err.Error(), "wanted", requiredlevel, "got", checklevel, "username", username)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		h(w, r, ps)
	}
}
