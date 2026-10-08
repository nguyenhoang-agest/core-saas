'use client';

import { createStore } from 'zustand/vanilla';
import { createStoreContext } from '@/core/store/create-store-context';

export interface ExampleState {
  filter: string;
  counter: number;
  setFilter: (filter: string) => void;
  increment: () => void;
  reset: () => void;
}

export interface ExampleStoreInitProps {
  initialFilter?: string;
  initialCounter?: number;
}

export function createExampleStore(initProps?: ExampleStoreInitProps) {
  return createStore<ExampleState>((set) => ({
    filter: initProps?.initialFilter ?? '',
    counter: initProps?.initialCounter ?? 0,
    setFilter: (filter) => set({ filter }),
    increment: () => set((state) => ({ counter: state.counter + 1 })),
    reset: () => set({ filter: '', counter: 0 }),
  }));
}

const {
  Provider: ExampleStoreProvider,
  useStore: useExampleStore,
} = createStoreContext(createExampleStore, 'ExampleStore');

export { ExampleStoreProvider, useExampleStore };
