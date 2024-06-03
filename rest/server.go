package rest

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwk"

	"github.com/FUMCGarland/hvac"
	"github.com/FUMCGarland/hvac/log"
)

var srv *http.Server
var sk jwk.Set

const sessionName string = "HVAC"
const jsonType = "application/json; charset=UTF-8"
const jsonTypeShort = "application/json"
const jsonStatusOK = `{"status":"ok"}`

// Start launches the HTTP server which is responsible for the frontend and the HTTP API.
func Start(c *hvac.Config, done <-chan bool) {
	var err error
	if ad, err = LoadAuth(c.HTTPAuthData); err != nil {
		log.Fatal(err.Error())
	}

	srv = &http.Server{
		Handler:           getServeMux(c),
		Addr:              c.HTTPaddr,
		WriteTimeout:      (30 * time.Second),
		ReadTimeout:       (30 * time.Second),
		ReadHeaderTimeout: (2 * time.Second),
	}

	// creates the keys if needed
	sk = getJWSigningKeys()

	log.Info("Starting up REST server", "on", c.HTTPaddr)
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err.Error())
		}
	}()

	<-done
	log.Info("Shutting down REST server")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Error(err.Error())
	}
}

func jsonError(e error) string {
	return fmt.Sprintf(`{"status":"error","error":"%s"}`, e.Error())
}

func contentTypeIs(req *http.Request, check string) bool {
	contentType := strings.Split(strings.Replace(req.Header.Get("Content-Type"), " ", "", -1), ";")[0]
	return strings.EqualFold(contentType, check)
}
