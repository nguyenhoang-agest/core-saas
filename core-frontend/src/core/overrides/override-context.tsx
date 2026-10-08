'use client';

import React, { createContext, useContext, useMemo } from 'react';

export type ComponentMap = Record<string, React.ComponentType<any>>;

interface OverrideContextValue {
  overrides: ComponentMap;
}

const OverrideContext = createContext<OverrideContextValue>({
  overrides: {},
});

export interface OverrideProviderProps {
  children: React.ReactNode;
  overrides?: ComponentMap;
}

export function OverrideProvider({ children, overrides = {} }: OverrideProviderProps) {
  const value = useMemo(() => ({ overrides }), [overrides]);

  return <OverrideContext.Provider value={value}>{children}</OverrideContext.Provider>;
}

/**
 * Hook to retrieve an overridden component or fallback to the provided default.
 */
export function useOverride<P extends Record<string, unknown>>(
  name: string,
  DefaultComponent: React.ComponentType<P>
): React.ComponentType<P> {
  const context = useContext(OverrideContext);
  return (context.overrides[name] as React.ComponentType<P>) || DefaultComponent;
}

/**
 * Helper component rendering either the override or default component with passed props.
 */
export function Override<P extends Record<string, unknown>>({
  name,
  defaultComponent,
  props,
}: {
  name: string;
  defaultComponent: React.ComponentType<P>;
  props?: P;
}) {
  const Component = useOverride(name, defaultComponent);
  return React.createElement(Component, (props ?? {}) as P);
}
