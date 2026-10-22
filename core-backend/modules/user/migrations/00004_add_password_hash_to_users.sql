-- +goose Up
-- +goose StatementBegin
ALTER TABLE users 
    ADD COLUMN IF NOT EXISTS password_hash VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users 
    DROP COLUMN IF EXISTS password_hash;
-- +goose StatementEnd
