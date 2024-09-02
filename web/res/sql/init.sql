
CREATE TABLE sensors (
  -- sensor unique id
  id TEXT PRIMARY KEY,

  -- display name
  name TEXT UNIQUE NOT NULL,

  -- display name
  color CHAR(6) NOT NULL,

  sort INT NOT NULL DEFAULT 0
);

CREATE TABLE history (
  -- entry id
  id INTEGER PRIMARY KEY AUTOINCREMENT,

  -- timestamp (seconds epoch)
  ts BIGINT NOT NULL,

  -- sensor data
  data JSON NOT NULL
);

CREATE INDEX in_history_ts ON history(ts);

-- create charts view
DROP VIEW IF EXISTS charts;
CREATE VIEW charts(data) AS
  -- normally we would use generate_series(), but the modernc.org/sqlite
  -- driver does not implement the series extension
  WITH RECURSIVE times(ts, s) AS (
    -- start time (24 hours ago)
    SELECT unixepoch() - unixepoch() % (15*60) - 24*60*60,
           datetime(unixepoch() - unixepoch() % (15*60) - 24*60*60, 'unixepoch', 'localtime')

    UNION ALL

    -- next value (increment by 15 minutes)
    SELECT ts + (15*60),
           datetime(ts + (15*60), 'unixepoch', 'localtime')
      FROM times
        -- next value (increment by 15 minutes)
     WHERE ts + (15*60) <= (unixepoch() - unixepoch() % (15*60))
  -- ), times AS (
  --   SELECT value AS ts,
  --          datetime(value, 'unixepoch', 'localtime') AS s
  --     FROM generate_series(
  --       -- start time (24 hours ago)
  --       unixepoch() - unixepoch() % (15*60) - 24*60*60,

  --       -- end time (most recent 15 minute tick)
  --       unixepoch() - unixepoch() % (15*60),

  --       -- time series increment (15 minutes)
  --       15*60
  --     )
  ), types AS (
    -- chart types
    SELECT column1 AS id, -- chart id ("t" or "h")
           column2 AS scale, -- scaling factor
           column3 AS offset, -- value offset
           column4 AS rounding -- decimal precision
      FROM (VALUES
        ('t', 1.8, 32, 2), -- temperature
        ('h', 100, 0, 1) -- humidity
      )
  )

  -- build a json-encoded hash which maps a chart type ('t' or 'h') to the
  -- chart data for the given chart type.
  SELECT json_group_object(types.id, json_object(
    -- timestamps (X axis labels)
    'labels', (SELECT json_group_array(times.s) FROM times ORDER BY ts),

    -- array of data sets for this chart type
    'datasets', (
      SELECT json_group_array(json_object(
               -- data set label (sensor name)
               'label', sensors.name,

               -- data set line style
               'borderWidth', 1,
               'borderColor', sensors.color,

               -- data set measurements
               'data', (
                 SELECT json_group_array(
                          -- do the following:
                          -- 1. extract measurement value at time from history table
                          -- 2. scale and offset value
                          -- 3. round value to decimal precision for chart type
                          ROUND(
                            types.scale * (b.data ->> ('$."' || sensors.id || '".' || types.id)) + types.offset,
                            types.rounding
                          )
                        )

                   FROM times a
                   LEFT JOIN history b
                     ON b.ts = a.ts

                  ORDER BY a.ts
               )
             ))

        FROM sensors

       -- sort sensors by sort column, then by name (case-insensitive)
       ORDER BY sensors.sort,
                LOWER(sensors.name)
    )
  )) FROM types;
