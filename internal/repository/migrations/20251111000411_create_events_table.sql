-- +goose Up
-- +goose StatementBegin
CREATE TABLE users_read_model (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  email TEXT NOT NULL, 
  password TEXT NOT NULL,

  UNIQUE(email)
);

CREATE TABLE users_events (
  id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  user_id UUID NOT NULL REFERENCES users_read_model(id),
  event_type TEXT NOT NULL,
  event_data JSONB NOT NULL,
  version SMALLINT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),

  UNIQUE(user_id, version)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users_read_model;
DROP TABLE users_events;
-- +goose StatementEnd
