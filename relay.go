package hvac

type Relay struct {
	Pin      uint8    // i2c pin on relay control board
	PumpID   PumpID   // non-zero if devices is a pump
	BlowerID BlowerID // non-zero if device is a blower

	Running   bool  // is the device currently running
	StartTime int64 // time of current start time.Unix
	StopTime  int64 // time to stop running - time.Unix
	RunTime   int64 // total run-time of the device (seconds since process start)
}
