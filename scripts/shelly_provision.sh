#/bin/bash

curl http://$1/settings/cloud\?enabled=0
curl http://$1/settings\?mqtt_enable=true\&mqtt_server=10.0.0.45:1883\&mqtt_user=sensor\&mqtt_pass=sensor\&discoverable=true\&temperature_units=F\&mqtt_update_period=3600
curl http://$1/reboot
echo

