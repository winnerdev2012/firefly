CREATE TABLE subscriptions (
  id             string          NOT NULL,
  namespace      string          NOT NULL,
  name           string          NOT NULL,
  dispatcher     string          NOT NULL,
  events         string          NOT NULL,
  filter_topic   string          NOT NULL,
  filter_context string          NOT NULL,
  filter_group   string          NOT NULL,
  options        blob            NOT NULL,
  created        int64           NOT NULL
);

CREATE UNIQUE INDEX subscriptions_primary ON subscriptions(id);
CREATE UNIQUE INDEX subscriptions_name ON subscriptions(namespace,name);
