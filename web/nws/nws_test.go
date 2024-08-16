package nws

import (
  "context"
  "encoding/json"
  "fmt"
  "testing"
)

func TestLatestObservations(t *testing.T) {
  ctx := context.Background()

  // tests expected to pass
  passTests := []string {
    "KDCA",
  }

  // run pass tests
  for _, test := range(passTests) {
    t.Run(test, func(t *testing.T) {
      // get latest observations
      headers, body, err := LatestObservations(ctx, test);
      t.Logf("headers = %#v", headers)
      t.Logf("body = %#v", string(body))
      t.Logf("err = %#v", err)
      if err != nil {
        t.Fatal(err)
      }

      // parse observations
      var observations ObservationsResponse
      if err := json.Unmarshal(body, &observations); err != nil {
        t.Fatal(err)
      }
      t.Logf("observations = %#v", observations)
    })
  }

  // tests expected to fail
  failTests := []string {
    "XXXX",
  }

  // run fail tests
  for _, test := range(failTests) {
    t.Run(test, func(t *testing.T) {
      headers, body, err := LatestObservations(ctx, test);
      t.Logf("headers = %#v", headers)
      t.Logf("body = %#v", string(body))
      t.Logf("err = %#v", err)
      if err == nil {
        t.Fatal("got success, expected error")
      }
    })
  }
}

func TestForecast(t *testing.T) {
  ctx := context.Background()

  // tests expected to pass
  passTests := []struct {
    wfo string
    x, y int
  } {
    { "LWX", 91, 70 },
  }

  // run pass tests
  for _, test := range(passTests) {
    // build test name
    testName := fmt.Sprintf("%s/%d,%d", test.wfo, test.x, test.y)

    t.Run(testName, func(t *testing.T) {
      // get forecast
      headers, body, err := Forecast(ctx, test.wfo, test.x, test.y);
      t.Logf("headers = %#v", headers)
      t.Logf("body = %#v", string(body))
      t.Logf("err = %#v", err)
      if err != nil {
        t.Fatal(err)
      }

      // parse forecast
      var forecast ForecastResponse
      if err := json.Unmarshal(body, &forecast); err != nil {
        t.Fatal(err)
      }
      t.Logf("forecast = %#v", forecast)
    })
  }

  // tests expected to fail
  failTests := []struct {
    name string // test name
    wfo string
    x, y int
  } {
    { "bad wfo", "XXX", 91, 70 },
    { "bad x", "LWX", 999999, 70 },
    { "bad y", "LWX", 91, 999999 },
  }

  // run fail tests
  for _, test := range(failTests) {
    t.Run(test.name, func(t *testing.T) {
      headers, body, err := Forecast(ctx, test.wfo, test.x, test.y);
      t.Logf("headers = %#v", headers)
      t.Logf("body = %#v", string(body))
      t.Logf("err = %#v", err)
      if err == nil {
        t.Fatal("got success, expected error")
      }
    })
  }
}
