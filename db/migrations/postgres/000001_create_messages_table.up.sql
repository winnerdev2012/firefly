CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

BEGIN;
CREATE TABLE messages (
  seq         SERIAL          PRIMARY KEY,
  id          UUID            NOT NULL,
  cid         CHAR(36),
  mtype       VARCHAR(64)     NOT NULL,
  author      VARCHAR(1024)   NOT NULL,
  created     BIGINT          NOT NULL,
  namespace   VARCHAR(64)     NOT NULL,
  topics      VARCHAR(1024)   NOT NULL,
  tag         VARCHAR(64)     NOT NULL,
  group_id    UUID,
  datahash    CHAR(64)        NOT NULL,
  hash        CHAR(64)        NOT NULL,
  pins        VARCHAR(1024)   NOT NULL,
  confirmed   BIGINT,
  tx_type     VARCHAR(64)     NOT NULL,
  batch_id    UUID
);

CREATE UNIQUE INDEX messages_id ON messages(id);
CREATE INDEX messages_created ON messages(created);
CREATE INDEX messages_confirmed ON messages(confirmed);
COMMIT;