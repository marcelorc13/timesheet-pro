-- +goose Up
-- +goose StatementBegin
CREATE TABLE organization_roles(
  id SMALLINT PRIMARY KEY,
  name TEXT NOT NULL
);
INSERT INTO organization_roles(id, name) VALUES (1, 'member'), (2, 'admin');

CREATE TABLE organization_users(
  id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  user_id UUID NOT NULL REFERENCES users(id),
  organization_id UUID NOT NULL REFERENCES organizations(id),
  organization_role_id SMALLINT NOT NULL REFERENCES organization_roles(id),
  joined_at TIMESTAMPTZ DEFAULT NOW(),

  UNIQUE(user_id, organization_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE organization_users;
DROP TABLE organization_roles;
-- +goose StatementEnd
