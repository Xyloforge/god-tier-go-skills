---
name: go-observability
description: "Instrumenting Go services — Prometheus metrics (counters/histograms/gauges + naming), structured/contextual logging with slog, and tracing hooks. Use when adding metrics, structured logs, or trace spans, or reviewing instrumentation. If you see fmt.Println debugging, counters without a _total suffix, latency as a gauge instead of a histogram, or log lines built by string concatenation → apply this skill. For error logging specifics → See `go-error-handling`. Do NOT use for performance profiling (→ See `go-performance`)."
origin: god_code
---

# Go Observability

Part of the **God-Tier Go** set. The three pillars — metrics, logs, traces — make
a running service legible. This distills how Prometheus and Moby instrument
themselves. Thin router; cited depth in `references/`.

## When to Activate

- Adding metrics (request rates, latencies, queue depths, error counts).
- Adding structured logs or wiring a contextual logger through requests.
- Adding tracing spans / propagating trace context.
- Reviewing instrumentation for correct metric types and naming.

## Decision Guide

| You want to measure… | Use | Depth |
|----------------------|-----|-------|
| A monotonically increasing count | `Counter` (name ends `_total`) | `references/metrics.md` |
| A value that goes up and down | `Gauge` | `references/metrics.md` |
| A distribution (latency, size) | `Histogram` (unit suffix, e.g. `_seconds`) | `references/metrics.md` |
| What happened, with fields | structured `slog` (key/values) | `references/logging-and-tracing.md` |
| Per-request context in logs | logger carried on `context.Context` | `references/logging-and-tracing.md` |
| Cross-service latency | tracing spans (OpenTelemetry) | `references/logging-and-tracing.md` |

## Core Rules

1. **Pick the right metric type:** `Counter` for cumulative counts, `Gauge` for instantaneous values, `Histogram` for distributions. Latency is a histogram, never a gauge.
2. **Name metrics by convention:** `namespace_subsystem_name`, `_total` suffix for counters, base-unit suffix for histograms (`_seconds`, `_bytes`). Always set `Help`.
3. **Labels are bounded dimensions, never unbounded IDs.** Label by `method`/`status`/`alertmanager`; never by user ID, request ID, or raw URL — that explodes cardinality.
4. **Log structured key/values, not formatted strings.** `slog.Info("msg", "key", val)` so logs are queryable in aggregation tools.
5. **Carry the logger on context** so every line in a request automatically has its request-scoped fields.
6. **Levels mean things:** `Error` = needs attention, `Warn` = suspicious, `Info` = lifecycle, `Debug` = diagnostics. Don't log-and-return errors twice (see `go-error-handling`).
7. **Never log secrets/PII.** Redact tokens, keys, and personal data (see `go-security`).

## Anti-Patterns

| Smell | Fix |
|-------|-----|
| `fmt.Println`/`log.Printf` debugging in prod | Structured `slog` with levels (Rule 4) |
| Counter named `requests` (no `_total`) | `requests_total` (Rule 2) |
| Latency exposed as a `Gauge` | `Histogram` with `_seconds` (Rule 1,2) |
| Label `user_id` / `request_id` on a metric | Bounded labels only; put IDs in logs/traces (Rule 3) |
| `log.Printf("user %s did %s", u, a)` | `slog.Info("action", "user", u, "action", a)` (Rule 4) |
| Metric with no `Help` text | Always set `Help` (Rule 2) |
| Logging an API key or password | Redact (Rule 7) |

## Checklist

- [ ] Each metric uses the correct type (counter/gauge/histogram).
- [ ] Counters end in `_total`; histograms carry a base-unit suffix; all have `Help`.
- [ ] Metric labels are low-cardinality; no IDs/emails/URLs as labels.
- [ ] Logs are structured key/values via `slog`, not formatted strings.
- [ ] A request-scoped logger is carried on `context.Context`.
- [ ] Log levels are used meaningfully; errors aren't double-logged.
- [ ] No secrets or PII appear in logs or metric labels.

## Deep Dives

- `references/metrics.md` — counters/gauges/histograms, naming, labels (cited: Prometheus `notifier/metric.go`).
- `references/logging-and-tracing.md` — `slog`, contextual loggers, tracing (cited: Moby).

## Assets

- `assets/recording-and-alert-rules.yml` — Prometheus recording + alert rules for RED-method metrics.

## Related

- → See `go-error-handling` — logging errors once, at the right layer.
- → See `go-context` — carrying a request logger and deadlines on context.
- → See `go-performance` — instrument first, then optimize the measured hot path.
- → See `go-security` — keeping secrets/PII out of telemetry.
