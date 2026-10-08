'use client';

import { createStore } from 'zustand/vanilla';
import { createStoreContext } from '@/core/store/create-store-context';

export interface IamState {
  selectedUserId: string | null;
  setSelectedUserId: (id: string | null) => void;
}

export function createIamStore() {
  return createStore<IamState>((set) => ({
    selectedUserId: null,
    setSelectedUserId: (id) => set({ selectedUserId: id }),
  }));
}

const { Provider: IamStoreProvider, useStore: useIamStore } =
  createStoreContext(createIamStore, 'IamStore');

export { IamStoreProvider, useIamStore };
