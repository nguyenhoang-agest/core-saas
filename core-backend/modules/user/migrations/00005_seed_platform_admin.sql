-- +goose Up
-- +goose StatementBegin
INSERT INTO users (id, tenant_id, name, email, role, status, password_hash, create_at, update_at)
VALUES (
    'admin-root-default',
    'default',
    'Platform Administrator',
    'admin@example.com',
    'owner',
    'active',
    '$2b$12$nomJzk29KC8iyFZ.bwRf4OEXIO.E9xJbRfVhFPu2si4Zf.u86NCR2',
    NOW(),
    NOW()
)
ON CONFLICT (id) DO UPDATE SET
    tenant_id = 'default',
    email = 'admin@example.com',
    role = 'owner',
    status = 'active',
    password_hash = '$2b$12$nomJzk29KC8iyFZ.bwRf4OEXIO.E9xJbRfVhFPu2si4Zf.u86NCR2',
    update_at = NOW();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE id = 'admin-root-default';
-- +goose StatementEnd
