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
