import React from 'react';
import { render, screen, fireEvent, renderHook, act } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { I18nProvider, useI18n } from './i18n-context';
import { LocaleSwitcher } from '@/components/locale-switcher';

describe('I18n System', () => {
  describe('formatError', () => {
    it('translates error messages to English when locale is en', () => {
      const wrapper = ({ children }: { children: React.ReactNode }) => (
        <I18nProvider initialLocale="en">{children}</I18nProvider>
      );
      const { result } = renderHook(() => useI18n(), { wrapper });

      const axiosError = {
        isAxiosError: true,
        response: {
          data: {
            status: 'error',
            message_code: 'E031002',
          },
        },
      };

      const formatted = result.current.formatError(axiosError);
      expect(formatted.code).toBe('E031002');
      expect(formatted.message).toBe('A tenant with this name already exists.');
    });

    it('translates error messages to Vietnamese when locale is vi', () => {
      const wrapper = ({ children }: { children: React.ReactNode }) => (
        <I18nProvider initialLocale="vi">{children}</I18nProvider>
      );
      const { result } = renderHook(() => useI18n(), { wrapper });

      const axiosError = {
        isAxiosError: true,
        response: {
          data: {
            status: 'error',
            message_code: 'E031002',
          },
        },
      };

      const formatted = result.current.formatError(axiosError);
      expect(formatted.code).toBe('E031002');
      expect(formatted.message).toBe('Tên Tenant này đã tồn tại trong hệ thống.');
    });

    it('translates field errors to Vietnamese', () => {
      const wrapper = ({ children }: { children: React.ReactNode }) => (
        <I18nProvider initialLocale="vi">{children}</I18nProvider>
      );
      const { result } = renderHook(() => useI18n(), { wrapper });

      const axiosError = {
        isAxiosError: true,
        response: {
          data: {
            status: 'error',
            message_code: 'E050000',
            error_detail: [
              {
                field: 'name',
                message_code: 'E050001',
              },
            ],
          },
        },
      };

      const formatted = result.current.formatError(axiosError);
      expect(formatted.code).toBe('E050000');
      expect(formatted.message).toBe('Dữ liệu nhập vào chưa hợp lệ.');
      expect(formatted.fieldErrors?.name).toBe('Trường này là bắt buộc.');
    });

    it('falls back to backend message when error code is not in frontend dictionary', () => {
      const wrapper = ({ children }: { children: React.ReactNode }) => (
        <I18nProvider initialLocale="vi">{children}</I18nProvider>
      );
      const { result } = renderHook(() => useI18n(), { wrapper });

      const axiosError = {
        isAxiosError: true,
        response: {
          data: {
            status: 'error',
            message_code: 'E999999',
            message: 'Thông điệp lỗi tùy chỉnh từ Backend',
          },
        },
      };

      const formatted = result.current.formatError(axiosError);
      expect(formatted.code).toBe('E999999');
      expect(formatted.message).toBe('Thông điệp lỗi tùy chỉnh từ Backend');
    });
  });

  describe('LocaleSwitcher component', () => {
    it('renders current language and allows switching to Vietnamese', () => {
      render(
        <I18nProvider initialLocale="en">
          <LocaleSwitcher />
        </I18nProvider>
      );

      // Default is EN
      expect(screen.getByText('EN')).toBeInTheDocument();

      // Click to open dropdown
      fireEvent.click(screen.getByTitle('Change application language'));

      // Check dropdown options
      expect(screen.getByText('Tiếng Việt')).toBeInTheDocument();
      expect(screen.getByText('English')).toBeInTheDocument();

      // Click Vietnamese
      fireEvent.click(screen.getByText('Tiếng Việt'));

      // Badge changes to VI
      expect(screen.getByText('VI')).toBeInTheDocument();
    });
  });

  describe('t() translation function', () => {
    it('translates common keys according to locale', () => {
      const wrapper = ({ children }: { children: React.ReactNode }) => (
        <I18nProvider initialLocale="vi">{children}</I18nProvider>
      );
      const { result } = renderHook(() => useI18n(), { wrapper });

      expect(result.current.t('common.search')).toBe('Tìm kiếm');
      expect(result.current.t('common.create')).toBe('Tạo mới');

      act(() => {
        result.current.setLocale('en');
      });

      expect(result.current.t('common.search')).toBe('Search');
      expect(result.current.t('common.create')).toBe('Create');
    });

    it('supports parameter interpolation in translations', () => {
      const wrapper = ({ children }: { children: React.ReactNode }) => (
        <I18nProvider initialLocale="en">{children}</I18nProvider>
      );
      const { result } = renderHook(() => useI18n(), { wrapper });

      expect(
        result.current.t('tenant.page.created_success', { name: 'Acme Corp' })
      ).toBe('Tenant "Acme Corp" successfully created!');

      act(() => {
        result.current.setLocale('vi');
      });

      expect(
        result.current.t('tenant.page.created_success', { name: 'Acme Corp' })
      ).toBe('Đã tạo thành công tenant "Acme Corp"!');
    });

    it('translates main action and section keys in Vietnamese', () => {
      const wrapper = ({ children }: { children: React.ReactNode }) => (
        <I18nProvider initialLocale="vi">{children}</I18nProvider>
      );
      const { result } = renderHook(() => useI18n(), { wrapper });

      expect(result.current.t('tenant.page.title')).toBe('Quản Lý Tenant');
      expect(result.current.t('tenant.page.provision_btn')).toBe('Khởi Tạo Tenant');
      expect(result.current.t('tenant.create.title')).toBe('Khởi Tạo Tenant Mới');
      expect(result.current.t('tenant.inspect.title')).toBe('Tra Cứu & Kiểm Tra Tenant');
      expect(result.current.t('tenant.inspect.inspect_btn')).toBe('Tra cứu');
      expect(result.current.t('tenant.inspect.switch_context')).toBe('Kích Hoạt Ngữ Cảnh');
    });
  });
});
