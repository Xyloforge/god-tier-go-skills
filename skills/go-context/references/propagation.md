# Context Propagation — cancellation, timeouts, deadlines

A `context.Context` flows *down* a call tree. Each layer either passes it through
unchanged or derives a more-constrained child (shorter deadline, its own cancel)
and passes that. Cancellation propagates from parent to all children.

## `ctx` is the first parameter — always

```go
// moby/daemon/stop.go:25 — context.Context is the first arg, by convention.
func (daemon *Daemon) ContainerStop(ctx context.Context, name string, options backend.ContainerStopOptions) error
```

A function that does I/O, blocks, or calls something that does must accept a
`ctx`. This is non-negotiable in service code: it's how a request timeout or a
client disconnect unwinds the entire stack.

## Derive a child with a bound, and always `defer cancel()`

When a step needs its own deadline, derive a child context from the incoming one.
The returned `cancel` **must** be called (defer it immediately) — otherwise the
parent retains the child until *it* is cancelled, leaking the timer goroutine.

```go
// moby/daemon/stop.go:90 — derive a bounded child from the parent ctx, defer
// cancel immediately, then pass the *child* into the blocking wait.
var subCtx context.Context
var cancel context.CancelFunc
if stopTimeout >= 0 {
	subCtx, cancel = context.WithTimeout(ctx, wait)
} else {
	subCtx, cancel = context.WithCancel(ctx)
}
defer cancel()

if status := <-ctr.State.Wait(subCtx, containertypes.WaitConditionNotRunning); status.Err() == nil {
	return nil
}
```

Why god-tier:
- The child inherits the parent's cancellation **and** adds its own bound — if the
  caller gives up, this step gives up too.
- `defer cancel()` runs on every return path, so there's no leaked timer even on
  the happy path. `go vet` flags a missing cancel.

## Honor `ctx.Done()` in anything that can block

A blocking loop or `select` must have a cancellation arm, or it ignores timeouts
and disconnects entirely.

```go
func (w *Worker) run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()      // propagate WHY we stopped (Canceled/DeadlineExceeded)
		case job := <-w.jobs:
			w.handle(ctx, job)    // pass ctx onward
		}
	}
}
```

Return `ctx.Err()` so callers can distinguish `context.Canceled` from
`context.DeadlineExceeded` with `errors.Is` (see `go-error-handling`).

## Rules of propagation

- **Pass the received `ctx` down.** Never replace it with `context.Background()`
  partway through a request — that severs cancellation for everything below.
- **Background/TODO only at the top.** `context.Background()` belongs in `main`,
  tests, and top-level init; `context.TODO()` marks a spot you haven't wired yet.
- **A child can only tighten, never loosen.** A 1s child of a 10s parent dies at
  1s; you cannot extend a parent's deadline from a child.
- **One cancel source unwinds the tree.** This is the whole point: a single
  `cancel()` or deadline at the top cascades to every derived context below.
