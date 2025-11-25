-- +goose Up
-- +goose StatementBegin
ALTER TABLE organization_users 
ADD CONSTRAINT organization_users_user_id_unique UNIQUE (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE organization_users 
DROP CONSTRAINT IF EXISTS organization_users_user_id_unique;
-- +goose StatementEnd
