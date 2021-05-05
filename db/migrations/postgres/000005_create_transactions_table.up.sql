BEGIN;
CREATE SEQUENCE transactions_seq;
CREATE TABLE transactions (
  id          CHAR(36)        NOT NULL PRIMARY KEY,
  seq         BIGINT          NOT NULL DEFAULT nextval('transactions_seq'),
  ttype       VARCHAR(64)     NOT NULL,
  author      VARCHAR(1024)   NOT NULL,
  created     BIGINT          NOT NULL,
  tracking_id VARCHAR(256),
  protocol_id VARCHAR(256),
  confirmed   BIGINT          NOT NULL,
  info        JSONB
);

CREATE INDEX transactions_search ON transactions(ttype,author,confirmed,created);
CREATE INDEX transactions_tracking_id ON transactions(tracking_id);
CREATE INDEX transactions_protocol_id ON transactions(protocol_id);
COMMIT;