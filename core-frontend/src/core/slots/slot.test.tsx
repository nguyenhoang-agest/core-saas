import { render, screen } from '@testing-library/react';
import { describe, it, expect, beforeEach } from 'vitest';
import React from 'react';
import { Slot } from './slot';
import { moduleRegistry } from '../modules/registry';

describe('Slot component', () => {
  beforeEach(() => {
    moduleRegistry.clear();
  });

  it('renders fallback when no contributions are registered', () => {
    render(<Slot name="empty-slot" fallback={<div>Empty Fallback</div>} />);
    expect(screen.getByText('Empty Fallback')).toBeInTheDocument();
  });

  it('renders nothing if no contributions and no fallback', () => {
    const { container } = render(<Slot name="empty-slot" />);
    expect(container.firstChild).toBeNull();
  });

  it('renders registered slot contributions in order with props', () => {
    const ItemA: React.FC<{ label?: string }> = ({ label }) => (
      <div data-testid="item-a">Item A: {label}</div>
    );
    const ItemB: React.FC<{ label?: string }> = ({ label }) => (
      <div data-testid="item-b">Item B: {label}</div>
    );

    moduleRegistry.register({
      id: 'test-module-1',
      name: 'Test Module 1',
      slots: [
        {
          id: 'slot-b',
          slotName: 'header-actions',
          order: 20,
          component: ItemB,
        },
        {
          id: 'slot-a',
          slotName: 'header-actions',
          order: 10,
          component: ItemA,
        },
      ],
    });

    render(<Slot name="header-actions" props={{ label: 'Shared' }} />);

    expect(screen.getByTestId('item-a')).toHaveTextContent('Item A: Shared');
    expect(screen.getByTestId('item-b')).toHaveTextContent('Item B: Shared');

    // Verify ordering
    const elements = screen.getAllByTestId(/^item-/);
    expect(elements[0]).toHaveAttribute('data-testid', 'item-a');
    expect(elements[1]).toHaveAttribute('data-testid', 'item-b');
  });
});
