import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import HomePage from './page';

describe('HomePage', () => {
  it('renders the core saas admin heading', () => {
    render(<HomePage />);
    expect(screen.getByText('Modular Multi-Tenant SaaS Console')).toBeInTheDocument();
  });

  it('renders the manage tenants action', () => {
    render(<HomePage />);
    expect(screen.getByText('Manage Tenants')).toBeInTheDocument();
    expect(screen.getByText('Tenant Workspaces')).toBeInTheDocument();
  });
});
