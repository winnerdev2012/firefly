CREATE TABLE operations (
  id          string         NOT NULL,
  tx_id       string         NOT NULL,
  optype      string         NOT NULL,
  opstatus    string         NOT NULL,
  recipient   string,
  plugin      string         NOT NULL,
  backend_id  string         NOT NULL,
  created     int64          NOT NULL,
  updated     int64,
  error       string         NOT NULL,
  info        blob
);

CREATE UNIQUE INDEX operations_primary ON operations(id);
CREATE INDEX operations_created ON operations(created);
