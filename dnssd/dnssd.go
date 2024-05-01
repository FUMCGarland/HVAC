package hvacdnssd

import (
	"context"
	"strconv"
	"strings"

	"github.com/FUMCGarland/hvac"
	"github.com/FUMCGarland/hvac/log"
	"github.com/brutella/dnssd"
)

func Start(ctx context.Context, c *hvac.Config) {
	parts := strings.Split(c.HTTPaddr, ":")
	p := parts[len(parts)-1]
	port, err := strconv.ParseInt(p, 10, 16)
	if err != nil {
		panic(err.Error())
	}

	svhttp, err := dnssd.NewService(dnssd.Config{
		Name:   "HVAC Controller",
		Type:   "_http._tcp",
		Ifaces: []string{"eth0"},
		Domain: "local",
		Port:   int(port),
	})
	if err != nil {
		panic(err.Error())
	}

	svmqtt, err := dnssd.NewService(dnssd.Config{
		Name:   "HVAC Controller",
		Type:   "_mqtt._tcp",
		Ifaces: []string{"eth0"},
		Domain: "local",
		Port:   1883,
	})
	if err != nil {
		panic(err.Error())
	}

	rp, err := dnssd.NewResponder()
	if err != nil {
		panic(err.Error())
	}

	_, err = rp.Add(svhttp)
	if err != nil {
		panic(err.Error())
	}

	_, err = rp.Add(svmqtt)
	if err != nil {
		panic(err.Error())
	}

	log.Info("starting DNSSD")
	rp.Respond(ctx)
}
