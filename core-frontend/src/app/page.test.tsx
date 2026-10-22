import { render, screen } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import HomePage from './page';

let mockStatus: 'authenticated' | 'unauthenticated' = 'unauthenticated';
let mockIsHost = false;

vi.mock('next-auth/react', () => ({
  useSession: () => ({
    data:
      mockStatus === 'authenticated'
        ? {
            user: {
              id: 'user-1',
              name: 'Demo Admin',
              email: 'admin@example.com',
              isHost: mockIsHost,
            },
          }
        : null,
    status: mockStatus,
  }),
}));

describe('HomePage', () => {
  beforeEach(() => {
    mockStatus = 'unauthenticated';
    mockIsHost = false;
  });

  it('renders the core saas admin heading', () => {
    render(<HomePage />);
    expect(screen.getByText('Modular Multi-Tenant SaaS Console')).toBeInTheDocument();
  });

  it('renders guest welcome actions and preview indicator for first-time / unauthenticated visitor', () => {
    mockStatus = 'unauthenticated';
    render(<HomePage />);

    // Calls to action for sign in / register
    expect(screen.getByText('Sign In to Console')).toBeInTheDocument();
    expect(screen.getByText('Create Organization')).toBeInTheDocument();

    // Guest mode indicator
    expect(screen.getByText('Public Preview • Guest')).toBeInTheDocument();
    expect(
      screen.getByText('(Sign in required to access management consoles)')
    ).toBeInTheDocument();

    // Module cards prompt sign in
    const signInLinks = screen.getAllByText('Sign In to Access');
    expect(signInLinks.length).toBe(4);
  });

  it('renders management actions and active workspace for authenticated tenant owner', () => {
    mockStatus = 'authenticated';
    mockIsHost = false;
    render(<HomePage />);

    expect(screen.getByText('Manage Users & Roles')).toBeInTheDocument();
    expect(screen.getByText('Workspace Settings')).toBeInTheDocument();
    expect(screen.queryByText('Manage Tenants')).not.toBeInTheDocument();

    expect(screen.getByText('Active Workspace:')).toBeInTheDocument();
    const openLinks = screen.getAllByText('Open Module');
    expect(openLinks.length).toBe(4);
  });

  it('renders platform tenant management action for authenticated host admin', () => {
    mockStatus = 'authenticated';
    mockIsHost = true;
    render(<HomePage />);

    expect(screen.getByText('Manage Tenants')).toBeInTheDocument();
    expect(screen.getByText('Manage Users & Roles')).toBeInTheDocument();
  });
});
