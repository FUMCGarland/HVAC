#/bin/bash

set mqtt_pass="sensor"

curl http://$1/settings
echo 
curl http://$1/settings/cloud\?enabled=0
curl http://$1/settings\?mqtt_enable=true\&mqtt_server=10.0.0.45:1883\&mqtt_user=sensor\&mqtt_pass=sensor\&discoverable=true\&mqtt_update_period=0\&temperature_units=F\&mqtt_update_period=3600
echo
curl http://$1/reboot

