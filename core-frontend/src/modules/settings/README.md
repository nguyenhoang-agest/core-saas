# Settings Module

## Purpose
Enables tenants to customize organization branding (primary colors, logo, theme tokens), configure notification channels, and manage domain aliases.

## Architecture
- `index.ts`: Registers routes (`/settings`), navigation entry, and permissions (`settings:view`, `settings:branding:update`).
- `pages/settings-page.tsx`: Settings dashboard shell.
- `store/settings-store.tsx`: Zustand store for tracking active settings tab and dirty form state.
- `api/settings-api.ts`: API endpoints for tenant metadata and theming.

## TODOs (Domain Implementation)
- [ ] Implement branding customizer form that preview-updates Antd theme tokens in real-time.
- [ ] Connect with backend `/api/v1/tenants/:id/settings` PUT endpoint.
- [ ] Implement custom domain validation and SSL certificate verification UI.
- [ ] Implement webhook notification endpoint settings and secret regeneration.
