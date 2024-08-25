
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
