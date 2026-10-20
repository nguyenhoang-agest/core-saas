'use client';

import { createStore } from 'zustand/vanilla';
import { createStoreContext } from '@/core/store/create-store-context';

export interface UserState {
  selectedUserId: string | null;
  setSelectedUserId: (id: string | null) => void;
  activeTab: 'users' | 'matrix';
  setActiveTab: (tab: 'users' | 'matrix') => void;
}

export function createUserStore() {
  return createStore<UserState>((set) => ({
    selectedUserId: null,
    setSelectedUserId: (id) => set({ selectedUserId: id }),
    activeTab: 'users',
    setActiveTab: (tab) => set({ activeTab: tab }),
  }));
}

const { Provider: UserStoreProvider, useStore: useUserStore } =
  createStoreContext(createUserStore, 'UserStore');

export { UserStoreProvider, useUserStore };
