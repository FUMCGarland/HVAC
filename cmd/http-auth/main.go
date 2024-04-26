package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"

	"github.com/FUMCGarland/hvac/rest"
)

const usage = "http-auth [-f filename] [-w] username password"

func main() {
	configPathPtr := flag.String("f", "/etc/hvac-http-auth.json", "Path to the auth file")
	writeBack := flag.Bool("w", false, "write the file back to disk")
	flag.Parse()
	username := flag.Arg(0)
	password := flag.Arg(1)
	if username == "" || password == "" {
		panic(usage)
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), rest.BcryptRounds)
	if err != nil {
		panic(err.Error())
	}

	c, err := rest.LoadAuth(*configPathPtr)
	if err != nil {
		panic(err.Error())
	}

	updated := false
	for k := range c {
		if c[k].Username == username {
			c[k].PwHash = string(bytes)
			updated = true
		}
	}
	if !updated {
		c = append(c, rest.AuthData{Username: username, PwHash: string(bytes), Level: 0})
	}

	m, _ := json.MarshalIndent(c, "", " ")
	if *writeBack {
		if err := os.WriteFile(*configPathPtr, m, 0600); err != nil {
			panic(err.Error())
		}
	} else {
		fmt.Println(string(m))
	}
}
