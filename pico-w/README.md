# Raspberry Pi Pico W

## Instructions

1. get pico w micropython from here: <https://micropython.org/download/RPI_PICO_W/>
2. plug in pico w with button pressed.  will show up as usb mass storage
   device.
3. copy uf2 from step #1 to usb mass storage device. example:
   `cp -v ~/dl/RPI_PICO_W-20240602-v1.23.0.uf2 /media/pabs/RPI-RP2/`
4. install pipenv (e.g. apt-get install pipenv`).
5. run `pipenv install` in this directory to set up virtual environment
   and install dependencies
6. make a copy of `example-configs/diningroom.py`. example:
   `cp example-configs/diningroom.py example-configs/basement.py`
7. edit copied config from previous step and set the following
   variables: WIFI_SSID, WIFI_PASSWORD, URL, SECRET, and HOSTNAME.
8. Run the following rshell below to copy the files to the pico w.

Rshell commands:

    # create /lib directory on pico w and copy bme280.py to it
    pipenv run rshell mkdir /pyboard/lib
    pipenv run rshell cp bme280.py /pyboard/lib/
    
    # copy new config and main.py to pico w
    pipenv run rshell cp configs/basement.py /pyboard/config.py
    pipenv run rshell cp main.py /pyboard/
    
    # enter repl and monitor activity
    pipenv run rshell repl

Afterwards you should be able to power-cycle the board and it will start
posting sensor readings to the given API URL every 10 seconds.
