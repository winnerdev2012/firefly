BEGIN;
CREATE TABLE blockchainevents (
  seq              INTEGER         PRIMARY KEY AUTOINCREMENT,
  id               UUID            NOT NULL,
  source           VARCHAR(256)    NOT NULL,
  namespace        VARCHAR(64)     NOT NULL,
  name             VARCHAR(256)    NOT NULL,
  protocol_id      VARCHAR(256)    NOT NULL,
  timestamp        BIGINT          NOT NULL,
  subscription_id  UUID,
  output           BYTEA,
  info             BYTEA,
  tx_type          VARCHAR(64),
  tx_id            UUID
);
COMMIT;
