package i18n_test

import (
	"context"
	"testing"

	"github.com/example/core-backend/pkg/i18n"
)

func TestNormalizeLocale(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", "en"},
		{"en", "en"},
		{"en-US", "en"},
		{"en_GB", "en"},
		{"vi", "vi"},
		{"vi-VN", "vi"},
		{"VI", "vi"},
		{"vi-VN,vi;q=0.9,en;q=0.8", "vi"},
		{"fr-FR,fr;q=0.9,en;q=0.8", "en"},
	}

	for _, tt := range tests {
		actual := i18n.NormalizeLocale(tt.input)
		if actual != tt.expected {
			t.Errorf("NormalizeLocale(%q) = %q; want %q", tt.input, actual, tt.expected)
		}
	}
}

func TestBundle_Translation(t *testing.T) {
	i18n.RegisterMessages(i18n.LocaleEN, map[string]string{
		"TEST_MSG_1": "Hello world",
		"TEST_MSG_2": "English only message",
	})
	i18n.RegisterMessages(i18n.LocaleVI, map[string]string{
		"TEST_MSG_1": "Xin chào thế giới",
	})

	ctxEn := i18n.WithLocale(context.Background(), "en")
	ctxVi := i18n.WithLocale(context.Background(), "vi")

	// Test Vietnamese
	if msg := i18n.T(ctxVi, "TEST_MSG_1"); msg != "Xin chào thế giới" {
		t.Errorf("expected 'Xin chào thế giới', got %q", msg)
	}

	// Test English
	if msg := i18n.T(ctxEn, "TEST_MSG_1"); msg != "Hello world" {
		t.Errorf("expected 'Hello world', got %q", msg)
	}

	// Test fallback to English when Vietnamese translation is missing
	if msg := i18n.T(ctxVi, "TEST_MSG_2"); msg != "English only message" {
		t.Errorf("expected fallback to 'English only message', got %q", msg)
	}

	// Test unknown message returns empty string
	if msg := i18n.T(ctxEn, "UNKNOWN_CODE"); msg != "" {
		t.Errorf("expected empty string for unknown code, got %q", msg)
	}
}
