---
name: go-function-design
description: "Design the shape of a Go function — small interface params, context-first/error-last signatures, when to extract shared helpers, functional options for constructors, and receiver choice. Use when defining any signature, constructor, or helper. If you see boolean-trap params, 8-arg constructors, or concrete types where an interface fits → apply this skill. For type definitions → See `go-structs-interfaces`; for the prose inside → See `go-clean-code`. Do NOT use for naming decisions (→ See `go-naming`)."
origin: god_code
---

# Go Function & Helper Design

Part of the **God-Tier Go** set. Where [[go-clean-code]] governs the *prose
inside* a function, this skill governs its *shape*: what it accepts, what it
returns, whether it should exist, and how it stays reusable. Every example is
cited from real code in this repo.

## When to Activate

- Designing a new function, method, or constructor signature.
- Deciding whether to extract a helper or inline logic.
- Building something configurable (a client, server, writer) with many options.
- Choosing pointer vs value receivers.

## Principles

### 1. Accept interfaces, return structs — and keep interfaces tiny

The most reusable functions ask for the *smallest behavior* they need, not a
concrete type. The Go stdlib is built on one-method interfaces, so any function
taking an `io.Reader` works with files, sockets, buffers, and `gzip` streams it
never heard of:

```go
// gostd/io/io.go:86 — one method. This is why io.Reader composes everywhere.
type Reader interface {
	Read(p []byte) (n int, err error)
}
```

Why god-tier: a parameter typed as the minimal interface maximizes the set of
callers and makes the function trivially testable with a fake. Return a concrete
struct, though — callers want the full type; let *them* narrow it.

### 2. Signature ordering: `context.Context` first, `error` last

Idiomatic Go puts `ctx context.Context` as the first parameter and `error` as
the last return value. It's a convention the whole ecosystem relies on, so
honoring it makes your API instantly legible:

```go
// Shape every fallible, cancelable function this way:
func FetchUser(ctx context.Context, id string) (*User, error)
//             ^^^ context first                       error last ^^^
```

Why god-tier: predictable shape removes a decision from every reader and every
caller. Never store a `Context` in a struct — pass it; see [[go-concurrency]].

### 3. No boolean traps — name the variant

`Process(data, true, false)` is unreadable at the call site. Replace boolean
parameters with distinct functions, an enum, or an option:

```go
// Bad — what do true, false mean here?
srv.Listen(addr, true, false)

// Good — intent is explicit at the call site.
srv.ListenTLS(addr)
srv.Listen(addr, WithKeepAlive(false))
```

Why god-tier: the call site is read far more often than the signature. Make it
self-documenting. For the *readability* angle of this, see [[go-clean-code]].

### 4. Functional options for extensible constructors

When a constructor has optional, growable configuration, take
`opts ...Option` instead of a widening parameter list or a config struct full of
zero-value ambiguity. Prometheus's chunk writer is a clean example:

```go
// prometheus-main/tsdb/chunks/chunks.go:314 — option type is a func over a
// private struct; each WithX returns one; the ctor folds them over a default.
type WriterOption func(*writerOptions)

func WithSegmentSize(segmentSize int64) WriterOption {
	return func(o *writerOptions) {
		if segmentSize <= 0 {
			segmentSize = DefaultChunkSegmentSize
		}
		o.segmentSize = segmentSize
	}
}

func NewWriter(dir string, opts ...WriterOption) (*Writer, error) {
	options := &writerOptions{segmentSize: DefaultChunkSegmentSize}
	for _, opt := range opts {
		opt(options)
	}
	...
}
```

Why god-tier: callers pass only what they care about; defaults stay in one
place; new options are added without breaking any existing call site. Note the
option itself even validates/normalizes (`<= 0` → default).

### 5. Extract a helper when it earns a name; design it to be shared

Extract logic into a helper when (a) it repeats, or (b) naming it makes the
caller read better. A god-tier shared helper takes the *minimal* inputs and
returns a *composable* result — like chi's `WithValue`, a tiny reusable unit
that returns a middleware closure:

```go
// chi-master/middleware/value.go:8 — minimal inputs (key,val), returns a
// reusable http.Handler decorator. Small, pure, composable.
func WithValue(key, val interface{}) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), key, val))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
```

Why god-tier: it depends on nothing but its arguments, so it's reusable and
testable in isolation. Don't extract a helper that needs five caller-locals
passed in — that coupling means it isn't really separable.

### 6. Receiver choice: pointer vs value, and be consistent

Use a pointer receiver if the method mutates the receiver, the struct is large,
or *any* method on the type needs a pointer (then make them all pointers for
consistency). Use a value receiver for small, immutable value types.

```go
// chi-master/mux.go:109 — *Mux: methods mutate routing state, so every method
// uses a pointer receiver. Consistency across the method set.
func (mx *Mux) Handle(pattern string, handler http.Handler) { ... }
```

Why god-tier: a mixed receiver set (some value, some pointer) is a classic
source of subtle bugs (lost mutations, accidental copies of a `sync.Mutex`).
Pick one per type.

## Anti-Patterns

| Smell | Fix |
|-------|-----|
| Param typed as concrete `*os.File` | Type it `io.Reader`/`io.Writer` (Principle 1) |
| `ctx` in the middle / not present | `ctx` first; thread it through (Principle 2) |
| `f(x, true, false, true)` | Named functions / options / enum (Principle 3) |
| Constructor with 8 positional args | `opts ...Option` (Principle 4) |
| Helper taking 6 caller-locals | Don't extract — the coupling says it's not separable (Principle 5) |
| Mixed value+pointer receivers on one type | One receiver kind per type (Principle 6) |
| Returning an interface "to be flexible" | Return the concrete struct; let callers narrow (Principle 1) |

## Checklist

- [ ] Each parameter is the smallest interface that does the job.
- [ ] Functions return concrete structs, not interfaces, unless an interface is genuinely needed.
- [ ] `context.Context` is the first parameter of every fallible/cancelable call; `error` is last.
- [ ] No boolean-trap parameters at any call site.
- [ ] Growable optional config uses functional options with defaults in one place.
- [ ] Extracted helpers depend only on their arguments (no smuggled caller state).
- [ ] One receiver kind (value or pointer) per type, chosen by mutation/size.

## Related

- [[go-clean-code]] — the prose inside the function (this skill is its shape).
- [[go-error-handling]] — designing the `error` return you place last.
- [[go-concurrency]] — why `context.Context` is passed, never stored.
- [[go-performance]] — when a hot signature should avoid interface boxing.
