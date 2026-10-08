'use client';

import { createStore } from 'zustand/vanilla';
import { createStoreContext } from '@/core/store/create-store-context';

export interface TemplateState {
  isInitialized: boolean;
  setInitialized: (val: boolean) => void;
}

export function createTemplateStore() {
  return createStore<TemplateState>((set) => ({
    isInitialized: true,
    setInitialized: (val) => set({ isInitialized: val }),
  }));
}

const { Provider: TemplateStoreProvider, useStore: useTemplateStore } =
  createStoreContext(createTemplateStore, 'TemplateStore');

export { TemplateStoreProvider, useTemplateStore };
