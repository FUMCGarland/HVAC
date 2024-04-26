# FUMC Garland HVAC controller

Multiple components:

## Primary Server (golang, on raspberry pi)

- Coordinate control of devices and UI
- provide https REST interface to UI
- communicate with control-points over MQTT

Phase 1:
- schedule relay state on control points (send changes via mqtt)
- manual override on control points (send changes via mqtt)

Phase 2:
- listen to mqtt for temp sensor data and adjust relays to keep zones in temp range

Phase 3:
- adjust flow controls (dampers, valves)

## MQTT server (process on Primary server pi)

## Control points (raspberry pi with 8-relay units)

- Listen to MQTT, change relay states on command
- Fall back to "failsafe" schedule mode if MQTT inaccessible 

## Temp sensor 
- Send data to MQTT

## Flow controls (phase 3)
- Arduino devices connected to existing valves / dampers
- Listen to MQTT, adjust based on command
- Failsafe when MQTT not available

# Model

## Zone
A collection of physical spaces (rooms, hallways, etc) which are controlled together.
(list in phase 1)

## Room
A physical space
(list in phase 1)

## Loop
A physical loop of pipe that moves either hot or cold water to a zone or zones.
(list in phase 1)

## Pump
Moves hot or cold water in a loop
(control in phase 1)

## Blower
Connected to loops and ducts, moves air to a zone
(control in phase 1)

## Valve
Controls flow on a loop
(list in phase 1, control in phase 3)

## Damper
Controls flow in a duct
(list in phase 1, control in phase 3)

## Schedule Entry
(phase 1)
A time for relays to be engaged.

## Target Temp Range
(phase 2)
The target temp for a zone at a given time

# Example
(phase 1)
Desire: cool 127 on Monday Morning
data: Room 127 is in zone D, which requires (pumps 1 & 2 + blower X) to cool.
command: Run (pumps 1&2 + blower X) on Monday from  5 AM - 7 AM

(phase 2)
Desire: keep Sch. hall comfortable on Wednesday evening
data: Sch. is in zone B, which requires (pump 1 & blower Z) to cool
Command: if temp in zone B is above 78, run (pump 1 & blower Z) until sensor reports temp <=70

(phase 3)
...
close down dampners and valves in areas near the bottom of the target range for their zones, open up those near the top of the range...

## assumptions/rules

Blowers can run without the loop pump
Loop pumps can only run when the blowers are running (lest the chiller freeze over)

Zones will have primary temp ranges
If a room is scheduled to be occupied, the zone will be adjusted to that room's selected range

## research
 
things to read: https://niektemme.com/2015/08/09/smart-thermostat/

