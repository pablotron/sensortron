CREATE VIEW charts(data) AS
  WITH times AS (
    SELECT value AS ts,
           datetime(value, 'unixepoch', 'localtime') AS s
      FROM generate_series(
        -- start time (24 hours ago)
        unixepoch() - unixepoch() % (15*60) - 24*60*60,

        -- end time (most recent 15 minute tick)
        unixepoch() - unixepoch() % (15*60),

        -- time series increment (15 minutes)
        15*60
      )
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

-- json_object('data', json_group_array(b.data ->>'$."' || ) from times a left join history b on b.ts = a.ts; -- SELECT a.ts, a.s, b.data ->> '$."outside".t' as t from times a left join history b on b.ts = a.ts order by a.ts desc limit 50;
