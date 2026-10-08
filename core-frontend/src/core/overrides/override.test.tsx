import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import React from 'react';
import { OverrideProvider, Override, useOverride } from './override-context';

describe('Component Override Map', () => {
  const DefaultLogo: React.FC<{ title: string }> = ({ title }) => (
    <div data-testid="default-logo">{title} (Default)</div>
  );

  const CustomLogo: React.FC<{ title: string }> = ({ title }) => (
    <div data-testid="custom-logo">{title} (Custom)</div>
  );

  it('renders default component when no override is provided', () => {
    render(<Override name="Logo" defaultComponent={DefaultLogo} props={{ title: 'CoreSaaS' }} />);
    expect(screen.getByTestId('default-logo')).toHaveTextContent('CoreSaaS (Default)');
  });

  it('renders custom component when override is registered in provider', () => {
    render(
      <OverrideProvider overrides={{ Logo: CustomLogo }}>
        <Override name="Logo" defaultComponent={DefaultLogo} props={{ title: 'CoreSaaS' }} />
      </OverrideProvider>
    );
    expect(screen.getByTestId('custom-logo')).toHaveTextContent('CoreSaaS (Custom)');
    expect(screen.queryByTestId('default-logo')).toBeNull();
  });
});
