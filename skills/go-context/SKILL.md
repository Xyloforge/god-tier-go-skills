---
name: go-context
description: "Idiomatic use of context.Context in Go — cancellation, deadlines/timeouts, and request-scoped values. Use when adding context to a function, propagating cancellation, setting timeouts, or storing request-scoped values. If you see a Context stored in a struct field, a context.Value with a string key, or a missing ctx.Done() in a blocking loop → apply this skill. For goroutine fan-out and errgroup → See `go-concurrency`. For ctx.Err()/retries → See `go-error-handling`. Do NOT use for general goroutine design."
origin: god_code
---

# Go Context

Part of the **God-Tier Go** set. `context.Context` is the cancellation and
request-scope spine of Go services. This is a thin router; the cited depth lives
in `references/`.

## When to Activate

- Adding `ctx context.Context` to a function or threading it through a call tree.
- Setting deadlines/timeouts, or propagating cancellation into blocking work.
- Storing/reading request-scoped values (trace IDs, auth) on a context.
- Reviewing code with `ctx` in a struct field or a `string`-keyed value.

## Decision Guide

| Situation | Do this | Depth |
|-----------|---------|-------|
| Function does I/O or can block | Take `ctx` as the **first** parameter | `references/propagation.md` |
| Need a time/operation bound | `context.WithTimeout`/`WithCancel`, **always `defer cancel()`** | `references/propagation.md` |
| Loop/select that can hang | Add a `case <-ctx.Done(): return ctx.Err()` | `references/propagation.md` |
| Attach request-scoped data | Custom **unexported pointer key type**, never a string | `references/values.md` |
| Tempted to put `ctx` in a struct | Don't — pass it per-call | `references/propagation.md` |

## Core Rules

1. **`ctx` is the first parameter**, named `ctx`, on every cancelable/fallible call.
2. **Never store a `Context` in a struct** — its lifetime is the call, not the object.
3. **Always `defer cancel()`** after `WithCancel`/`WithTimeout`/`WithDeadline` — not cancelling leaks the timer and goroutine.
4. **Context values use a private key type**, not a `string` (collision- and alloc-safe). See chi's pattern in `references/values.md`.
5. **Context values are for request scope only** — trace IDs, auth, deadlines; never for optional function parameters.
6. **Propagate, don't swallow** — pass the received `ctx` down; derive children, never `context.Background()` mid-tree.

## Anti-Patterns

| Smell | Fix |
|-------|-----|
| `ctx` stored as a struct field | Pass per-call as first arg (Rule 2) |
| `context.WithTimeout(...)` with no `defer cancel()` | Always `defer cancel()` (Rule 3) |
| `context.WithValue(ctx, "userID", id)` (string key) | Private key type (Rule 4) |
| `context.Background()` deep inside a request | Thread the real `ctx` down (Rule 6) |
| Blocking `select`/channel op with no `ctx.Done()` | Add the cancellation case (Rule 3) |
| Passing config through `ctx.Value` | Use explicit parameters (Rule 5) |

## Checklist

- [ ] `ctx context.Context` is the first parameter of every blocking/fallible function.
- [ ] No `Context` is stored in a struct field.
- [ ] Every `WithCancel`/`WithTimeout`/`WithDeadline` is paired with `defer cancel()`.
- [ ] Context value keys are an unexported pointer/struct type, never a string.
- [ ] `ctx.Value` carries only request-scoped data, not optional args.
- [ ] Blocking loops/selects honor `<-ctx.Done()`.

## Deep Dives

- `references/propagation.md` — cancellation, timeouts, `defer cancel()`, `ctx.Done()` (cited: Moby, stdlib).
- `references/values.md` — the typed-key pattern and request scope (cited: chi `context.go`).

## Related

- → See `go-concurrency` — goroutine lifecycle, `errgroup`, fan-out.
- → See `go-error-handling` — returning `ctx.Err()`, bounded retries.
- → See `go-function-design` — why `ctx` is a parameter, not stored state.
