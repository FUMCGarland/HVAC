package hvacdnssd

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/FUMCGarland/hvac"
	"github.com/FUMCGarland/hvac/log"
	"github.com/brutella/dnssd"
)

// Start runs the DNSSD announcer - it remains alive as long as the ctx is valid.
func Start(ctx context.Context, c *hvac.Config) {
	parts := strings.Split(c.HTTPaddr, ":")
	p := parts[len(parts)-1]
	port, err := strconv.ParseInt(p, 10, 16)
	if err != nil {
		log.Fatal(err.Error())
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err.Error())
	}

	svhttp, err := dnssd.NewService(dnssd.Config{
		Name: fmt.Sprintf("HVAC Controller - %s (rest)", hostname),
		Type: "_http._tcp",
		Port: int(port),
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	svmqtt, err := dnssd.NewService(dnssd.Config{
		Name: fmt.Sprintf("HVAC Controller - %s (mqtt)", hostname),
		Type: "_mqtt._tcp",
		Port: 1883,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	rp, err := dnssd.NewResponder()
	if err != nil {
		log.Fatal(err.Error())
	}

	h, err := rp.Add(svhttp)
	if err != nil {
		log.Fatal(err.Error())
	}
	h.UpdateText(map[string]string{"JOSH": "Joyous Online Scheduler for HVAC"}, rp)

	h, err = rp.Add(svmqtt)
	if err != nil {
		log.Fatal(err.Error())
	}
	h.UpdateText(map[string]string{"api_proto": "mqtt", "destination_port": "1883"}, rp)

	log.Debug("starting DNSSD")
	if err := rp.Respond(ctx); err != nil && err != context.Canceled {
		log.Error("unable to start DNSSD responder", "error", err.Error())
	}
}
