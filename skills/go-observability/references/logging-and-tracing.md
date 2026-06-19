# Structured Logging & Tracing

## Structured logging — key/values, not formatted strings

A log line built with `fmt`/`Printf` is a string an aggregator can't query.
Structured logging emits machine-parseable key/value pairs. Use the stdlib
`log/slog` (Go 1.21+):

```go
// Bad — a string; you can't filter by user or action in Loki/ELK.
log.Printf("user %s performed %s on %s", userID, action, resource)

// Good — queryable fields.
slog.Info("user action",
	"user", userID,
	"action", action,
	"resource", resource,
)
```

Why god-tier: `slog.Info("...", k, v, ...)` produces structured output (JSON via
`slog.NewJSONHandler`) that aggregation tools index by field. You can then query
`action="delete" AND resource="billing"` instead of grepping prose.

## Contextual logging — carry the logger on context

Derive a logger enriched with request-scoped fields once, attach it to the
`context.Context`, and every downstream call logs with those fields automatically.
Moby does exactly this — `log.G(ctx)` fetches the context's logger:

```go
// moby/daemon/internal/stream/attach.go:135 — log via the context-carried
// logger, so the line inherits whatever fields were attached upstream.
defer log.G(ctx).Debug("attach done")
```

The pattern, with `slog`:

```go
// At the request boundary, enrich and stash the logger.
logger := slog.With("request_id", reqID, "route", route)
ctx = ctxlog.WithLogger(ctx, logger)

// Anywhere downstream, retrieve and log — fields come for free.
ctxlog.From(ctx).Info("validated payload", "bytes", n)
```

This is why request IDs belong in logs (carried on context), not in metric
labels (`references/metrics.md`) — logs are high-cardinality by design.

## Log levels

- **Error** — something needs human attention; an operation failed.
- **Warn** — recoverable but suspicious; degraded mode.
- **Info** — normal lifecycle events (startup, request completed).
- **Debug** — diagnostics, off in production by default.

Don't log *and* return the same error (that double-reports) — return it with
context low, log it once at the top. See `go-error-handling`.

## Tracing

For latency that crosses services, metrics tell you *that* it's slow and tracing
tells you *where*. Use OpenTelemetry:

- Propagate trace context across boundaries (it rides on `context.Context`, like
  the logger — see `go-context`).
- Create a span per significant unit of work; record key attributes (low
  cardinality, like metric labels) and errors on the span.
- Correlate the three pillars by stamping the `trace_id` into structured logs so
  a slow trace links straight to its logs.

## Never log secrets or PII

Tokens, API keys, passwords, and personal data must be redacted before logging
(see `go-security`). Structured fields make redaction easier — filter known
sensitive keys at the handler — but the discipline is yours to enforce.
