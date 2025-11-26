-- +goose Up
-- +goose StatementBegin
CREATE TABLE addresses(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  organization_id UUID NOT NULL REFERENCES organizations(id),
  zip_code TEXT NOT NULL,
  complement TEXT NOT NULL,
  public_place TEXT NOT NULL,
  city TEXT NOT NULL,
  state TEXT NOT NULL,
  UNIQUE(organization_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE addresses;
-- +goose StatementEnd
