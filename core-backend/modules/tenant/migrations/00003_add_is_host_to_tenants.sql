-- +goose Up
-- +goose StatementBegin
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS is_host BOOLEAN NOT NULL DEFAULT FALSE;
CREATE INDEX IF NOT EXISTS idx_tenants_is_host ON tenants(is_host);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_tenants_is_host;
ALTER TABLE tenants DROP COLUMN IF EXISTS is_host;
-- +goose StatementEnd
