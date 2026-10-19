# User Module

## Purpose
Canonical blueprint for backend domain modules following the **Tenant Pattern & Structure**.

---

## Layout
- `module.go`: Module lifecycle entrypoint implementing `kernel.Module`. Wires GORM `*gorm.DB`, repository, service, and Fiber handler.
- `api/`: Public contracts, DTOs, and event topics (`events.go`).
- `internal/model/`: Domain entity models with GORM tags and `BaseModel` (`base_model.go`).
- `internal/dto/`: Request/Response data transfer objects with validator tags.
- `internal/errcode/`: Standardized machine-readable error codes.
- `internal/repository/`: Data access implementations wrapping GORM `*gorm.DB`.
- `internal/service/`: Business domain orchestration and eventbus publishing.
- `internal/handler/`: Fiber v3 HTTP handlers and route registration.
- `migrations/`: Embedded Goose SQL migrations.
