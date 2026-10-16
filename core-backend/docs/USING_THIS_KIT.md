# Using This SaaS Starter Kit

This starter kit is a modular monolith template written in Go 1.27. It provides common SaaS infrastructure skeletons (`pkg/`), a reference vertical slice (`modules/example`), and stubbed domain modules that can be retained or pruned.

---

## 1. Quick Start

### Step 1: Clone the Template
```bash
git clone <starter-kit-url> my-saas-backend
cd my-saas-backend
```

### Step 2: Configure Environment & Database
Run the automated environment setup helper:

```bash
make env      # Generates .env from .env.example
make db-up    # Starts PostgreSQL 17 container and verifies readiness
# Or simply run: ./scripts/setup-env.sh
```

### Step 3: Run the Application
```bash
make run
```
The monolith server will automatically run all module migrations and listen on `http://localhost:8080`.

Verify health check and database connectivity:
```bash
curl http://localhost:8080/healthz
# Response: {"checks":{"database":"ok"},"status":"ok"}
```

---

## 2. Adding New Domain Modules

When implementing a new feature domain (e.g. `billing`, `orders`, `catalog`):

```bash
make new-module NAME=billing
```

This generator:
1. Copies `modules/_template` to `modules/billing`.
2. Renames Go packages and adjusts imports to your module path.
3. Automatically registers `billing.New(pool, bus)` into `cmd/monolith/modules.go`.

See [ADDING_A_MODULE.md](ADDING_A_MODULE.md) for step-by-step guidance.

---

## 3. Pulling Upstream Updates from the Starter Kit

As the upstream starter kit evolves with bug fixes, security patches, or infrastructure improvements, you can sync your project with upstream without losing your domain work.

### Add Upstream Remote
```bash
git remote add upstream <starter-kit-repository-url>
git fetch upstream
```

### Strategy A: Selective Cherry-Picking (Recommended)
If upstream commits update shared infrastructure (`pkg/`, `deploy/`, CI files):
```bash
git cherry-pick <commit-hash>
```

### Strategy B: Merging Upstream Releases
```bash
git merge upstream/main --no-commit
```

### What to Expect When Conflicts Occur
- **`go.mod`**: Conflicts happen if dependencies are bumped. Keep whichever version is newer or run `go mod tidy` after merging.
- **`cmd/monolith/modules.go`**: Conflicts occur when upstream adds or modifies module registrations. Keep your registered modules and accept upstream changes to shared wiring.
- **`pkg/`**: Should generally be accepted from upstream unless you customized shared infra contracts.
- **`modules/*`**: Your custom business modules won't conflict because upstream does not have them.
