# Interface Design

## Keep them tiny — one method composes with everything

The Go standard library is built from one-method interfaces, and that is exactly
why it composes so well:

```go
// gostd/io/io.go:86,99,107 — three one-method interfaces that the entire
// ecosystem builds on. Anything that Reads works with anything that Writes.
type Reader interface { Read(p []byte) (n int, err error) }
type Writer interface { Write(p []byte) (n int, err error) }
type Closer interface { Close() error }
```

Larger behaviors are *composed* from these, not declared as one fat interface:

```go
// gostd/io — composition by embedding small interfaces.
type ReadCloser interface {
	Reader
	Closer
}
```

Why god-tier: a one-method interface has a huge set of implementers and is
trivial to fake in tests. Every method you add shrinks that set. Aim for 1–3
methods; if an interface has many, it's probably a struct in disguise.

## Define interfaces at the consumer, not the implementer

The Go idiom is: **the package that *uses* a behavior declares the interface for
it**, listing only the methods it actually calls. The implementing package
exposes a concrete type with concrete methods and imports nothing extra.

```go
// In the CONSUMER package — declare exactly what you need:
type userStore interface {
	UserByID(ctx context.Context, id string) (*User, error)
}

func NewHandler(s userStore) *Handler { return &Handler{store: s} }
```

The implementer (a `*PostgresStore`, say) never references `userStore` — it just
has a `UserByID` method, and satisfies the interface implicitly. Benefits:

- The consumer depends on a minimal, local abstraction it controls.
- The implementer has no import dependency on its consumers.
- Tests in the consumer package fake `userStore` in a few lines.

## Accept interfaces, return concrete types

Abstract the *inputs*; keep the *outputs* concrete so callers get the full type
and can decide what to narrow:

```go
// Good — input abstracted, output concrete.
func NewScanner(r io.Reader) *Scanner   // takes any Reader, returns *Scanner

// Avoid — returning an interface hides capability and forces type assertions.
func NewScanner(r io.Reader) ScannerLike
```

## Verify satisfaction at compile time

Because satisfaction is implicit, a refactor can silently break it. Pin it:

```go
// Compile-time guarantee that *Mux implements http.Handler; fails to build if
// the method set ever drifts. (chi relies on *Mux satisfying http.Handler.)
var _ http.Handler = (*Mux)(nil)
```
