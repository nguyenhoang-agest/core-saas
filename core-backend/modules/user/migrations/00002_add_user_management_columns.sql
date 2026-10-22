-- +goose Up
-- +goose StatementBegin
ALTER TABLE users 
    ADD COLUMN IF NOT EXISTS tenant_id VARCHAR(64) NOT NULL DEFAULT 'default',
    ADD COLUMN IF NOT EXISTS email VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS role VARCHAR(50) NOT NULL DEFAULT 'member',
    ADD COLUMN IF NOT EXISTS permissions TEXT[] NOT NULL DEFAULT '{}',
    ADD COLUMN IF NOT EXISTS phone VARCHAR(50),
    ADD COLUMN IF NOT EXISTS avatar_url VARCHAR(512),
    ADD COLUMN IF NOT EXISTS last_login_at TIMESTAMPTZ;

ALTER TABLE users DROP COLUMN IF EXISTS description;

CREATE INDEX IF NOT EXISTS idx_users_tenant ON users(tenant_id);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'uq_users_tenant_email'
    ) THEN
        ALTER TABLE users ADD CONSTRAINT uq_users_tenant_email UNIQUE(tenant_id, email);
    END IF;
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP CONSTRAINT IF EXISTS uq_users_tenant_email;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_tenant;
ALTER TABLE users 
    DROP COLUMN IF EXISTS last_login_at,
    DROP COLUMN IF EXISTS avatar_url,
    DROP COLUMN IF EXISTS phone,
    DROP COLUMN IF EXISTS permissions,
    DROP COLUMN IF EXISTS role,
    DROP COLUMN IF EXISTS email,
    DROP COLUMN IF EXISTS tenant_id;
-- +goose StatementEnd
