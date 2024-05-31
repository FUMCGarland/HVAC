package rest

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/FUMCGarland/hvac"
	"github.com/FUMCGarland/hvac/log"

	"github.com/julienschmidt/httprouter"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jws"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

const jwtSignerFilename = "signer.jwk"

func mintjwt(username string, level authLevel) (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}

	key, ok := getJWSigningKeys().Key(0)
	if !ok {
		return "", fmt.Errorf("encryption jwk not set")
	}

	jwts, err := jwt.NewBuilder().
		IssuedAt(time.Now()).
		Subject(string(username)).
		Claim("level", level).
		Issuer(hostname).
		JwtID(generateID(16)).
		Audience([]string{sessionName}).
		Expiration(time.Now().Add(time.Hour * 24 * 28)).
		Build()
	if err != nil {
		return "", err
	}

	hdrs := jws.NewHeaders()
	signed, err := jwt.Sign(jwts, jwt.WithKey(jwa.RS256, key, jws.WithProtectedHeaders(hdrs)))
	if err != nil {
		return "", err
	}

	return string(signed[:]), nil
}

func login(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if err := req.ParseMultipartForm(1024 * 64); err != nil {
		log.Warn(err.Error())
		http.Error(res, jsonError(err), http.StatusNotAcceptable)
		return
	}

	username := req.PostFormValue("username")
	if username == "" {
		err := fmt.Errorf("username not set")
		log.Error(err.Error())
		http.Error(res, jsonError(err), http.StatusNotAcceptable)
		return
	}

	password := req.FormValue("password")
	if password == "" {
		err := fmt.Errorf("password not set")
		log.Error(err.Error())
		http.Error(res, jsonError(err), http.StatusNotAcceptable)
		return
	}

	var level authLevel = 0
	authenticated := false
	for k := range ad {
		if strings.EqualFold(ad[k].Username, username) { // ignore case
			if err := bcrypt.CompareHashAndPassword([]byte(ad[k].PwHash), []byte(password)); err != nil {
				log.Error("login failed", "err", err)
				http.Error(res, "Invalid username/password", http.StatusNotAcceptable)
				return
			}
			username = ad[k].Username // use case in config file
			level = ad[k].Level
			authenticated = true
			break
		}
	}

	if !authenticated {
		log.Error("user not found", "username", username)
		http.Error(res, "Invalid username/password", http.StatusNotAcceptable)
		return
	}

	JWT, err := mintjwt(username, level)
	if err != nil {
		log.Error(err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("login", "username", username, "level", level)
	headers(res, req)
	res.Header().Set("content-type", "application/jwt")
	http.SetCookie(res, &http.Cookie{
		Name:     "jwt",
		Value:    JWT,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24 * 365),
		MaxAge:   0,
		Secure:   false,
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
	})
	fmt.Fprint(res, JWT)
}

func generateID(size int) string {
	var characters = strings.Split("abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ", "")
	var buf = make([]byte, size)

	for i := 0; i < size; i++ {
		r, err := rand.Int(rand.Reader, big.NewInt(int64(len(characters))))
		if err != nil {
			log.Error(err.Error())
		}
		b := []byte(characters[r.Int64()])
		buf[i] = b[0]
	}
	return string(buf)
}

func initkey() error {
	c := hvac.GetConfig()

	raw, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal("failed to generate new RSA private key", "error", err.Error())
	}

	key, err := jwk.FromRaw(raw)
	if err != nil {
		log.Fatal("failed to create symmetric key", "error", err.Error())
	}

	_ = key.Set(jwk.KeyIDKey, generateID(16))

	buf, err := json.MarshalIndent(key, "", "  ")
	if err != nil {
		log.Fatal("failed to marshal key into JSON", "error", err.Error())
	}

	pubpath := path.Join(c.StateStore, jwtSignerFilename)
	pubfd, err := os.OpenFile(pubpath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatal("failed to open jwt signer for writing", "error", err.Error())
	}
	defer pubfd.Close()

	if _, err := pubfd.Write(buf); err != nil {
		log.Fatal("unable to write jwt signer", "error", err.Error())
	}
	return nil
}

func getJWSigningKeys() jwk.Set {
	c := hvac.GetConfig()

	keys, err := jwk.ReadFile(path.Join(c.StateStore, jwtSignerFilename))
	if err != nil {
		log.Warn("unable to load jwk signer, creating new", "error", err.Error())
		// first run, or old keys deleted, start anew
		if err := initkey(); err != nil {
			log.Fatal(err.Error())
		}
		// try 2
		keys, err = jwk.ReadFile(path.Join(c.StateStore, jwtSignerFilename))
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	return keys
}
