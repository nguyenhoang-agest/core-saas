-- +goose Up
-- +goose StatementBegin
INSERT INTO tenants (id, name, description, status, isolate_mode, is_host, create_at, update_at)
VALUES ('default', 'Platform Host', 'Platform root administration workspace', 'active', 'shared', TRUE, NOW(), NOW())
ON CONFLICT (id) DO UPDATE SET is_host = TRUE, update_at = NOW();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM tenants WHERE id = 'default' AND is_host = TRUE;
-- +goose StatementEnd
