# sensortron

Sensortron source code.  

- `pico-w/`: Micropython source code for Pico W which will read from an
  I2C BME280 temperature sensor every 10 seconds and post the
  measurements to an HTTP URL.
- `web/`: Go source code and Dockerfile to build container which shows a
  minimal web interface and also exposes a REST API to accept
  temperature sensor readings.

Setup instructions are available in each directory.
