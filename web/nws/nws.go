// National Weather Service API wrapper
package nws

import (
  "context"
  "encoding/json"
  "fmt"
  "io"
  "net/http"
)

// NWS API problem (error).
type Problem struct {
  Type string `json:"type"` // problem type
  Title string `json:"title"` // problem title
  Status int `json:"status"` // status code
  Detail string `json:"detail"` // detail
  Instance string `json:"instance"` // instance
  CorrelationId string `json:"correlationId"` // correlation ID
}

// return problem error message
func (p Problem) Error() string {
  return fmt.Sprintf("%s: %s", p.Title, p.Detail)
}

// default client user agent 
var defaultUserAgent = "(github.com/pablotron/sensortron, sensortron@pablotron.org)"

// NWS API client.
type Client struct {
  UserAgent string // user-agent string
}

// observation property
type Property struct {
  UnitCode string `json:"unitCode"` // unit code
  Value *float32 `json:"value"` // value
  QualityControl *string `json:"qualityControl"` // quality control
}

// latest observations response
type ObservationsResponse struct {
  Id string `json:"id"` // observation ID
  Type string `json:"type"` // observation type

  Geometry struct {
    Type string `json:"type"` // geometry type
    Coordinates []float32 `json:"coordinates"` // coordinates
  } `json:"geometry"` // geometry

  Properties struct {
    Elevation Property `json:"elevation"` // elevation
    Station string `json:"station"` // station
    Timestamp string `json:"timestamp"` // timestamp
    RawMessage string `json:"rawMessage"` // raw message
    TextDescription string `json:"textDescription"` // text description
    Icon string `json:"icon"` // icon

    // TODO: presentWeather

    Temperature Property `json:"temperature"` // temperature
    DewPoint Property `json:"dewPoint"` // dew point
    WindDirection Property `json:"windDirection"` // wind direction
    WindSpeed Property `json:"windSpeed"` // wind speed
    WindGust Property `json:"windGust"` // wind gust
    BarometricPressure Property `json:"barometricPressure"` // barometric pressure
    SeaLevelPressure Property `json:"seaLevelPressure"` // sea level pressure
    Visibility Property `json:"visibility"` // visibility
    MaxTemperatureLast24Hours Property `json:"maxTemperatureLast24Hours"` // max temp last 24 hours
    MinTemperatureLast24Hours Property `json:"minTemperatureLast24Hours"` // min temp last 24 hours
    PrecipitationLastHour Property `json:"precipitationLastHour"` // precipitation last hour
    PrecipitationLast3Hours Property `json:"precipitationLast3Hours"` // precipitation last 3 hours
    PrecipitationLast6Hours Property `json:"precipitationLast6Hours"` // precipitation last 6 hours
    RelativeHumidity Property `json:"relativeHumidity"` // relative humidity
    WindChill Property `json:"windChill"` // wind chill
    HeatIndex Property `json:"heatIndex"` // heat index

    CloudLayers []struct {
      Base Property `json:"base"` // base
      Amount string `json:"amount"` // amount
    } `json:"cloudLayers"` // cloud layers
  } `json:"properties"` // observed properties
}

// internal shared http client
var httpClient http.Client

// fetch data from NWS API.
func (c Client) fetch(ctx context.Context, path string) (http.Header, []byte, error) {
  // build request url
  url := fmt.Sprintf("https://api.weather.gov/%s", path)

  // get user agent
  userAgent := defaultUserAgent
  if c.UserAgent != "" {
    userAgent = c.UserAgent
  }

  // create request
  req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
  if err != nil {
    return http.Header{}, []byte{}, err
  }

  // add request headers
  req.Header.Add("accept", "application/geo+json")
  req.Header.Add("user-agent", userAgent)

  // send request, get response
  resp, err := httpClient.Do(req)
  if err != nil {
    return http.Header{}, []byte{}, err
  }
  defer resp.Body.Close()

  // check response code
  if resp.StatusCode != 200 {
    // parse problem
    var problem Problem
    if err := json.NewDecoder(resp.Body).Decode(&problem); err != nil {
      // json decode failed: return response headers and json decode error
      return resp.Header, []byte{}, err
    }

    // nws api call failed: return response headers and problem
    return resp.Header, []byte{}, &problem
  }

  // read response body
  body, err := io.ReadAll(resp.Body)
  if err != nil {
    return http.Header{}, []byte{}, err
  }

  // return headers and response body
  return resp.Header, body, nil
}

// Get latest observations for given station.
func (c Client) LatestObservations(ctx context.Context, stationId string) (http.Header, []byte, error) {
  // example: https://api.weather.gov/stations/KDCA/observations/latest?require_qc=false
  return c.fetch(ctx, fmt.Sprintf("/stations/%s/observations/latest", stationId))
}

// Forecast period.
type Period struct {
  Number int `json:"number"` // number
  Name string `json:"name"` // name
  StartTime string `json:"startTime"` // start time
  EndTime string `json:"endTime"` // end time
  IsDaytime bool `json:"isDaytime"` // is daytime?
  Temperature int `json:"temperature"` // temperature
  TemperatureUnit string `json:"temperatureUnit"` // temperature unit
  TemperatureTrend string `json:"temperatureTrend"` // temperature trend
  ProbabilityOfPrecipitation Property `json:"probabilityOfPrecipitation"` // probability of precipitation
  WindSpeed string `json:"windSpeed"` // wind speed
  WindDirection string `json:"windDirection"` // wind direction
  Icon string `json:"icon"` // icon
  ShortForecast string `json:"shortForecast"` // short forecast
  DetailedForecast string `json:"detailedForecast"` // detailed forecast
}

// Forecast response.
type ForecastResponse struct {
  Type string `json:"type"` // type
  Geometry struct {
    Type string `json:"type"` // type
    Coordinates [][][]float32 `json:"coordinates"` // coordinates
  } `json:"geometry"` // geometry

  Properties struct {
    Units string `json:"units"` // units
    ForecastGenerator string `json:"forecastGenerator"` // forecast generator
    GeneratedAt string `json:"generatedAt"` // generated at
    UpdateTime string `json:"updateTime"` // update time
    ValidTimes string `json:"validTimes"` // valid times
    Elevation Property `json:"elevation"` // elevation

    Periods []Period `json:"periods"` // periods
  } `json:"properties"` // properties
}

// Get forecast for given forecast office and grid coordinates.
func (c Client) Forecast(ctx context.Context, wfo string, x, y int) (http.Header, []byte, error) {
  // example: https://api.weather.gov/gridpoints/LWX/91,70/forecast?units=us
  return c.fetch(ctx, fmt.Sprintf("/gridpoints/%s/%d,%d/forecast?units=us", wfo, x, y))
}

// Default NWS API client
var DefaultClient Client

// Get latest observations for given station.
func LatestObservations(ctx context.Context, stationId string) (http.Header, []byte, error) {
  return DefaultClient.LatestObservations(ctx, stationId)
}

// Get forecast for given forecast office and grid coordinates.
func Forecast(ctx context.Context, wfo string, x, y int) (http.Header, []byte, error) {
  return DefaultClient.Forecast(ctx, wfo, x, y)
}
