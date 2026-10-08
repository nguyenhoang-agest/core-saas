import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import HomePage from './page';

describe('HomePage', () => {
  it('renders the core saas frontend heading', () => {
    render(<HomePage />);
    expect(screen.getByText('Core SaaS Frontend')).toBeInTheDocument();
  });

  it('renders the tech stack tags', () => {
    render(<HomePage />);
    expect(screen.getAllByText(/Next.js/i).length).toBeGreaterThan(0);
    expect(screen.getAllByText(/Ant Design/i).length).toBeGreaterThan(0);
  });
});
