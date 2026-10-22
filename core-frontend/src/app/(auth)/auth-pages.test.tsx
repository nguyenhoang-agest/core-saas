import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import LoginPage from './login/page';
import RegisterPage from './register/page';
import { I18nProvider } from '@/core/i18n/i18n-context';
import { signIn } from 'next-auth/react';

const mockPush = vi.fn();
const mockRefresh = vi.fn();

vi.mock('next/navigation', () => ({
  useRouter: () => ({
    push: mockPush,
    refresh: mockRefresh,
  }),
  useSearchParams: () => ({
    get: vi.fn().mockReturnValue(null),
  }),
}));

vi.mock('next-auth/react', () => ({
  signIn: vi.fn(),
  signOut: vi.fn(),
  useSession: () => ({ data: null, status: 'unauthenticated' }),
}));

function renderWithI18n(ui: React.ReactNode, locale: 'en' | 'vi' = 'en') {
  return render(<I18nProvider initialLocale={locale}>{ui}</I18nProvider>);
}

describe('LoginPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders all login form inputs and links', () => {
    renderWithI18n(<LoginPage />);

    expect(screen.getByRole('heading', { name: /Sign in to your account/i })).toBeInTheDocument();
    expect(screen.getByLabelText(/Email Address/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/^Password$/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /Sign In/i })).toBeInTheDocument();
    expect(screen.getByRole('link', { name: /Create a workspace/i })).toHaveAttribute('href', '/register');
  });

  it('toggles password visibility', () => {
    renderWithI18n(<LoginPage />);

    const passwordInput = screen.getByLabelText(/^Password$/i);
    const toggleButton = screen.getByLabelText(/Show password/i);

    expect(passwordInput).toHaveAttribute('type', 'password');

    fireEvent.click(toggleButton);
    expect(passwordInput).toHaveAttribute('type', 'text');

    fireEvent.click(screen.getByLabelText(/Hide password/i));
    expect(passwordInput).toHaveAttribute('type', 'password');
  });

  it('shows error banner when credentials are rejected', async () => {
    vi.mocked(signIn).mockResolvedValueOnce({
      error: 'E010004',
      status: 401,
      ok: false,
      url: null,
    });

    renderWithI18n(<LoginPage />);

    fireEvent.change(screen.getByLabelText(/Email Address/i), {
      target: { value: 'alice@example.com' },
    });
    fireEvent.change(screen.getByLabelText(/^Password$/i), {
      target: { value: 'WrongPassword' },
    });

    fireEvent.click(screen.getByRole('button', { name: /Sign In/i }));

    await waitFor(() => {
      expect(screen.getByRole('alert')).toBeInTheDocument();
      expect(screen.getByText(/Incorrect email or password/i)).toBeInTheDocument();
    });
  });

  it('redirects to root upon successful login', async () => {
    vi.mocked(signIn).mockResolvedValueOnce({
      error: null as any,
      status: 200,
      ok: true,
      url: '/',
    });

    renderWithI18n(<LoginPage />);

    fireEvent.change(screen.getByLabelText(/Email Address/i), {
      target: { value: 'alice@example.com' },
    });
    fireEvent.change(screen.getByLabelText(/^Password$/i), {
      target: { value: 'CorrectPassword' },
    });

    fireEvent.click(screen.getByRole('button', { name: /Sign In/i }));

    await waitFor(() => {
      expect(mockPush).toHaveBeenCalledWith('/');
    });
  });
});

describe('RegisterPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    global.fetch = vi.fn();
  });

  it('renders all organization registration form fields', () => {
    renderWithI18n(<RegisterPage />);

    expect(screen.getByRole('heading', { name: /Create your organization/i })).toBeInTheDocument();
    expect(screen.getByLabelText(/Organization \/ Workspace Name/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Your Full Name/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Work Email/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/^Password$/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/Confirm Password/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /Create Workspace & Account/i })).toBeInTheDocument();
    expect(screen.getByRole('link', { name: /Sign in/i })).toHaveAttribute('href', '/login');
  });

  it('validates password mismatch on client side', async () => {
    renderWithI18n(<RegisterPage />);

    fireEvent.change(screen.getByLabelText(/Organization \/ Workspace Name/i), {
      target: { value: 'Acme Corp' },
    });
    fireEvent.change(screen.getByLabelText(/Your Full Name/i), {
      target: { value: 'Alice' },
    });
    fireEvent.change(screen.getByLabelText(/Work Email/i), {
      target: { value: 'alice@acme.com' },
    });
    fireEvent.change(screen.getByLabelText(/^Password$/i), {
      target: { value: 'Password123' },
    });
    fireEvent.change(screen.getByLabelText(/Confirm Password/i), {
      target: { value: 'DifferentPassword' },
    });

    fireEvent.click(screen.getByRole('button', { name: /Create Workspace & Account/i }));

    await waitFor(() => {
      expect(screen.getByRole('alert')).toBeInTheDocument();
      expect(screen.getByText(/Passwords do not match/i)).toBeInTheDocument();
    });
    expect(global.fetch).not.toHaveBeenCalled();
  });

  it('submits registration and signs in automatically on success', async () => {
    vi.mocked(global.fetch).mockResolvedValueOnce({
      ok: true,
      status: 201,
      json: async () => ({
        code: 'E000000',
        data: {
          tenant_id: 'new-tenant-123',
          user: { id: 'new-user-123', email: 'alice@acme.com' },
        },
      }),
    } as any);

    vi.mocked(signIn).mockResolvedValue({
      error: null as any,
      status: 200,
      ok: true,
      url: '/',
    });

    renderWithI18n(<RegisterPage />);

    fireEvent.change(screen.getByLabelText(/Organization \/ Workspace Name/i), {
      target: { value: 'Acme Corp' },
    });
    fireEvent.change(screen.getByLabelText(/Your Full Name/i), {
      target: { value: 'Alice' },
    });
    fireEvent.change(screen.getByLabelText(/Work Email/i), {
      target: { value: 'alice@acme.com' },
    });
    fireEvent.change(screen.getByLabelText(/^Password$/i), {
      target: { value: 'Password123' },
    });
    fireEvent.change(screen.getByLabelText(/Confirm Password/i), {
      target: { value: 'Password123' },
    });

    fireEvent.click(screen.getByRole('button', { name: /Create Workspace & Account/i }));

    await waitFor(() => {
      expect(global.fetch).toHaveBeenCalledWith(
        '/api/register',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({
            org_name: 'Acme Corp',
            name: 'Alice',
            email: 'alice@acme.com',
            password: 'Password123',
          }),
        })
      );
      expect(signIn).toHaveBeenCalledWith('credentials', expect.objectContaining({
        email: 'alice@acme.com',
        password: 'Password123',
      }));
      expect(mockPush).toHaveBeenCalledWith('/');
    });
  });
});
