# __MODULE_TITLE__ (`__MODULE_NAME__`)

## Purpose
Domain module responsible for __MODULE_DESCRIPTION__.

## Directory Structure
- `index.ts`: Module declaration implementing `ModuleDefinition`.
- `pages/`: Page views exported for Next.js thin route wrappers.
- `components/`: Internal module UI components and slot contributions.
- `api/`: API clients, types, and React Query hooks using `tenantKey`.
- `store/`: Hydration-safe Zustand store created via `createStoreContext`.

## Module Registration
To enable this module, include it in `config/app.config.ts`:
```ts
import { __CAMEL_NAME__Module } from '@/modules/__MODULE_NAME__';

export const appConfig = {
  enabledModules: [
    // ...
    __CAMEL_NAME__Module,
  ],
};
```

## TODOs
- [ ] Define backend API schemas and contracts.
- [ ] Implement data fetching endpoints in `api/`.
- [ ] Construct domain views in `pages/`.
- [ ] Register required slot contributions or menu items in `index.ts`.
- [ ] Write integration and smoke tests.
