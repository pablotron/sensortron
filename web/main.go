// world's crappiest sensor display
package main

import (
  "encoding/json"
  "embed"
  io_fs "io/fs"
  "log"
  "net/http"
  "slices"
  "strings"
  "github.com/go-chi/chi/v5"
  "github.com/go-chi/chi/v5/middleware"
)

//go:embed res/assets/*
var resFs embed.FS

// bme280 sensor readings
type SensorReadings struct {
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

// latest readings
var latest = make(map[string]SensorReadings)

// /api/read handler
func doApiRead(w http.ResponseWriter, r *http.Request) {
  // decode sensor readings from request body
  var sr SensorReadings
  if err := json.NewDecoder(r.Body).Decode(&sr); err != nil {
    log.Print(err) // log error
    http.Error(w, "error", 500)
    return
  }

  // get unique id and pseudo-mac
  id := r.Header.Get("x-unique-id")
  mac := r.Header.Get("x-pseudo-mac-sha256")

  // save latest reading, log result
  latest[id] = sr
  log.Printf("id = %s, mac = %s, data = %#v", id, mac, sr)

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
  Temp float32 `json:"temp"` // temperature
}

// Get a slice of poll rows sorted by room name.
func getPollRows() []PollRow {
  // build list of sensor readings
  rows := make([]PollRow, 0, len(latest))
  for id, sensorRow := range(latest) {
    // get room name
    name := names[id]
    if name == "" {
      name = id
    }

    // append to rows
    rows = append(rows, PollRow {
      Id: id,
      Name: name,
      Temp: sensorRow.T,
    })
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

// home page html
const HOME_HTML = `<!DOCTYPE html>
<html data-bs-theme='dark'>
  <head>
    <meta charset='utf-8'/>
    <meta name='viewport' content='width=device-width, initial-scale=1'/>

    <title>SENSORTRON (beta)</title>
    <link rel='stylesheet' type='text/css' href='/assets/bootstrap-5.3.3/css/bootstrap.min.css'/>
  </head>

  <body>
    <nav class='navbar navbar-expand-lg bg-body-tertiary'>
      <div class='container-fluid'>
        <a
          href='/'
          class='navbar-brand'
          title='SENSORTRON (beta)'
          aria-label='SENSORTRON (beta)'
        >SENSORTRON (beta)</a>

        <button
          type='button'
          class='navbar-toggler'
          data-bs-toggle='collapse'
          data-bs-target='#navbar-content'
          aria-controls='navbar-content'
          aria-expanded='false'
          title='Toggle navigation'
          aria-label='Toggle navigation'
        >
          <span class='navbar-toggler-icon'></span>
        </button>

        <div id='navbar-content' class='collapse navbar-collapse'>
          <ul class='navbar-nav me-auto mb-2 mb-lg-0'>
            <li class='nav-item'>
              <a
                href='/'
                class='nav-link active'
                aria-current='page'
                title='Home'
                aria-label='Home'
              >Home</a>
            </li>

            <li class='nav-item'>
              <a
                href='/about/'
                class='nav-link'
                title='About'
                aria-label='About'
              >About</a>
            </li>
          </ul>
        </div>
      </div>
    </nav>

    <div class='container'>
      <br/>

      <div class='card'>
        <div class='card-header'>
          <b>Current Measurements</b>
        </div>

        <table id='current-temps' class='table table-hover'>
        </table>

        <div class='card-footer'>
          <div class='btn-group' role='group'>
            <input
              type='radio'
              id='unit-celcius'
              class='unit btn-check'
              name='unit'
              value='c'
              autocomplete='off'
              title='Celcius'
              aria-label='Celcius'
            />
            <label
              for='unit-celcius'
              class='unit btn btn-sm btn-outline-primary'
              title='Celcius'
              aria-label='Celcius'
            >
              Celcius
            </label>

            <input
              type='radio'
              id='unit-fahrenheit'
              class='unit btn-check'
              name='unit'
              value='f'
              autocomplete='off'
              title='Fahrenheit'
              aria-label='Fahrenheit'
              checked
            />
            <label
              for='unit-fahrenheit'
              class='unit btn btn-sm btn-outline-primary'
              title='Fahrenheit'
              aria-label='Fahrenheit'
            >
              Fahrenheit
            </label>
          </div>
        </div>
      </div>
    </div>

    <script type='text/javascript' src='/assets/bootstrap-5.3.3/js/bootstrap.bundle.js'></script>
    <script type='text/javascript' src='/assets/home.js'></script>
  </body>
</html>`

// home page handler
func doHome(w http.ResponseWriter, r *http.Request) {
  w.Header().Add("content-type", "text/html; charset=utf-8")

  if _, err := w.Write([]byte(HOME_HTML)); err != nil {
    log.Print(err) // log error
    return
  }
}

// about page html
const ABOUT_HTML = `<!DOCTYPE html>
<html data-bs-theme='dark'>
  <head>
    <meta charset='utf-8'/>
    <meta name='viewport' content='width=device-width, initial-scale=1'/>

    <title>ABOOT SENSORTRON (beta)</title>
    <link rel='stylesheet' type='text/css' href='/assets/bootstrap-5.3.3/css/bootstrap.min.css'/>
  </head>

  <body>
    <nav class='navbar navbar-expand-lg bg-body-tertiary'>
      <div class='container-fluid'>
        <a
          href='/'
          class='navbar-brand'
          title='SENSORTRON (beta)'
          aria-label='SENSORTRON (beta)'
        >SENSORTRON (beta)</a>

        <button
          type='button'
          class='navbar-toggler'
          data-bs-toggle='collapse'
          data-bs-target='#navbar-content'
          aria-controls='navbar-content'
          aria-expanded='false'
          title='Toggle navigation'
          aria-label='Toggle navigation'
        >
          <span class='navbar-toggler-icon'></span>
        </button>

        <div id='navbar-content' class='collapse navbar-collapse'>
          <ul class='navbar-nav me-auto mb-2 mb-lg-0'>
            <li class='nav-item'>
              <a
                href='/'
                class='nav-link'
                title='Home'
                aria-label='Home'
              >Home</a>
            </li>

            <li class='nav-item'>
              <a
                href='/about/'
                class='nav-link active'
                aria-current='page'
                title='About'
                aria-label='About'
              >About</a>
            </li>
          </ul>
        </div>
      </div>
    </nav>

    <div class='container'>
      <br/>

      <div class='card'>
        <div class='card-header'>
          <b>ABOOT SENSORTRON (beta)</b>
        </div>

        <div class='card-body'>
          synergy
        </dev>
      </div>
    </div>

    <script type='text/javascript' src='/assets/bootstrap-5.3.3/js/bootstrap.bundle.js'></script>
    <script type='text/javascript' src='/assets/about.js'></script>
  </body>
</html>`

// about page handler
func doAbout(w http.ResponseWriter, r *http.Request) {
  w.Header().Add("content-type", "text/html; charset=utf-8")

  if _, err := w.Write([]byte(ABOUT_HTML)); err != nil {
    log.Print(err) // log error
    return
  }
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

  // add routes
  r.Post("/api/read", doApiRead)
  r.Post("/api/poll", doApiPoll)
  r.Get("/", doHome)
  r.Get("/about/", doAbout)
  r.Handle("/assets/*", http.StripPrefix("/assets", http.FileServerFS(assetsDir)))

  // serve
  log.Fatal(http.ListenAndServe(":1979", r))
}
