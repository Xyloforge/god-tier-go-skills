---
name: go-error-handling
description: "God-tier Go error handling — wrap with %w, sentinel vs typed errors, errors.Is/As, never swallow, nil-guard wrappers, retries/backoff, and idiomatic message style. Use when returning, wrapping, inspecting, or logging errors. If you see bare `return err`, `err.Error()` string matching, swallowed errors, or unbounded retries → apply this skill. For errgroup propagation → See `go-concurrency`; for not leaking internals in logs → See `go-observability`. Do NOT use for error string naming (→ See `go-naming`)."
origin: god_code
---

# Go Error Handling & Resilience

Part of the **God-Tier Go** set. In Go, errors are values you design, not
exceptions you throw. This skill is how production systems (Vault, Kubernetes,
Prometheus) make errors traceable, matchable, and recoverable. Every example is
cited from real code in this repo.

## When to Activate

- Returning, wrapping, or inspecting any `error`.
- Designing a package's error API (sentinels vs typed errors).
- Adding retries, backoff, timeouts, or graceful degradation.
- Reviewing code with bare `return err` or `err.Error()` string matching.

## Principles

### 1. Wrap with context using `%w`

A bare `return err` discards where the failure happened. Wrap with `fmt.Errorf`
and the `%w` verb to add context *and* preserve the chain for `errors.Is`/`As`.
Kubernetes' kubeadm wrapper does exactly this:

```go
// kubernetes/cmd/kubeadm/app/util/errors/errors.go:55 — note the nil-guard:
// wrapping nil returns nil, so callers can wrap unconditionally.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return &errorWithStack{
		msg:   fmt.Errorf("%s: %w", message, err), // %w keeps the chain
		cause: err,
		stack: callStack(),
	}
}
```

Why god-tier: `%w` (not `%v`) means a caller three layers up can still
`errors.Is(err, ErrTarget)`. The nil-guard lets the helper be called without an
`if err != nil` at every site.

### 2. Sentinel errors for expected, matchable conditions

When callers need to branch on a specific condition, expose a package-level
sentinel and let them match it with `errors.Is`. Vault does this throughout:

```go
// vault-main/sdk/database/helper/connutil/connutil.go:12
var ErrNotInitialized = errors.New("connection has not been initialized")
```

```go
// prometheus-main/tsdb/head_append_v2.go:186 — callers match sentinels, never
// string-compare. errors.Is walks the whole %w chain.
case errors.Is(appErr, storage.ErrOutOfOrderSample):
	...
case errors.Is(appErr, storage.ErrTooOldSample):
	...
```

Why god-tier: `errors.Is` matches through any depth of wrapping; string matching
(`err.Error() == "..."`) breaks the moment someone adds context. Export the
sentinel; never make callers guess the message.

### 3. Typed errors + `errors.As` when callers need the details

If the caller needs structured data from the failure (a field, a code), define a
type that implements `error` and (when wrapping) `Unwrap()`, then extract it with
`errors.As`:

```go
// kubernetes/cmd/kubeadm/app/util/errors/errors.go:41 — a typed error that
// participates in the chain by implementing Error() and Unwrap().
type errorWithStack struct {
	msg   error
	cause error
	stack string
}
func (e errorWithStack) Error() string { return e.msg.Error() }
func (e errorWithStack) Unwrap() error { return e.cause }
```

```go
// Caller extracts the concrete type out of an arbitrarily-wrapped error:
var se *errorWithStack
if errors.As(err, &se) {
	log.Print(se.stack)
}
```

Why god-tier: `Unwrap()` is what makes the type discoverable through layers of
`%w`. Implement it whenever your error wraps another.

### 4. Never swallow; handle exactly once at the right layer

An error is either handled or returned — never silently dropped, and never both
logged *and* returned (that double-reports). Push the decision to the layer that
can actually act on it.

```go
// Bad — swallowed: the caller cannot know anything failed.
v, _ := strconv.Atoi(s)

// Bad — handled twice: logged here AND returned, so it's logged again above.
if err != nil {
	log.Printf("parse: %v", err)
	return err
}

// Good — add context and return; let the top layer log/decide once.
if err != nil {
	return fmt.Errorf("parse port %q: %w", s, err)
}
```

Why god-tier: single ownership of each error keeps logs clean and behavior
predictable. The lowest layer adds context; one upper layer decides.

### 5. Resilience: retry transient failures with bounded backoff

Not every error is permanent. For transient ones (network, contention), retry
with a *bounded* count and backoff — and respect context cancellation so a
caller can give up.

```go
// Shape of a god-tier retry: bounded, backed off, context-aware.
func withRetry(ctx context.Context, attempts int, fn func() error) error {
	var err error
	for i := 0; i < attempts; i++ {
		if err = fn(); err == nil {
			return nil
		}
		if !isTransient(err) {
			return err // permanent: fail fast, don't waste attempts
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff(i)): // e.g. exponential
		}
	}
	return fmt.Errorf("after %d attempts: %w", attempts, err)
}
```

Why god-tier: unbounded or context-blind retries turn a blip into an outage.
Distinguish transient from permanent; always honor cancellation. See
[[go-concurrency]] for context propagation.

### 6. Message style: lowercase, no trailing punctuation, no redundancy

Error strings are fragments that get wrapped into larger sentences, so the Go
convention is lowercase with no ending punctuation. The Vault sentinel in
Principle 2 (`"connection has not been initialized"`) follows it exactly.

Because wrapping *compounds* messages, avoid stacking a `failed to` prefix at
every layer — `failed to start: failed to dial: failed to resolve` reads badly.
State the action, let `%w` chain the cause.

```go
// Good — fragment composes cleanly when wrapped.
errors.New("connection has not been initialized")

// Avoid — capital + punctuation breaks when embedded mid-sentence.
errors.New("Connection has not been initialized.")
```

Why god-tier: consistent fragments wrap into readable chains. (Note: large
codebases vary — some Vault paths do use `failed to …` — but the lowercase /
no-punctuation rule is near-universal and is the Go Code Review Comments
standard.)

## Anti-Patterns

| Smell | Fix |
|-------|-----|
| `return err` with no context | Wrap: `fmt.Errorf("doing X: %w", err)` (Principle 1) |
| `if err.Error() == "not found"` | Sentinel + `errors.Is` (Principle 2) |
| `fmt.Errorf("...: %v", err)` when callers must match | Use `%w` (Principle 1) |
| Custom error type without `Unwrap()` | Add `Unwrap()` so `errors.As` works (Principle 3) |
| `v, _ := f()` | Don't discard; handle or wrap+return (Principle 4) |
| Log *and* return the same error | Pick one — return down low, log once up top (Principle 4) |
| `for { if fn() == nil { break } }` retry | Bound it, back off, honor ctx (Principle 5) |
| `errors.New("Failed To Do X.")` | lowercase, no punctuation (Principle 6) |

## Checklist

- [ ] Every returned error adds context with `%w` or is a deliberate sentinel.
- [ ] No `err.Error()` string comparison anywhere — use `errors.Is`/`As`.
- [ ] Custom error types that wrap implement `Unwrap()`.
- [ ] No discarded errors (`_`), and no error both logged and returned.
- [ ] Retries are bounded, backed off, and respect `context` cancellation.
- [ ] Error strings are lowercase, no trailing punctuation, non-redundant.

## Related

- [[go-clean-code]] — clarity in the error path.
- [[go-function-design]] — designing the `error` you return last.
- [[go-concurrency]] — `errgroup` propagation and context cancellation.
- [[go-testing]] — table-driven tests over the error branches.
