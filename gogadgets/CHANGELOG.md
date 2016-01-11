# CHANGELOG

## 0.0.1
 - first version

## 0.0.2
 - added a thermostat
   the config for a heater should look like:
         {
             "host": "http://192.168.1.30:6111",
             "gadgets": [
                 {
                     "location": "the lab",
                     "name": "temperature",
                     "pin": {
                         "type": "thermometer",
                         "OneWireId": "28-0000041cb544",
                         "Units": "F"
                     }
                 },
                 {
                     "location": "the lab",
                     "name": "heater",
                     "pin": {
                         "type": "thermostat",
                         "port": "8",
                         "pin": "11",
                         "args": {
                             "type": "heater",
                             "sensor": "the lab temperature",
                             "high": 150.0,
                             "low": 120.0
                         }
                     }
                 }
             ]
         }
         
