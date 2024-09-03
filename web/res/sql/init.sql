
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

--
-- charts view: view with a single row containing a `data` column.  The
-- `data` column is a JSON-encoded map of chart type to chart data.
--
-- * chart type: one of 't' or 'h', indicating temperature or humidity,
--   respectively.
-- * chart data: data needed to generate a line chart using chartjs
--   where the X axis is time and the Y axis is the measurement value.
--
BEGIN;
DROP VIEW IF EXISTS charts;
CREATE VIEW charts(data) AS
  -- modernc.org does not implement the "series" extension so we cannot
  -- use generate_series().  instead, we're a modified version of the
  -- recursive CTE suggested here:
  --
  --   https://www.sqlite.org/series.html
  --
  WITH RECURSIVE times(ts, s) AS (
    -- start time (48 hours ago)
    SELECT unixepoch() - unixepoch() % (15*60) - 48*60*60,
           datetime(unixepoch() - unixepoch() % (15*60) - 48*60*60, 'unixepoch', 'localtime')

    UNION ALL

    -- next value (increment by 15 minutes)
    SELECT ts + (15*60),
           datetime(ts + (15*60), 'unixepoch', 'localtime')
      FROM times
        -- next value (increment by 15 minutes)
     WHERE ts + (15*60) <= (unixepoch() - unixepoch() % (15*60))
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
               'backgroundColor', sensors.color,

               -- cubic interpolation (not working)
               -- ref: https://www.chartjs.org/docs/latest/samples/line/interpolation.html
               'tension', 0.4,

               -- span gaps (also not working)
               -- ref: https://www.chartjs.org/docs/latest/charts/line.html#line-styling
               -- 'spanGaps', true,
               -- 'showLine', true,

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

        FROM (
          -- sort sensors by sort, then name (case-insensitive)
          SELECT * FROM sensors ORDER BY sort, LOWER(name)
        ) sensors
    )
  )) FROM types;
COMMIT;
