package main

import (
  "os"
  "strconv"
)

// Sensortron configuration.
type Config struct {
  DbPath string // database path
  StationId string // observation station ID
  Wfo string // weather forecasting office (WFO)
  GridX int // forecast grid X
  GridY int // forecast grid Y
}

// Get string from environment variable or return default value.
func envStr(key, defaultValue string) string {
  if s := os.Getenv(key); s != "" {
    return s
  } else {
    return defaultValue
  }
}

// Get integer value from environment variable or return default value.
func envInt(key string, defaultValue int) int {
  // get value from environment
  s := os.Getenv(key)
  if s == "" {
    return defaultValue
  }

  // convert to integer
  if r, err := strconv.Atoi(s); err != nil {
    panic(err)
  } else {
    return r
  }
}

// Get configuration from environment.
func NewConfigFromEnv() Config {
  return Config {
    DbPath: envStr("SENSORTRON_DB_PATH", "/data/sensortron.db"),
    StationId: envStr("SENSORTRON_STATION_ID", "KDCA"),
    Wfo: envStr("SENSORTRON_FORECAST_WFO", "LWX"),
    GridX: envInt("SENSORTRON_FORECAST_GRID_X", 91),
    GridY: envInt("SENSORTRON_FORECAST_GRID_Y", 70),
  }
}
