// world's crappiest sensor display
package main

import (
  "encoding/csv"
  "encoding/json"
  "embed"
  "fmt"
  "io"
  io_fs "io/fs"
  "log"
  "net/http"
  "slices"
  "strings"
  "time"

  "github.com/go-chi/chi/v5"
  "github.com/go-chi/chi/v5/middleware"
)

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

// /api/poll handler
func doApiPoll(w http.ResponseWriter, r *http.Request) {
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

// Get filename timestamp suffix in YYYYMMDDHHMMSS format.
func timestampSuffix() string {
  now := time.Now()
  return fmt.Sprintf("%04d%02d%02d%02d%02d%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
}

// /api/download/current handler
func doApiDownloadCurrent(w http.ResponseWriter, r *http.Request) {
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

// /api/forecast handler
func doApiForecast(w http.ResponseWriter, r *http.Request) {
  w.Header().Add("content-type", "application/json")

  // open mock forecast
  f, err := resFs.Open("res/forecast.json")
  if err != nil {
    log.Print(fmt.Errorf("forecast.json: %w", err)) // log error
    http.Error(w, "error", 500)
    return
  }
  defer f.Close()

  // send forecast
  if _, err := io.Copy(w, f); err != nil {
    log.Print(err) // log error
    http.Error(w, "error", 500)
    return
  }
}

// /api/download/forecast handler
func doApiDownloadForecast(w http.ResponseWriter, r *http.Request) {
  // build content-disposition header value
  disposition := fmt.Sprintf("attachment; filename=\"sensortron-current-data-%s.csv\"", timestampSuffix())

  // TODO: replace this with real forecast download

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

// mime types to compress
var compressedTypes = []string {
  "application/json",
  "text/css",
  "test/csv",
  "text/javascript",
  "text/html",
}

func main() {
  // get assets subdirectory
  assetsDir, err := io_fs.Sub(resFs, "res/assets")
  if err != nil {
    panic(err)
  }

  // init router
  r := chi.NewRouter()
  r.Use(middleware.Logger)
  r.Use(middleware.Compress(5, compressedTypes...))

  // add routes
  r.Post("/api/read", doApiRead)
  r.Post("/api/poll", doApiPoll)
  r.Get("/api/download/current", doApiDownloadCurrent)
  r.Post("/api/forecast", doApiForecast)
  r.Get("/api/download/forecast", doApiDownloadForecast)
  r.Get("/", doHtmlFile("home.html"))
  r.Get("/about/", doHtmlFile("about.html"))
  r.Handle("/assets/*", http.StripPrefix("/assets", http.FileServerFS(assetsDir)))

  // serve
  log.Fatal(http.ListenAndServe(":1979", r))
}
