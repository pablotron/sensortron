#
# diningroom.py: example sensortron configuration file
#
# edit the configuration below and then copy this file to the pico w
# with the following command:
#
#   pipenv run rshell cp example-configs/diningroom.py /pyboard/config.py
#

# wifi ssid and password (required)
WIFI_SSID = 'YOUR-WIFI-SSID'
WIFI_PASSWORD = 'YOUR-WIFI-PASSWORD'

# api url (required)
URL = 'http://sensors.home.pmdn.org:1979/api/read'

# pseudo-mac secret (optional, defaults to '' if unspecified)
SECRET = 'SOME-SECRET-STRING'

# host name (optional, defaults to pico-$UNIQUE_ID if unspecified)
HOSTNAME = 'pico-diningroom'
