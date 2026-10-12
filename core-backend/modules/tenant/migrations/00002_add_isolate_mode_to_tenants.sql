-- +goose Up
-- +goose StatementBegin
ALTER TABLE tenants 
    ADD COLUMN IF NOT EXISTS isolate_mode VARCHAR(50) NOT NULL DEFAULT 'shared',
    ADD COLUMN IF NOT EXISTS db_config JSONB;

CREATE INDEX IF NOT EXISTS idx_tenants_isolate_mode ON tenants(isolate_mode);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_tenants_isolate_mode;
ALTER TABLE tenants 
    DROP COLUMN IF EXISTS db_config,
    DROP COLUMN IF EXISTS isolate_mode;
-- +goose StatementEnd
