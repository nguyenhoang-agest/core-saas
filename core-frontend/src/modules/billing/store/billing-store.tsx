'use client';

import { createStore } from 'zustand/vanilla';
import { createStoreContext } from '@/core/store/create-store-context';

export interface BillingState {
  selectedPlan: string;
  setSelectedPlan: (plan: string) => void;
}

export function createBillingStore() {
  return createStore<BillingState>((set) => ({
    selectedPlan: 'starter',
    setSelectedPlan: (plan) => set({ selectedPlan: plan }),
  }));
}

const { Provider: BillingStoreProvider, useStore: useBillingStore } =
  createStoreContext(createBillingStore, 'BillingStore');

export { BillingStoreProvider, useBillingStore };
