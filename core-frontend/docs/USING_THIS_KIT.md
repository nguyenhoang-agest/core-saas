# Using This Starter Kit

This guide walks you through using the SaaS Starter Kit: from cloning and initializing a new project to managing modules and pulling upstream kit updates over time.

---

## 1. Quick Start

### Step 1: Clone the Template Repository
```bash
git clone <starter-kit-repository-url> my-new-saas
cd my-new-saas
```

### Step 2: Initialize Your Project
Run the automated initialization script to rename the package, set branding strings, and update environment defaults:

```bash
# Using pnpm
pnpm init-kit --name my-new-saas

# Using npm
npm run init-kit -- --name my-new-saas
```

### Step 3: Configure Environment Variables
Copy `.env.example` to `.env.local` and adjust any secrets:
```bash
cp .env.example .env.local
```

### Step 4: Install Dependencies & Run
```bash
# Install dependencies
pnpm install # or npm install

# Start development server
pnpm dev     # or npm run dev
```

Visit [http://localhost:3000](http://localhost:3000) to view your new SaaS dashboard.

---

## 2. Managing Modules

### Adding a New Module
To scaffold a new domain module with isolated state, API, routing, and smoke tests:
```bash
pnpm new-module analytics
# or: npm run new-module analytics
```
See [ADDING_A_MODULE.md](./ADDING_A_MODULE.md) for detailed implementation instructions.

### Enabling or Disabling Modules
Modules are explicitly declared in `config/app.config.ts`. To enable or disable a module, adjust the `enabledModules` list:

```ts
import { exampleModule } from '@/modules/example';
import { iamModule } from '@/modules/iam';
import { billingModule } from '@/modules/billing';

export const appConfig = {
  enabledModules: [
    exampleModule,
    iamModule,
    // billingModule, // Comment out to disable billing
  ],
};
```
When a module is removed from `enabledModules`, its navigation menu items, route bindings, and slot contributions are immediately detached from the application shell.

---

## 3. Pulling Upstream Starter Kit Updates

Because this starter kit is maintained as a clean single application repository (no monorepo complexity, no private npm publishing overhead), downstream projects can sync upstream improvements via Git remotes.

### Step 1: Configure the Upstream Remote
In your cloned downstream repository, add the original kit repository as `upstream`:

```bash
git remote add upstream <starter-kit-repository-url>
```

Verify your remotes:
```bash
git remote -v
# origin   https://github.com/my-org/my-new-saas.git (fetch & push)
# upstream <starter-kit-repository-url> (fetch & push)
```

### Step 2: Fetch and Merge Upstream Changes
When bug fixes, performance optimizations, or framework updates are published to the starter kit:

```bash
# Fetch latest updates from the kit
git fetch upstream

# Merge upstream updates into your working branch
git merge upstream/main -m "Merge upstream starter kit updates"
```

### Why This Architecture Stays Clean
- **Clean Boundaries**: All core infrastructure lives in `src/core/`, while your domain business logic lives inside `src/modules/<your-modules>/`.
- **Zero Infiltration**: Upstream commits primarily modify `src/core/`, `proxy.ts`, or documentation, resulting in virtually no merge conflicts with your domain modules.
- **Easy Maintenance**: After merging, run the verification suite to ensure all tests pass:
  ```bash
  npx tsc --noEmit && npm run lint && npm test && npm run build
  ```
