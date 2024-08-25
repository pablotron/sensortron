// world's crappiest sensor display
package main

import (
  "context"
  db_sql "database/sql"
  "encoding/csv"
  "encoding/json"
  "embed"
  "fmt"
  "io"
  io_fs "io/fs"
  "log"
  "net/http"
  "sensortron/nws"
  "slices"
  "strings"
  "time"

  "github.com/go-chi/chi/v5"
  "github.com/go-chi/chi/v5/middleware"
)

// global database
var db *db_sql.DB

//go:embed res
var resFs embed.FS

// bme280 sensor data
type SensorData struct {
  T float32 `json:"t"` // temperature, in degrees celcius
  H float32 `json:"h"` // humidity, percentage
  P float32 `json:"p"` // pressure, in pascals
  E string `json:"e"` // rfc3339-formatted timestamp
}

// id to friendly name
var names = map[string]string {
  "e6614c311b4b4537": "living room",
  "e6614c311b267237": "dining room",
  "e6614c311b867937": "nadine's office",
  "keylime": "bedroom",
}

// latest values
var latest = make(map[string]SensorData)

// /api/read handler
func doApiRead(w http.ResponseWriter, r *http.Request) {
  // decode sensor readings from request body
  var data SensorData
  if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
    log.Print(err) // log error
    http.Error(w, "error", 500)
    return
  }

  // get unique id and pseudo-mac
  id := r.Header.Get("x-unique-id")
  mac := r.Header.Get("x-pseudo-mac-sha256")

  // save latest values, log result
  latest[id] = data
  log.Printf("id = %s, mac = %s, data = %#v", id, mac, data)

  // respond with success
  w.Header().Add("content-type", "application/json")
  if _, err := w.Write([]byte("null")); err != nil {
    log.Print(err) // log error
    return
  }
}

type PollRow struct {
  Id string `json:"id"` // sensor ID
  Name string `json:"name"` // room name
  Data SensorData `json:"data"` // sensor data
}

// Get a slice of poll rows sorted by room name.
func getPollRows() []PollRow {
  // build list of sensor readings
  rows := make([]PollRow, 0, len(latest))
  for id, sr := range(latest) {
    // get room name
    name := names[id]
    if name == "" {
      name = id
    }

    // append to rows
    rows = append(rows, PollRow { id, name, sr })
  }

  // sort rows by name
  slices.SortFunc(rows, func(a, b PollRow) int {
		return strings.Compare(a.Name, b.Name)
	})

  return rows
}

// /api/home/current/poll handler
func doApiHomeCurrentPoll(w http.ResponseWriter, r *http.Request) {
  // get sorted list of latest readings
  rows := getPollRows()

  // return rows
  w.Header().Add("content-type", "application/json")
  if err := json.NewEncoder(w).Encode(&rows); err != nil {
    log.Print(err) // log error
    http.Error(w, "error", 500)
    return
  }
}

// /api/home/current/edit handler
func doApiHomeCurrentEdit(w http.ResponseWriter, r *http.Request) {
  // parse args
  var args map[string]string
  if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
    log.Print(err) // log error
    http.Error(w, "error", 500)
    return
  }

  if len(args["name"]) > 0 {
    // save name
    names[args["id"]] = args["name"]
  }

  // respond with success
  w.Header().Add("content-type", "application/json")
  if _, err := w.Write([]byte("null")); err != nil {
    log.Print(err) // log error
    return
  }
}

// Get filename timestamp suffix in YYYYMMDDHHMMSS format.
func timestampSuffix() string {
  now := time.Now()
  return fmt.Sprintf("%04d%02d%02d%02d%02d%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
}

// /api/home/current/download handler
func doApiHomeCurrentDownload(w http.ResponseWriter, r *http.Request) {
  // build content-disposition header value
  disposition := fmt.Sprintf("attachment; filename=\"sensortron-current-data-%s.csv\"", timestampSuffix())

  // build csv rows
  rows := make([][]string, 0, len(latest) + 1)
  rows = append(rows, []string { "Location", "Temperature (F)", "Humidity (%)" })
  for _, row := range(getPollRows()) {
    rows = append(rows, []string {
      row.Name,
      fmt.Sprintf("%2.2f", row.Data.T * 9.0/5.0 + 32.0),
      fmt.Sprintf("%2.1f", row.Data.H * 100.0),
    })
  }

  // send headers
  w.Header().Add("content-type", "text/csv; charset=utf-8")
  w.Header().Add("content-disposition", disposition)

  // send rows
  if err := csv.NewWriter(w).WriteAll(rows); err != nil {
    log.Print(err) // log error
    http.Error(w, "error", 500)
    return
  }
}

// latest cached forecast
var cachedForecast = []byte("null")

// /api/home/forecast/poll handler
func doApiHomeForecastPoll(w http.ResponseWriter, r *http.Request) {
  w.Header().Add("content-type", "application/json")

  // send cached forecast
  if _, err := w.Write(cachedForecast); err != nil {
    log.Print(err) // log error
    http.Error(w, "error", 500)
    return
  }
}

// home forecast csv column headers
var homeForecastCsvCols = []string {
  "Period",
  "Temperature (F)",
  "Wind Speed",
  "Wind Direction",
  "Short Forecast",
  "Detailed Forecast",
}

// /api/home/forecast/download handler
func doApiHomeForecastDownload(w http.ResponseWriter, r *http.Request) {
  // build content-disposition header value
  disposition := fmt.Sprintf("attachment; filename=\"sensortron-forecast-%s.csv\"", timestampSuffix())

  // parse cached forecast
  var forecast nws.ForecastResponse
  if err := json.Unmarshal(cachedForecast, &forecast); err != nil {
    log.Print(err) // log error
    http.Error(w, "error", 500)
    return
  }

  // build csv rows
  rows := make([][]string, 0, len(latest) + 1)
  rows = append(rows, homeForecastCsvCols)
  for _, row := range(forecast.Properties.Periods) {
    rows = append(rows, []string {
      row.Name,
      fmt.Sprintf("%d", row.Temperature),
      row.WindSpeed,
      row.WindDirection,
      row.ShortForecast,
      row.DetailedForecast,
    })
  }

  // send headers
  w.Header().Add("content-type", "text/csv; charset=utf-8")
  w.Header().Add("content-disposition", disposition)

  // send rows
  if err := csv.NewWriter(w).WriteAll(rows); err != nil {
    log.Print(err) // log error
    http.Error(w, "error", 500)
    return
  }
}

// return handler which sends named HTML file.
func doHtmlFile(name string) func(http.ResponseWriter, *http.Request) {
  // build full path
  path := fmt.Sprintf("res/html/%s", name)

  // return handler
  return func(w http.ResponseWriter, r *http.Request) {
    // open file
    f, err := resFs.Open(path)
    if err != nil {
      log.Print(fmt.Errorf("%s: %w", name, err)) // log error
      http.Error(w, "error", 400)
      return
    }
    defer f.Close()

    if _, err := io.Copy(w, f); err != nil {
      log.Print(fmt.Errorf("%s: %w", name, err)) // log error
      http.Error(w, "error", 400)
      return
    }
  }
}

// response types to compress
var compressedResponseTypes = []string {
  "application/json",
  "text/css",
  "test/csv",
  "text/javascript",
  "text/html",
}

// fetch latest NWS observations from given station and convert them to
// sensor data
func fetchLatestObservations(ctx context.Context, stationId string) (SensorData, error) {
  var observations nws.ObservationsResponse

  // get latest observations
  _, body, err := nws.LatestObservations(ctx, stationId)
  if err != nil {
    return SensorData{}, err
  }

  // parse observations body
  if err = json.Unmarshal(body, &observations); err != nil {
    return SensorData{}, err
  }

  // build and return sensor data
  return SensorData {
    T: *observations.Properties.Temperature.Value,
    H: *observations.Properties.RelativeHumidity.Value / 100.0,
    P: *observations.Properties.BarometricPressure.Value,
    E: observations.Properties.Timestamp,
  }, nil
}

func main() {
  // get config from environment
  config := NewConfigFromEnv()
  log.Printf("config = %#v", config)

  // get assets subdirectory from embedded resources
  assetsDir, err := io_fs.Sub(resFs, "res/assets")
  if err != nil {
    panic(err)
  }

  // connect to database
  db = dbOpen(config.DbPath)
  defer db.Close()

  // fetch outside temperature every 10 minutes
  go func() {
    ctx := context.Background()

    // parse sleep duration
    delay, err := time.ParseDuration("10m")
    if err != nil {
      panic(err)
    }

    // loop forever
    for {
      // get latest observations as sensordata
      if data, err := fetchLatestObservations(ctx, config.StationId); err != nil {
        log.Print(err)
      } else {
        log.Printf("NWS data = %#v", data)
        latest["outside"] = data
      }

      // sleep until next fetch interval
      time.Sleep(delay)
    }
  }()

  // fetch current forecast every 4 hours
  go func() {
    ctx := context.Background()

    // parse sleep duration
    delay, err := time.ParseDuration("4h")
    if err != nil {
      panic(err)
    }

    // loop forever
    for {
      // get current forecast
      _, body, err := nws.Forecast(ctx, config.Wfo, config.GridX, config.GridY)
      if err != nil {
        log.Print(err)
      } else {
        // update cached forecast
        cachedForecast = body
        log.Printf("NWS forecast = %s", string(body))
      }

      // sleep until next fetch interval
      time.Sleep(delay)
    }
  }()

  // add history entry every 15 minutes
  go func() {
    // loop forever
    for {
      now := time.Now()

      // is minute of current time divisible by 15?
      if now.Minute() % 15 == 0 {
        // get unix timestamp, rounded to nearest 15 minutes
        unix := now.Unix() - (now.Unix() % 900)

        // log current entries to database
        if err := dbHistoryInsert(db, unix, latest); err != nil {
          log.Print(err)
        }
      }

      // sleep for 1 minute (FIXME)
      time.Sleep(time.Minute)
    }
  }()

  // init router and middleware
  r := chi.NewRouter()
  r.Use(middleware.Logger)
  r.Use(middleware.Compress(5, compressedResponseTypes...))

  // add routes
  r.Post("/api/read", doApiRead)
  r.Post("/api/home/current/poll", doApiHomeCurrentPoll)
  r.Post("/api/home/current/edit", doApiHomeCurrentEdit)
  r.Get("/api/home/current/download", doApiHomeCurrentDownload)
  r.Post("/api/home/forecast/poll", doApiHomeForecastPoll)
  r.Get("/api/home/forecast/download", doApiHomeForecastDownload)
  r.Get("/", doHtmlFile("home.html"))
  r.Get("/about/", doHtmlFile("about.html"))
  r.Get("/forecasts/", doHtmlFile("forecasts.html"))
  r.Handle("/assets/*", http.StripPrefix("/assets", http.FileServerFS(assetsDir)))

  // serve
  log.Fatal(http.ListenAndServe(":1979", r))
}
