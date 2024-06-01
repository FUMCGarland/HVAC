# FUMC Garland HVAC controller

Multiple components:

## Primary Server (golang, on raspberry pi)

- Coordinate control of devices and UI
- provide https REST interface to UI
- communicate with control-points over MQTT

Phase 1 (done):
- schedule relay states
- manually control relays

Phase 2: (done)
- listen to mqtt for temp sensor data and adjust relays to keep zones in temp range
- based on room schedule occupancy

Phase 3: (started)
- ML/AI optimization of the schedule based on occupancy

Phase 4:
- adjust flow controls (dampers on blowers/ducts, valves on water loops)

## MQTT server / http/REST server / logic

## Control points (raspberry pi with 8-relay units)
- Listen to MQTT, change relay states on command
- "dumb" units which start/stop based on duration
- Using Waveshare RPi Relay Board ( https://www.waveshare.com/wiki/RPi_Relay_Board_(B) )
  
## Temp sensor 
- Send data to MQTT (Shelly devices will do this job, others too)
- Using Shelly H&T

## Cold/Hot water flow controls (phase N)
- Arduino devices connected to existing valves / dampers
- Listen to MQTT, adjust based on command

# Model

## Zone
A collection of physical spaces (rooms) which are controlled together, old-school schedule controls zones

## Room
A physical space - temperature recorded to room, occupancy schedule to roooms

## Loop
A physical loop of pipe that moves either hot or cold water to a zone or zones.

## Pump
Moves hot or cold water in a loop

## Blower
Connected to loops and ducts, moves air to a zone

## Valve
Controls flow on a loop
(phase 4)

## Damper
Controls flow in a duct
(phase 4)

## Schedule Entry
A time for relays to be engaged.

## Target Temp Range
The target temp for a zone at a given time

# Example
(phase 1 - schedule based - complete)
Desire: cool 127 on Monday Morning
data: Room 127 is in zone D, which requires (pumps 1 + blower X) to cool.
command: Run zone D on Monday from  5 AM - 7 AM

(phase 2 - occupancy/temp based - complete)
Desire: keep Sch. hall comfortable on Wednesday evening
Command: Sch. Hall occupied Wed 5-8 PM
Outcome: if temp in zone B is above 78, run (pump 1 & blower Z) until sensor reports temp <= 72 starting at 4 PM

(phase 3 - occupancy/temp based + Machine Learning prediction - in progress)
Same as Phase 2, but looking at weather forcast + historical data to determine when to best run the zones to keep them in target range

(phase 4 - not started)
close down dampners and valves in areas near the bottom of the target range for their zones, open up those near the top of the range...

## assumptions/rules

Blowers can run without the loop pump
Loop pumps can only run when the blowers are running (lest the chiller freeze over)
Boilers run when water temp is too low, we will not control them
Chiller runs hybrid, we must start/stop it, but it has some self-protection logic

Zones will have primary temp ranges
If a room is scheduled to be occupied, the zone will be adjusted to that room's selected range (avg. across zone)

## research
 
things to read:

https://niektemme.com/2015/08/09/smart-thermostat/

https://medium.com/devthoughts/linear-regression-with-go-ff1701455bcd

https://www.waveshare.com/wiki/RPi_Relay_Board_(B)
