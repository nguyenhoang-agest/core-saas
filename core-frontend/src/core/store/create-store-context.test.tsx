import { render, screen, act } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import React from 'react';
import { createStore } from 'zustand/vanilla';
import { createStoreContext } from './create-store-context';

interface CounterState {
  count: number;
  increment: () => void;
}

function createCounterStore(initProps?: { initialCount?: number }) {
  return createStore<CounterState>((set) => ({
    count: initProps?.initialCount ?? 0,
    increment: () => set((state) => ({ count: state.count + 1 })),
  }));
}

const { Provider: CounterProvider, useStore: useCounterStore } = createStoreContext(
  createCounterStore,
  'CounterStore'
);

function CounterDisplay() {
  const count = useCounterStore((s) => s.count);
  const increment = useCounterStore((s) => s.increment);

  return (
    <div>
      <span data-testid="count">{count}</span>
      <button onClick={increment} data-testid="btn-increment">
        Increment
      </button>
    </div>
  );
}

describe('createStoreContext', () => {
  it('throws error when hook is called outside Provider', () => {
    expect(() => render(<CounterDisplay />)).toThrowError(
      /CounterStore.*Hook must be used within its corresponding Provider/
    );
  });

  it('renders initial state and responds to state updates', () => {
    render(
      <CounterProvider initialValues={{ initialCount: 10 }}>
        <CounterDisplay />
      </CounterProvider>
    );

    expect(screen.getByTestId('count')).toHaveTextContent('10');

    act(() => {
      screen.getByTestId('btn-increment').click();
    });

    expect(screen.getByTestId('count')).toHaveTextContent('11');
  });

  it('isolates state between separate Provider instances (no singleton leak)', () => {
    render(
      <div>
        <CounterProvider initialValues={{ initialCount: 5 }}>
          <div data-testid="tree-1">
            <CounterDisplay />
          </div>
        </CounterProvider>
        <CounterProvider initialValues={{ initialCount: 100 }}>
          <div data-testid="tree-2">
            <CounterDisplay />
          </div>
        </CounterProvider>
      </div>
    );

    const tree1 = screen.getByTestId('tree-1');
    const tree2 = screen.getByTestId('tree-2');

    expect(tree1.querySelector('[data-testid="count"]')).toHaveTextContent('5');
    expect(tree2.querySelector('[data-testid="count"]')).toHaveTextContent('100');
  });
});
