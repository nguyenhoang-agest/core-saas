# Conventions

## 1. Naming Conventions

- **Module Directories**: Singular, lowercase names (e.g. `billing`, `inventory`, `tenant`).
- **Tables**: Pluralized or snake_cased prefixed by module name where appropriate (e.g. `example_items`, `iam_users`, `tenants`).
- **Go Types**: Standard Go PascalCase for exported types (`Item`, `Service`, `Repository`).
- **Events**: Dot-delimited lowercase strings describing past occurrences (e.g. `tenant.created`, `example.item.created`).

---

## 2. Error Handling

Use `pkg/errors` instead of creating arbitrary unstructured errors:

```go
import "github.com/example/core-backend/pkg/errors"

// In domain / application code:
if item == nil {
    return nil, errors.NotFound("item not found", nil)
}
if opts.Name == "" {
    return nil, errors.InvalidArgument("name is required", nil)
}
```

- When returning errors in HTTP handlers, call `httpx.Error(w, err)`.
- `httpx.Error` automatically inspects `errors.AsAppError` and emits the corresponding HTTP status code (404, 400, 401, 403, 409, 422, 500) and standard JSON payload:
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "item not found"
  }
}
```

---

## 3. Database Migrations

- Migrations live inside each module in `modules/<name>/migrations/`.
- File naming: `00001_description.sql`, `00002_description.sql`.
- Format: Standard Goose annotations (`-- +goose Up`, `-- +goose Down`).
- Each module has its own migration tracking table (`goose_version_<module>`).
- Multi-Tenancy Row Level Security (RLS) template:
```sql
-- +goose Up
ALTER TABLE my_table ENABLE ROW LEVEL SECURITY;

CREATE POLICY my_table_tenant_isolation ON my_table
    FOR ALL
    USING (tenant_id = current_setting('app.current_tenant_id', true))
    WITH CHECK (tenant_id = current_setting('app.current_tenant_id', true));

-- +goose Down
DROP POLICY IF EXISTS my_table_tenant_isolation ON my_table;
```

---

## 4. Events & Pub/Sub

- Asynchronous inter-module messaging uses `pkg/eventbus`.
- Events are plain Go structs passed through typed generics:
```go
// Publishing
eventbus.Publish(ctx, bus, "billing.invoice.paid", InvoicePaidEvent{
    InvoiceID: id,
    TenantID:  tenantID,
})

// Subscribing
eventbus.Subscribe[InvoicePaidEvent](bus, "billing.invoice.paid", func(ctx context.Context, e InvoicePaidEvent) error {
    // React to event
    return nil
})
```

---

## 5. Testing

- Keep module slice tests focused on end-to-end user journeys (HTTP endpoint down to database or in-memory stub).
- Run test suites with race detection enabled:
```bash
make test  # runs go test -race ./...
```
