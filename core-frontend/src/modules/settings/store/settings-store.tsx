'use client';

import { createStore } from 'zustand/vanilla';
import { createStoreContext } from '@/core/store/create-store-context';

export interface SettingsState {
  activeTab: string;
  setActiveTab: (tab: string) => void;
}

export function createSettingsStore() {
  return createStore<SettingsState>((set) => ({
    activeTab: 'general',
    setActiveTab: (tab) => set({ activeTab: tab }),
  }));
}

const { Provider: SettingsStoreProvider, useStore: useSettingsStore } =
  createStoreContext(createSettingsStore, 'SettingsStore');

export { SettingsStoreProvider, useSettingsStore };
