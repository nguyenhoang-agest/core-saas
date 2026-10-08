-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS example_items (
    id VARCHAR(64) PRIMARY KEY,
    tenant_id VARCHAR(64) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_example_items_tenant_id ON example_items(tenant_id);

ALTER TABLE example_items ENABLE ROW LEVEL SECURITY;

CREATE POLICY example_items_tenant_isolation ON example_items
    FOR ALL
    USING (tenant_id = current_setting('app.current_tenant_id', true))
    WITH CHECK (tenant_id = current_setting('app.current_tenant_id', true));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP POLICY IF EXISTS example_items_tenant_isolation ON example_items;
DROP TABLE IF EXISTS example_items;
-- +goose StatementEnd
