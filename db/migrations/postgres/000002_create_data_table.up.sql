BEGIN;
CREATE TABLE data (
  seq            SERIAL          PRIMARY KEY,
  id             UUID            NOT NULL,
  validator      VARCHAR(64)     NOT NULL,
  namespace      VARCHAR(64)     NOT NULL,
  def_name       VARCHAR(64)     NOT NULL,
  def_version    VARCHAR(64)     NOT NULL,
  hash           CHAR(64)        NOT NULL,
  created        BIGINT          NOT NULL,
  value          JSONB           NOT NULL
);
CREATE UNIQUE INDEX data_id ON data(id);
CREATE INDEX data_hash ON data(namespace,hash);
CREATE INDEX data_created ON data(namespace,created);
COMMIT;