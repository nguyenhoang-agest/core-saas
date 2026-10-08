'use client';

import React from 'react';
import { moduleRegistry } from '../modules/registry';

export interface SlotProps<P = Record<string, unknown>> {
  name: string;
  props?: P;
  fallback?: React.ReactNode;
  className?: string;
  as?: React.ElementType;
}

export function Slot<P extends Record<string, unknown> = Record<string, unknown>>({
  name,
  props,
  fallback = null,
  className,
  as: Container = React.Fragment,
}: SlotProps<P>) {
  const contributions = moduleRegistry.getSlotContributions(name);

  if (contributions.length === 0) {
    return <>{fallback}</>;
  }

  const renderedContributions = contributions.map((contribution) =>
    React.createElement(contribution.component as React.ComponentType<P>, {
      key: contribution.id,
      ...((props ?? {}) as P),
    })
  );

  if (Container === React.Fragment) {
    return <>{renderedContributions}</>;
  }

  return <Container className={className}>{renderedContributions}</Container>;
}
