# Core SaaS Starter Kit

A production-grade, modular SaaS starter kit built on Next.js 16 (App Router), Ant Design v6, React Query v5, Zustand v5, NextAuth v4, and Vitest.

Designed as a single application repository with strict modular boundaries. Clone it, rename it with `pnpm init-kit --name <project>`, and add domain modules with `pnpm new-module <name>`.

---

## Tech Stack

- **Framework**: Next.js 16.4 (Turbopack, App Router, Edge Proxy request interception)
- **UI Library**: Ant Design v6 with `@ant-design/nextjs-registry` and `@ant-design/cssinjs`
- **Data Fetching**: TanStack React Query v5 with automatic tenant-scoped query cache keys
- **State Management**: Zustand v5 with SSR-safe store factory & React Context pattern
- **Authentication**: NextAuth v4 with JWT strategy, session tokens, and `<Can>` UI permission guards
- **Validation & Env**: Zod v3 environment validation
- **Testing**: Vitest v3 with jsdom & React Testing Library
- **Architecture**: Modular domain slices (`src/modules/*`) with ESLint boundary enforcement

---

## Getting Started

### 1. Environment Configuration
From repository root:
```bash
make env
# Or manually inside core-frontend:
cp .env.example .env.local
```

### 2. Install & Run
```bash
npm install
npm run dev
# Or from repository root:
make dev-frontend
```
Open [http://localhost:3000](http://localhost:3000) in your browser.

---

## UI & Design System

The frontend strictly implements an **Outlined Minimalist / Flat 2.0 Monochrome UI**:
- **Border-focused hierarchy**: Thin 1px borders (`border-neutral-200`) instead of heavy drop shadows (`shadow-none`).
- **Monochrome Grayscale**: 100% grayscale palette (`bg-white`, `bg-neutral-50`, `text-neutral-900`), reserving status colors strictly for live indicators (e.g. green dot for active status).
- **Internationalization (i18n)**:
  - English (`en`) and Vietnamese (`vi`) with parameter interpolation: `useI18n().t('key', { param })`.
  - Two-tier error resolution: Client dictionary lookup $\rightarrow$ Backend returned message $\rightarrow$ Fallback.
  - Minimalist `LocaleSwitcher` on the topbar.

---

## Creating Modules

```bash
# Scaffold a new domain module
pnpm new-module analytics
# or: npm run new-module analytics
```

Modules are self-contained in `src/modules/<name>/` and declared in `config/app.config.ts`.

---

## Documentation

- [**Architecture Overview**](./docs/ARCHITECTURE.md): Deep-dive into core subsystems, extension points, slots, and multi-tenancy.
- [**Adding a Module**](./docs/ADDING_A_MODULE.md): Guide for creating, implementing, and registering domain modules.
- [**Using This Kit**](./docs/USING_THIS_KIT.md): Guide on cloning, configuration, and pulling upstream updates via Git remotes.

---

## Quality Verification

```bash
npx tsc --noEmit   # Strict TypeScript check
npm run lint       # ESLint & Architectural boundary checks
npm test           # Vitest unit & integration test suite
npm run build      # Production Next.js build
```
