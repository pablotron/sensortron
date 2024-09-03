package main

import (
  db_sql "database/sql"
  "encoding/json"
  "errors"
  _ "embed"
  _ "modernc.org/sqlite"
)

var tableExistsSql = `
  SELECT EXISTS(
    SELECT true
      FROM sqlite_master
     WHERE type = 'table'
       AND tbl_name = ?
  )
`

var errTableExistsFailed = errors.New("tableExists failed")

func dbTableExists(db *db_sql.DB, tableName string) (bool, error) {
  rows, err := db.Query(tableExistsSql, tableName)
  if err != nil {
    return false, err
  }
  defer rows.Close()

  if !rows.Next() {
    return false, errTableExistsFailed
  }

  // scan into
  var r bool
  if err := rows.Scan(&r); err != nil {
    return false, err
  }

  // check for row iteration error
  if rows.Err() != nil {
    return false, rows.Err()
  }

  // return result
  return r, nil
}

//go:embed res/sql/init.sql
var dbInitSql string

// Check for "sensors" table and if it does not exist, then initialize
// database.
func dbMaybeInit(db *db_sql.DB) error {
  // check to see if sensors table exists
  sensorsExists, err := dbTableExists(db, "sensors")
  if err != nil {
    return err
  }

  if sensorsExists {
    // database is already initialized, do nothing
    return nil
  }

  // init db, return result
  _, err = db.Exec(dbInitSql)
  return err
}

// Connect to database and initialize it.  Panics on error.
func dbOpen(dbPath string) *db_sql.DB {
  // connect to db, check for error
  db, err := db_sql.Open("sqlite", dbPath)
  if err != nil {
    panic(err)
  }

  // initialize database
  if err := dbMaybeInit(db); err != nil {
    panic(err)
  }

  // return database
  return db
}

var historyInsertSql = `
  INSERT INTO history(ts, data) VALUES (?, ?)
`

// add entry to history table.
func dbHistoryInsert(db *db_sql.DB, ts int64, data map[string]SensorData) error {
  // encode data
  buf, err := json.Marshal(data)
  if err != nil {
    return err
  }

  // insert row
  _, err = db.Exec(historyInsertSql, ts, buf)
  return err
}

var sensorInsertSql = `
  INSERT INTO sensors(id, name, color, sort)
    SELECT a.column1, -- id
           a.column1, -- name
           a.column2, -- color
           0          -- sort
      FROM (VALUES (?, ?)) a
     WHERE a.column1 NOT IN (SELECT id FROM sensors)
`

// Insert sensor into sensors table with if the sensor does not already
// exist in the sensors table.
func dbSensorInsert(db *db_sql.DB, id, color string) error {
  // insert row
  _, err := db.Exec(sensorInsertSql, id, color)
  return err
}

var chartDataSql = `
  SELECT data FROM charts
`

// Read chart data from charts view and return string containing
// JSON-encoded chart data.
func dbChartData(db *db_sql.DB) (string, error) {
  var s string
  err := db.QueryRow(chartDataSql).Scan(&s)
  return s, err
}

var sensorsSql = `
  SELECT id,
         name,
         color,
         sort
    FROM sensors
   ORDER BY sort, LOWER(name)
`

// Get ordered list of sensors from sensors table.
func dbSensors(db *db_sql.DB) ([]Sensor, error) {
  r := make([]Sensor, 0)

  // exec query, get result
  rows, err := db.Query(sensorsSql)
  if err != nil {
    return r, err
  }
  defer rows.Close()

  // build result
  for rows.Next() {
    // get sensor properties
    var s Sensor
    if err := rows.Scan(&s.Id, &s.Name, &s.Color, &s.Sort); err != nil {
      return r, err
    }

    // append sensor to results
    r = append(r, s)
  }

  // check for error
  if err := rows.Err(); err != nil {
    return r, err
  }

  // return result
  return r, nil
}

var sensorUpdateSql = `
  UPDATE sensors
     SET name = ?, color = ?, sort = ?
   WHERE id = ?
`

// Update sensor in sensors table.
func dbSensorUpdate(db *db_sql.DB, s Sensor) error {
  // update row
  _, err := db.Exec(sensorUpdateSql, s.Name, s.Color, s.Sort, s.Id)
  return err
}
