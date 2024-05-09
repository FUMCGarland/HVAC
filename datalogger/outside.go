package datalogger

import (
	"github.com/FUMCGarland/hvac"
	"github.com/FUMCGarland/hvac/log"

	owm "github.com/briandowns/openweathermap"
)

// my intention is that this system will run for decades without needing updates
// including a third-party API like this makes it all very precarious
// it should fail-safe should the API key expire, or the service shutdown

func getOutsideTemp(c *hvac.Config) (float64, int) {
	if err := owm.ValidAPIKey(c.OpenWeatherMapKey); err != nil {
		return 0.0, 0
	}

	w, err := owm.NewCurrent("F", "EN", c.OpenWeatherMapKey)
	if err != nil {
		log.Error(err.Error())
		return 0.0, 0
	}

	if err := w.CurrentByID(c.OpenWeatherMapID); err != nil {
		log.Error(err.Error())
		return 0.0, 0
	}
	return w.Main.Temp, w.Main.Humidity
}
