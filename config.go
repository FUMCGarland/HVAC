package hvac

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/FUMCGarland/hvac/log"
)

// cmdChan is the channel used to pass commands to the MQTT server
var cmdChan chan MQTTRequest

// Config is the system configuration and run-time data
// a subset of Config is used for the startup configuration
type Config struct {
	StateStore        string       // the directory in which to store running state
	SystemMode        SystemModeT  // heating/cooling
	ControlMode       ControlModeT // off/manual/schedule/temp
	MQTT              *MQTTConfig  // config for the mqtt server
	HTTPaddr          string       // the address on which to listen :8080
	HTTPStaticDir     string       // the directory that contains the built webui
	HTTPAuthData      string       // the file that contains the HTTP authentication creds
	DataLogFile       string       // the name of the datalog files
	OpenWeatherMapKey string       // API key from OpenWeatherMap
	OpenWeatherMapID  int          // locaiton ID in OpenWeatherMap
	ChillerLockout    bool         // is the chiller locked-out due to being too cold?
	BoilerLockout     bool         // is the boiler locked-out due to being to warm?
	Blowers           []Blower
	Chillers          []Chiller
	Dampers           []Damper
	Loops             []Loop
	Pumps             []Pump
	Rooms             []Room
	Valves            []Valve
	Zones             []Zone
}

// MQTTConfig is the configuration of the MQTT subsystem
type MQTTConfig struct {
	ID         string // something randomish (fumcg-hvac-server)
	Auth       string // filename (/etc/hvac-mqtt-auth.json)
	Root       string // base mqtt path component(s) (fumcg)
	ListenAddr string // (":1883")
}

// c is the global running config
// start with sane defaults if the config file isn't fully populated
var c *Config = &Config{
	StateStore: "/var/hvac",
	MQTT: &MQTTConfig{
		Root:       "hvac",
		Auth:       "/etc/hvac-mqtt-auth.json",
		ID:         "hvac",
		ListenAddr: ":1883",
	},
	HTTPaddr:          ":8080",
	HTTPStaticDir:     "/usr/local/hvac",
	HTTPAuthData:      "/etc/hvac-http-auth.json",
	DataLogFile:       "/var/hvac/datalog.csv",
	OpenWeatherMapKey: "",      // no key, do not check OWM
	OpenWeatherMapID:  4693003, // Garland, TX
	ChillerLockout:    false,
	BoilerLockout:     false,
	SystemMode:        SystemModeCool,
}

// init() considered harmful, except that this is trivial and sets up a global
func init() {
	cmdChan = make(chan MQTTRequest)
}

// LoadConfig is called from main() to set up the running configuration
func LoadConfig(filename string) (*Config, error) {
	raw, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err.Error())
	}

	// overwrite the defaults with what is in the file
	if err := json.Unmarshal(raw, c); err != nil {
		log.Fatal(err.Error())
	}

	if err := validate(); err != nil {
		log.Fatal("validate", "config", c, "error", err.Error())
	}

	if err := c.loadFromStore(); err != nil {
		log.Fatal("loadFromStore", "config", c, "error", err.Error())
	}

	return c, nil
}

// GetConfig returns the global running config, used by various sub-modules
func GetConfig() *Config {
	return c
}

// validate is called by LoadConfig to make sure that things are sane
// TODO fill in all the rest
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
			if zone.ID == room.CoolZone {
				found = true
				break
			}
		}
		if !found {
			return (fmt.Errorf("room %d set to invalid cool zone %d", room.ID, room.CoolZone))
		}
	}
	for _, room := range c.Rooms {
		for _, zone := range c.Zones {
			if zone.ID == room.HeatZone {
				found = true
				break
			}
		}
		if !found {
			return (fmt.Errorf("room %d set to invalid heat zone %d", room.ID, room.HeatZone))
		}
	}
	return nil
}

// GetMQTTChan returns the command channel to the MQTT subsystem
func GetMQTTChan() chan MQTTRequest {
	return cmdChan
}

// StopAll sends a command to all devices to shut down.
func StopAll() {
	for k := range c.Chillers {
		if !c.Chillers[k].Running {
			continue
		}
		c.Chillers[k].ID.Stop("stop all")
		c.Chillers[k].Running = false
		time.Sleep(1 * time.Second) // prevent current inrush/overload
	}

	for k := range c.Pumps {
		if !c.Pumps[k].Running {
			continue
		}
		c.Pumps[k].ID.Stop("stop all")
		c.Pumps[k].Running = false
		time.Sleep(1 * time.Second) // prevent overload
	}

	for k := range c.Blowers {
		if !c.Blowers[k].Running {
			continue
		}
		c.Blowers[k].ID.Stop("stop all")
		c.Blowers[k].Running = false
		time.Sleep(1 * time.Second) // prevent overload
	}
}
