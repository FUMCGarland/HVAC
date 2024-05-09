package hvac

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/FUMCGarland/hvac/log"
)

var c *Config
var cmdChan chan MQTTRequest

type Config struct {
	// the directory in which to store running state
	StateStore string
	// heating/cooling
	SystemMode SystemModeT
	// off/manual/schedule/temp
	ControlMode ControlModeT
	// config for the mqtt server
	MQTT *MQTTConfig
	// the address on which to listen :8080
	HTTPaddr string
	// the directory that contains the built webui
	HTTPStaticDir string
	// the file that contains the HTTP authentication creds
	HTTPAuthData string
	// the location of the datalog files
	DataLogFile       string
	OpenWeatherMapKey string
	OpenWeatherMapID  int
	Blowers           []Blower
	Chillers          []Chiller
	Dampers           []Damper
	Loops             []Loop
	Pumps             []Pump
	Rooms             []Room
	Valves            []Valve
	Zones             []Zone
}

type MQTTConfig struct {
	ID         string // something randomish (fumcg-hvac-server)
	Auth       string // filename (/etc/hvac-mqtt-auth.json)
	Root       string // base mqtt path component(s) (fumcg)
	ListenAddr string // (":1883")
}

var defaults *Config = &Config{
	StateStore: "/var/hvac",
	MQTT: &MQTTConfig{
		Root:       "fumcg",
		Auth:       "/etc/hvac-mqtt-auth.json",
		ID:         "fumcg",
		ListenAddr: ":1883",
	},
	HTTPaddr:          ":8080",
	HTTPStaticDir:     "/usr/local/hvac",
	HTTPAuthData:      "/etc/hvac-http-auth.json",
	DataLogFile:       "/var/hvac/datalog.csv",
	OpenWeatherMapKey: "",
	OpenWeatherMapID:  4693003,
}

func init() {
	cmdChan = make(chan MQTTRequest)
}

func LoadConfig(filename string) (*Config, error) {
	// config.Load is called early, this is probably an OK place for this
	_ = log.Start()

	raw, err := os.ReadFile(filename)
	if err != nil {
		panic(err.Error())
	}

	in := defaults
	// overwrite the defaults with what is in the file
	if err := json.Unmarshal(raw, &in); err != nil {
		panic(err.Error())
	}

	c = in

	if err := validate(); err != nil {
		log.Error("config", "config", c)
		panic(err.Error())
	}

	c.loadFromStore()

	return c, nil
}

func GetConfig() *Config {
	if c == nil {
		panic("GetConfig() called before LoadConfig()")
	}

	return c
}

func validate() error {
	if err := validateBlower(); err != nil {
		return err
	}
	if err := validatePumps(); err != nil {
		return err
	}
	if err := validateRooms(); err != nil {
		return err
	}
	return nil
}

func validateBlower() error {
	log.Debug("blowers", "blowers", c.Pumps)
	found := false
	for _, blower := range c.Blowers {
		// loop
		found = false
		for _, loop := range c.Loops {
			if loop.ID == blower.HotLoop {
				found = true
				break
			}
		}
		if !found && blower.HotLoop != 0 {
			return (fmt.Errorf("blower %d set to invalid hot loop %d", blower.ID, blower.HotLoop))
		}
		found = false
		for _, loop := range c.Loops {
			if loop.ID == blower.ColdLoop {
				found = true
				break
			}
		}
		if !found {
			return (fmt.Errorf("blower %d set to invalid cold loop %d", blower.ID, blower.ColdLoop))
		}

		// zone
		found = false
		for _, zone := range c.Zones {
			if zone.ID == blower.Zone {
				found = true
				break
			}
		}
		if !found {
			return (fmt.Errorf("blower %d set to invalid zone %d", blower.ID, blower.Zone))
		}
	}

	return nil
}

func validatePumps() error {
	log.Debug("pumps", "pumps", c.Pumps)
	found := false
	for _, pump := range c.Pumps {
		for _, loop := range c.Loops {
			if loop.ID == pump.Loop {
				found = true
				break
			}
		}
		if !found {
			return (fmt.Errorf("pump %d set to invalid loop %d", pump.ID, pump.Loop))
		}
	}
	return nil
}

func validateRooms() error {
	log.Debug("rooms", "rooms", c.Rooms)
	found := false
	for _, room := range c.Rooms {
		for _, zone := range c.Zones {
			if zone.ID == room.Zone {
				found = true
				break
			}
		}
		if !found {
			return (fmt.Errorf("room %d set to invalid zone %d", room.ID, room.Zone))
		}
	}
	return nil
}

func GetMQTTChan() chan MQTTRequest {
	return cmdChan
}

func StopAll() {
	for k := range c.Pumps {
		if !c.Pumps[k].Running {
			continue
		}
		c.Pumps[k].ID.Stop("manual")
		c.Pumps[k].Running = false
		time.Sleep(1 * time.Second)
	}

	for k := range c.Blowers {
		if !c.Blowers[k].Running {
			continue
		}
		c.Blowers[k].ID.Stop("manual")
		c.Blowers[k].Running = false
		time.Sleep(1 * time.Second)
	}
}
