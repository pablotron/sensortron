# main.py: pico w micropython which reads bme280 sensor values
# every 10 seconds and posts them to an HTTP endpoint.
#
# notes:
# - reads configuration from config.py
# - expects BME280 sensor to be attached to I2C0 on GP4 and GP5
#
# request notes:
# - sends POST request with JSON-encoded body to config.URL
# - each request includes the following headers
#   - x-unique-id: hex-encoded machine unique ID
#   - x-pseudo-mac-sha256: hex-encoded digest of secret, machine unique
#     ID, and sensor values

import bme280
import config
import hashlib
import json
import gc
import network
import requests
import time

# valid states
STATE_IDLE = rp2.const(0) # idle state, led off
STATE_BUSY = rp2.const(1) # busy state, led flashes
STATE_FAILED = rp2.const(2) # failed state, led on

# init state
state = STATE_IDLE

# get hex-encoded unique ID
id_hex = machine.unique_id().hex()

# get bme280, led, and wifi handles
bme = bme280.BME280(i2c=machine.I2C(0))
led = machine.Pin('LED')
wifi = network.WLAN()

def timestamp():
  """get current time as RFC3339-formatted string"""
  return '{:04}-{:02}-{:02}T{:02}:{:02}:{:02}Z'.format(*(time.gmtime()))

def log(s):
  """print log message with timestamp prefix"""
  print(f'{timestamp()} {s}')

def get_sensor_vals(bme):
  """get dict of bme280 sensor values"""

  # read current timestamp
  e = timestamp()

  # read raw temperature, pressure, and humidity values
  t, p, h = bme.read_compensated_data()

  # adjust values
  t = round(t / 100, 2) # temperature (C)
  p = p // 256 # pressure (Pa)
  h = round((h // 1024) / 100, 2) # humidity (%)

  # return dict of values
  return { 't': t, 'p': p, 'h': h, 'e': e }

def on_led_timer(timer):
  """LED timer callback: read state and adjust led to compenste"""
  if state == STATE_IDLE:
    led.off() # disable led
  elif state == STATE_BUSY:
    led.toggle() # flash led
  elif state == STATE_FAILED:
    led.on() # enable led
  else:
    log(f'unknown state: {state}')

def on_wifi_timer(timer):
  # get wifi status
  wifi_status = wifi.status()

  if wifi_status < 0:
    # wifi connect failed, log error
    log(f'wifi connect failed, status = {wifi_status}')
    timer.deinit()
    state = STATE_FAILED
  elif wifi_status == network.STAT_GOT_IP:
    # wifi connect succeeded, continue
    log(f'connected to {config.WIFI_SSID}')
    timer.deinit()

    # init read timer which fires every 10s, reads the sensor values,
    # and posts them to an HTTP URL
    machine.Timer(-1, mode=machine.Timer.PERIODIC, period=10000, callback=on_read_timer)

    state = STATE_IDLE # set idle state
  else:
    log(f'waiting, status = {wifi_status}') # print wifi status

def on_read_timer(timer):
  """read timer callback: read sensor values and post them to url"""
  try:
    state = STATE_BUSY # set busy state

    # get sensor values
    vals = json.dumps(get_sensor_vals(bme))
    log(f'vals = {vals}') # log values

    # get secret from config
    secret = (config.SECRET if hasattr(config, 'SECRET') else '')
    
    # build pseudo-mac hex digest by concatenating the following:
    # - secret
    # - unique ID,
    # - JSON-encoded sensor values
    digest = hashlib.sha256(secret + id_hex + vals).digest().hex()

    # build request headers
    headers = {
      'x-unique-id': id_hex,
      'x-pseudo-mac-sha256': digest,
      'content-type': 'application/json',
    }

    # post request
    log(f'posting to {config.URL}') # log url
    resp = requests.post(config.URL, headers=headers, data=bytes(vals, 'utf-8'))
    log(f'response code = {resp.status_code}') # log response code

    # trigger gc (prevent OOM)
    gc.collect()

    state = STATE_IDLE # clear state
  except Exception as e:
    state = STATE_BUSY # set busy state
    raise e # re-raise exception

def main():
  """init, connect to wifi, start read timer"""
  try:
    # init led timer which fires every 500ms and sets the LED based on
    # the current state
    machine.Timer(-1, mode=machine.Timer.PERIODIC, period=500, callback=on_led_timer)

    # set country code (default to 'US' if unspecified)
    # country = (config.COUNTRY if hasattr(config, 'COUNTRY') else 'US')
    # log(f'country = {country}')
    # rp2.country(country)

    # set hostname (default to 'pico-$UNIQUE_ID' if unspecified)
    hostname = (config.HOSTNAME if hasattr(config, 'HOSTNAME') else f'pico-{id_hex}')
    log(f'hostname = {hostname}')
    network.hostname(hostname)
    time.sleep(1)

    # connect to wifi
    state = STATE_BUSY # set busy state
    log(f'wifi: activating')
    wifi.active(True)
    log(f'wifi: connecting to {config.WIFI_SSID}...')
    wifi.connect(config.WIFI_SSID, config.WIFI_PASSWORD)

    # wait for wifi connect result
    log('wifi: waiting for IP')
    machine.Timer(-1, mode=machine.Timer.PERIODIC, period=1000, callback=on_wifi_timer)
  except Exception as e:
    state = STATE_FAILED # set failed state
    raise e # re-raise exception

# run main
main()
