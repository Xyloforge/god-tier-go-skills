# Struct Design — composition, zero values, layout

## Compose with embedding

Go has no inheritance. You build bigger types by **embedding** smaller ones,
which promotes their fields and methods to the outer type. Moby's
`RestartManager` embeds two `sync` primitives and gets their methods for free:

```go
// moby/daemon/internal/restartmanager/restartmanager.go:22 — embeds sync.Mutex
// and sync.Once, so rm.Lock()/rm.Unlock() and rm.Do(...) work directly.
type RestartManager struct {
	sync.Mutex
	sync.Once
	policy       container.RestartPolicy
	restartCount int
	timeout      time.Duration
	active       bool
	cancel       chan struct{}
	canceled     bool
}
```

Embedding expresses "this type *is* lockable / *is* once-guarded." Use a **named
field** instead when the relationship is "has-a" and you don't want method
promotion:

```go
type Server struct {
	logger *slog.Logger // has-a; we don't want Server to expose logger's methods
}
```

Caution: embedding a `sync.Mutex` (as above) means the outer struct **must not be
copied** once used — copying duplicates the lock. Pass `*RestartManager`, never
`RestartManager`. See `go-concurrency` on copy semantics.

## Make the zero value useful

The best constructor is none. Embedded `sync.Mutex` and `sync.Once` above are
usable at their zero value, so `RestartManager` only needs a constructor for the
fields that genuinely require initialization (the channel):

```go
// moby/daemon/internal/restartmanager/restartmanager.go:34 — New only sets what
// can't be a useful zero value (the cancel channel); the sync fields need nothing.
func New(policy container.RestartPolicy, restartCount int) *RestartManager {
	return &RestartManager{policy: policy, restartCount: restartCount, cancel: make(chan struct{})}
}
```

A small value struct can skip the constructor entirely:

```go
// chi-master/context.go:160 — tiny value type, zero value meaningful, no ctor.
type contextKey struct {
	name string
}
```

Rule: provide a constructor only when the zero value would be invalid (needs an
open channel, a non-nil map, a dialed connection). Otherwise let `var x T` work.

## Field layout

- **Group fields by meaning**, not by type. Keep related fields together so the
  struct reads as a description of the thing.
- **Embedded types go first**, before named fields (as in `RestartManager`).
- **Don't hand-tune field order for memory padding** unless a benchmark proves it
  matters — clarity wins by default (see `go-performance` for when it doesn't).
- Tag fields only where needed (`json:"..."`, `db:"..."`); keep tags accurate —
  they're part of the type's contract.
