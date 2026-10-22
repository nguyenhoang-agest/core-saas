package i18n

import (
	"context"
	"strings"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

const (
	LocaleEN = "en"
	LocaleVI = "vi"
)

var SupportedLocales = []string{LocaleEN, LocaleVI}

type contextKey string

const localeContextKey contextKey = "locale"

type Bundle struct {
	mu     sync.RWMutex
	bundle *i18n.Bundle
}

var (
	defaultBundle *Bundle
	once          sync.Once
)

func GetDefaultBundle() *Bundle {
	once.Do(func() {
		b := i18n.NewBundle(language.English)
		defaultBundle = &Bundle{
			bundle: b,
		}
	})
	return defaultBundle
}

// RegisterMessages adds message translations for a specific locale tag (e.g. "en", "vi").
func RegisterMessages(locale string, messages map[string]string) {
	GetDefaultBundle().RegisterMessages(locale, messages)
}

func (b *Bundle) RegisterMessages(locale string, messages map[string]string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	tag, err := language.Parse(locale)
	if err != nil {
		tag = language.English
	}

	for id, text := range messages {
		_ = b.bundle.AddMessages(tag, &i18n.Message{
			ID:    id,
			Other: text,
		})
	}
}

// NormalizeLocale turns an input tag (e.g. "vi-VN", "en_US", "VI") into "en" or "vi".
func NormalizeLocale(raw string) string {
	raw = strings.TrimSpace(strings.ToLower(raw))
	if raw == "" {
		return LocaleEN
	}

	// Handle standard Accept-Language format like "vi-VN,vi;q=0.9,en;q=0.8"
	if strings.Contains(raw, ",") {
		parts := strings.Split(raw, ",")
		for _, part := range parts {
			sub := strings.TrimSpace(strings.Split(part, ";")[0])
			normalized := NormalizeLocale(sub)
			if normalized == LocaleVI {
				return LocaleVI
			}
		}
		return LocaleEN
	}

	if strings.HasPrefix(raw, "vi") {
		return LocaleVI
	}

	return LocaleEN
}

// WithLocale returns a context carrying the normalized locale.
func WithLocale(ctx context.Context, locale string) context.Context {
	return context.WithValue(ctx, localeContextKey, NormalizeLocale(locale))
}

// GetLocale extracts the normalized locale from context, defaulting to "en".
func GetLocale(ctx context.Context) string {
	if ctx == nil {
		return LocaleEN
	}
	if val, ok := ctx.Value(localeContextKey).(string); ok && val != "" {
		return val
	}
	return LocaleEN
}

// T translates a message ID according to the locale in ctx. If template data is provided, it is interpolated.
func T(ctx context.Context, messageID string, templateData ...map[string]any) string {
	return GetDefaultBundle().T(ctx, messageID, templateData...)
}

func (b *Bundle) T(ctx context.Context, messageID string, templateData ...map[string]any) string {
	if messageID == "" {
		return ""
	}

	locale := GetLocale(ctx)

	b.mu.RLock()
	defer b.mu.RUnlock()

	var data map[string]any
	if len(templateData) > 0 {
		data = templateData[0]
	}

	localizer := i18n.NewLocalizer(b.bundle, locale)
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})
	if err == nil && msg != "" {
		return msg
	}

	// Fallback to English if primary locale didn't match
	if locale != LocaleEN {
		fallbackLocalizer := i18n.NewLocalizer(b.bundle, LocaleEN)
		fallbackMsg, fallbackErr := fallbackLocalizer.Localize(&i18n.LocalizeConfig{
			MessageID:    messageID,
			TemplateData: data,
		})
		if fallbackErr == nil && fallbackMsg != "" {
			return fallbackMsg
		}
	}

	return ""
}
