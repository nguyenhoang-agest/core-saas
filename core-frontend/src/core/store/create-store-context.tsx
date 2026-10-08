'use client';

import React, { createContext, useContext, useRef } from 'react';
import { useStore, type StoreApi } from 'zustand';

export interface StoreContextResult<TStore, TProps> {
  Provider: React.FC<{ children: React.ReactNode; initialValues?: TProps }>;
  useStore: <TSelected = TStore>(selector?: (state: TStore) => TSelected) => TSelected;
  StoreContext: React.Context<StoreApi<TStore> | null>;
}

/**
 * Creates a Zustand store factory integrated with React Context.
 * Ensures SSR safety by instantiating the store per-request/per-tree
 * instead of relying on global singletons.
 */
export function createStoreContext<TStore, TProps = Record<string, unknown>>(
  storeCreator: (initProps?: TProps) => StoreApi<TStore>,
  contextName = 'Store'
): StoreContextResult<TStore, TProps> {
  const StoreContext = createContext<StoreApi<TStore> | null>(null);

  function Provider({
    children,
    initialValues,
  }: {
    children: React.ReactNode;
    initialValues?: TProps;
  }) {
    const storeRef = useRef<StoreApi<TStore> | null>(null);
    if (!storeRef.current) {
      storeRef.current = storeCreator(initialValues);
    }

    return <StoreContext.Provider value={storeRef.current}>{children}</StoreContext.Provider>;
  }

  function useStoreHook<TSelected = TStore>(
    selector: (state: TStore) => TSelected = (s) => s as unknown as TSelected
  ): TSelected {
    const store = useContext(StoreContext);
    if (!store) {
      throw new Error(`[${contextName}] Hook must be used within its corresponding Provider.`);
    }
    return useStore(store, selector);
  }

  return {
    Provider,
    useStore: useStoreHook,
    StoreContext,
  };
}
