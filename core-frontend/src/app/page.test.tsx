import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import HomePage from './page';

describe('HomePage', () => {
  it('renders the core saas admin heading', () => {
    render(<HomePage />);
    expect(screen.getByText('Modular Multi-Tenant SaaS Console')).toBeInTheDocument();
  });

  it('renders the manage tenants and users actions', () => {
    render(<HomePage />);
    expect(screen.getByText('Manage Tenants')).toBeInTheDocument();
    expect(screen.getByText('Manage Users & Roles')).toBeInTheDocument();
    expect(screen.getByText('Tenant Workspaces')).toBeInTheDocument();
    expect(screen.getByText('Users & Access Control')).toBeInTheDocument();
  });

  it('renders active workspace context indicator', () => {
    render(<HomePage />);
    expect(screen.getByText('Active Workspace:')).toBeInTheDocument();
    expect(screen.getByText('Architecture & Security Engineering Pillars')).toBeInTheDocument();
  });
});
