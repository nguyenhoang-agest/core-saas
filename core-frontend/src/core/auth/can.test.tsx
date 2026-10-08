import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import React from 'react';
import { Can } from './can';
import * as nextAuthReact from 'next-auth/react';

vi.mock('next-auth/react', () => ({
  useSession: vi.fn(),
}));

describe('Can component', () => {
  it('renders fallback when session has no matching permission', () => {
    vi.mocked(nextAuthReact.useSession).mockReturnValue({
      data: {
        user: { permissions: ['billing:read'] },
      } as any,
      status: 'authenticated',
      update: vi.fn(),
    });

    render(
      <Can permission="billing:write" fallback={<div>Access Denied</div>}>
        <div>Sensitive Action</div>
      </Can>
    );

    expect(screen.getByText('Access Denied')).toBeInTheDocument();
    expect(screen.queryByText('Sensitive Action')).toBeNull();
  });

  it('renders children when session has matching permission', () => {
    vi.mocked(nextAuthReact.useSession).mockReturnValue({
      data: {
        user: { permissions: ['billing:write'] },
      } as any,
      status: 'authenticated',
      update: vi.fn(),
    });

    render(
      <Can permission="billing:write" fallback={<div>Access Denied</div>}>
        <div>Sensitive Action</div>
      </Can>
    );

    expect(screen.getByText('Sensitive Action')).toBeInTheDocument();
    expect(screen.queryByText('Access Denied')).toBeNull();
  });

  it('renders children when user has wildcard (*) permission', () => {
    vi.mocked(nextAuthReact.useSession).mockReturnValue({
      data: {
        user: { permissions: ['*'] },
      } as any,
      status: 'authenticated',
      update: vi.fn(),
    });

    render(
      <Can permission="settings:delete">
        <div>Admin Button</div>
      </Can>
    );

    expect(screen.getByText('Admin Button')).toBeInTheDocument();
  });
});
