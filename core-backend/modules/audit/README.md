# Audit Module (STUB)

## Purpose
Records structured audit events and security trails for system actions across all tenants.

## Planned Ports & Interfaces
- **Ports (Domain)**:
  - `Repository`: Append-only storage for immutable audit log entries.
- **Adapters**:
  - `PostgresAuditRepository`: Inserts into partitioned or indexed PostgreSQL table.
- **Events Subscribed To**:
  - `*`: Listens to security-critical domain events (e.g. `iam.user.authenticated`, `tenant.created`) to write audit trails.

## TODO Checklist
- [ ] Implement append-only insertion in `PostgresAuditRepository`
- [ ] Connect audit subscriber to `pkg/eventbus` wildcard or specific security topics
- [ ] Add query filters for action, actor, and time ranges in HTTP handler
- [ ] Evaluate cold-storage archival strategy (S3 / Parquet) for high-volume logs
