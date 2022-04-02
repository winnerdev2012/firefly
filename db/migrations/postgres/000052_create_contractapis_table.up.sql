BEGIN;
CREATE TABLE contractapis (
  seq               SERIAL          PRIMARY KEY,
  id                UUID            NOT NULL,
  interface_id      UUID            NOT NULL,
  location          TEXT,
  name              VARCHAR(64)     NOT NULL,
  namespace         VARCHAR(64)     NOT NULL
);

CREATE UNIQUE INDEX contractapis_namespace_name ON contractapis(namespace,name);
COMMIT;