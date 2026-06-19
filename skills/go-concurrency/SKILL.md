---
name: go-concurrency
description: "God-tier Go concurrency — goroutine lifecycle and leak prevention, context cancellation, channels vs mutexes, errgroup fan-out with limits, sync primitives, and race-free design. Use when starting a goroutine, channel, or sync primitive, or fanning out work. If you see fire-and-forget goroutines, unbounded fan-out, a channel guarding one field, or a copied mutex → apply this skill. For context propagation → See `go-context`; for pool reuse → See `go-performance`. Do NOT use for general function design (→ See `go-function-design`)."
origin: god_code
---

# Go Concurrency

Part of the **God-Tier Go** set. Concurrency is cheap to start and expensive to
get wrong: leaks, races, and deadlocks. This skill is how production systems
(Kubernetes, Moby, the stdlib) keep goroutines bounded, cancelable, and safe.
Every example is cited from real code in this repo.

## When to Activate

- Starting any goroutine, channel, or `sync` primitive.
- Fanning out work concurrently (worker pools, parallel visits).
- Threading cancellation/deadlines through a call tree.
- Reviewing code for leaks, races, or unbounded concurrency.

## Principles

### 1. Every goroutine must have a known exit — prevent leaks

A goroutine that can block forever is a leak. The fix is twofold: give it a
cancellation signal (`ctx.Done()`), and make any channel it sends on *buffered*
so it can't block after the receiver has gone. Moby's attach loop does both:

```go
// moby/daemon/internal/stream/attach.go:133 — buffered (cap 1) so the inner
// goroutine can always send and exit even if no one is reading; the outer
// select races group completion against context cancellation.
errs := make(chan error, 1)
go func() {
	groupErr := make(chan error, 1) // buffered: sender never blocks
	go func() { groupErr <- group.Wait() }()
	select {
	case <-ctx.Done():
		// close all pipes so the blocked work unwinds
		...
	case err := <-groupErr:
		...
	}
}()
```

Why god-tier: the `cap 1` buffer is the difference between a clean exit and a
permanently parked goroutine. Always ask: "how does this goroutine end?"

### 2. `context.Context` is the cancellation spine — pass it, don't store it

Thread `ctx` as the first parameter through the call tree; select on
`ctx.Done()` in any loop or blocking wait. Never stash a `Context` in a struct
field — its lifetime belongs to the call, not the object.

```go
// Idiomatic shutdown of a worker loop:
func (w *Worker) run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err() // single, predictable exit
		case job := <-w.jobs:
			w.handle(ctx, job)
		}
	}
}
```

Why god-tier: one cancellation source propagates everywhere, so a single
`cancel()` (or deadline) unwinds the whole tree. See [[go-error-handling]] for
returning `ctx.Err()` and bounded retries.

### 3. Bound your fan-out with `errgroup` + `SetLimit`

Spawning one goroutine per item is a denial-of-service against your own process.
Use `errgroup` to collect the first error *and* cancel siblings, and `SetLimit`
to cap concurrency. Kubernetes' concurrent visitor is the template:

```go
// kubernetes/staging/src/k8s.io/cli-runtime/pkg/resource/visitor.go:210
func (l ConcurrentVisitorList) Visit(fn VisitorFunc) error {
	g := errgroup.Group{}
	concurrency := 1            // safe default: sequential
	if l.concurrency > concurrency {
		concurrency = l.concurrency
	}
	g.SetLimit(concurrency)    // bound the parallelism
	for i := range l.visitors {
		g.Go(func() error {
			return l.visitors[i].Visit(fn)
		})
	}
	return g.Wait()            // first error wins; waits for all
}
```

Why god-tier: bounded, error-aware, and defaults to sequential — concurrency is
opt-in, not a footgun. `g.Wait()` is the single join point.

### 4. Channels to transfer ownership; mutexes to protect state

"Share memory by communicating" — pass data through a channel when you're
handing off ownership. But a plain mutex is clearer and faster when you're just
guarding a few fields. Pick by intent, don't cargo-cult channels everywhere.

```go
// Good — mutex guards shared state. Simple, fast, obvious.
type Counter struct {
	mu sync.Mutex
	n  int
}
func (c *Counter) Inc() { c.mu.Lock(); c.n++; c.mu.Unlock() }
```

Why god-tier: channels for *flow*, mutexes for *state*. Using a channel to guard
a single int is slower and harder to read than a `sync.Mutex`.

### 5. `sync.Once` / `sync.WaitGroup` for the jobs they're built for

Reach for the right primitive instead of hand-rolling. `sync.Once` guarantees
exactly-once init even under concurrent callers — Moby embeds it directly:

```go
// moby/daemon/internal/restartmanager/restartmanager.go:24 — embedded Once
// guarantees the restart logic initializes exactly once across goroutines.
type restartManager struct {
	sync.Mutex
	sync.Once
	...
}
```

Why god-tier: `Once.Do(f)` is race-free and clearer than a `mu`-guarded
`initialized bool`. Use `WaitGroup` to await N goroutines; use `Once` for
exactly-once.

### 6. Respect copy semantics — never copy a sync type after use

`sync.Mutex`, `sync.Pool`, `sync.WaitGroup`, and `sync.Once` must not be copied
once used. The stdlib states the contract explicitly:

```go
// gostd/sync/pool.go:44
// A Pool must not be copied after first use.
//
// gostd/sync/pool.go:22
// A Pool is safe for use by multiple goroutines simultaneously.
type Pool struct { ... }
```

Why god-tier: copying a used `sync` value silently breaks mutual exclusion —
`go vet` catches it, and passing these types by pointer (or never copying the
enclosing struct) avoids it. This is why mixed receivers (see
[[go-function-design]]) are dangerous on types embedding a mutex.

## Anti-Patterns

| Smell | Fix |
|-------|-----|
| `go f()` with no way to stop it | Bind to `ctx`; ensure a known exit (Principle 1) |
| Unbuffered send a dead receiver may abandon | Buffer the channel (cap 1) (Principle 1) |
| `context.Context` stored in a struct field | Pass it as first param (Principle 2) |
| One goroutine per item, unbounded | `errgroup` + `SetLimit` (Principle 3) |
| Channel used to guard a single field | Use a `sync.Mutex` (Principle 4) |
| `if !inited { init(); inited = true }` race | `sync.Once.Do` (Principle 5) |
| Copying a struct that embeds a `Mutex` | Use a pointer; never copy after use (Principle 6) |
| Shipping concurrent code untested for races | `go test -race ./...` (Checklist) |

## Checklist

- [ ] Every goroutine has a defined stop condition (ctx, closed channel, or join).
- [ ] Channels a goroutine sends on are buffered when the receiver may vanish.
- [ ] `context.Context` is passed as the first arg, never stored in a struct.
- [ ] Fan-out concurrency is bounded (`errgroup.SetLimit` or a worker pool).
- [ ] Channels move ownership; mutexes guard state — chosen deliberately.
- [ ] `sync.Once`/`WaitGroup` used instead of hand-rolled flags/counters.
- [ ] No `sync` type is copied after first use.
- [ ] Tested with `go test -race ./...`.

## Related

- [[go-error-handling]] — `errgroup` error propagation, `ctx.Err()`, retries.
- [[go-performance]] — `sync.Pool` reuse and lock contention.
- [[go-function-design]] — why context is a parameter and receiver consistency.
- [[go-testing]] — the `-race` detector and concurrent test design.
