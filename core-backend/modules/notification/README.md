# Notification Module (STUB)

## Purpose
Dispatches transactional emails, webhooks, and push notifications triggered by domain events across the system.

## Planned Ports & Interfaces
- **Ports (Domain)**:
  - `Sender`: Interface for outbound message transport (`Send(ctx, msg)`).
  - `TemplateRenderer`: Compiling dynamic email templates.
- **Adapters**:
  - `LogSender`: Development placeholder that logs outbound notifications via `slog`.
  - `SmtpSender`: Outbound SMTP client (e.g. AWS SES, Resend, Postmark).
  - `WebhookSender`: HTTP webhook dispatcher with HMAC signature headers.
- **Events Listened To**:
  - `iam.user.registered`: Sends welcome verification email.
  - `tenant.created`: Sends workspace confirmation.

## TODO Checklist
- [ ] Connect `Sender` port to email delivery service (e.g., Resend, AWS SES)
- [ ] Implement HTML template rendering with `html/template`
- [ ] Wire background worker or eventbus listener for asynchronous dispatch
- [ ] Add retry backoff policy for failed webhook deliveries
