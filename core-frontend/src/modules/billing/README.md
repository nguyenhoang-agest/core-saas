# Billing Module

## Purpose
Manages commercial subscriptions, billing contacts, payment gateways (Stripe, Paddle, LemonSqueezy), tax IDs, seat additions, and PDF invoice receipts.

## Architecture
- `index.ts`: Registers routes (`/billing`), navigation item, and permissions (`billing:view`, `billing:plans:update`, `billing:invoices:view`).
- `pages/billing-page.tsx`: Billing dashboard view.
- `store/billing-store.tsx`: Zustand store managing active subscription checkout selection and currency filters.
- `api/billing-api.ts`: API client functions for querying current subscription and available plans.

## TODOs (Domain Implementation)
- [ ] Connect with billing backend or Stripe Customer Portal `/api/v1/billing/create-portal-session`.
- [ ] Render interactive pricing comparison table with monthly/annual toggle.
- [ ] Implement Payment Methods manager (Cards, Bank transfers, SEPA).
- [ ] Implement Invoice History data table with one-click PDF receipt downloads.
- [ ] Add webhook listeners on the backend for `invoice.payment_succeeded` and `customer.subscription.deleted`.
