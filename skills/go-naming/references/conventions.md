# Go Naming Conventions — with real identifiers

Every rule below is shown with an identifier that actually ships in this repo.

## Packages

Short, lowercase, singular, evocative — and meaningful at the call site, because
the package name prefixes every exported call.

- `gostd/io` → `io.Reader`, `io.Copy`. One syllable, says everything.
- `prometheus-main/tsdb` → `tsdb.NewWriter`. Domain term, not `database`.
- `chi-master/middleware` → `middleware.WithValue`.

Avoid `utils`, `helpers`, `common`, `base`, `misc` — they describe nothing and
become dumping grounds. Name a package for what it *provides*. Avoid stutter: the
type in package `http` is `Server`, giving `http.Server` — never `http.HTTPServer`.

## MixedCaps everywhere — including constants

Go uses `MixedCaps`/`mixedCaps`, never `snake_case` or `ALL_CAPS`.

```go
// prometheus-main/tsdb/chunks/chunks.go:316 — exported func, MixedCaps.
func WithSegmentSize(segmentSize int64) WriterOption
```

```go
// prometheus-main/tsdb/chunks/chunks.go:310 — unexported field, mixedCaps,
// boolean named as a predicate (no "is"/"flag" noise).
useUncachedIO bool
```

Constants follow the same rule: `MaxSize`, not `MAX_SIZE`. Underscores appear
only in test function names (`Test_parseURL_invalid`).

## Interfaces — behavior, and `-er` for one method

```go
// gostd/io/io.go:86,99,107 — the canonical one-method, verb+er interfaces.
type Reader interface { Read(p []byte) (n int, err error) }
type Writer interface { Write(p []byte) (n int, err error) }
type Closer interface { Close() error }
```

Name interfaces for what they *do*, not who implements them. Avoid the
`SomethingInterface`/`ISomething` styles from other languages.

## Constructors — `New` / `NewType`, returning a concrete type

```go
// prometheus-main/tsdb/chunks/chunks.go:336 — NewWriter: the package exposes
// several constructible types, so the type is in the name.
func NewWriter(dir string, opts ...WriterOption) (*Writer, error)
```

When a package builds one primary type, just `New` (e.g. `zeropool.New`).

## Errors — `Err` values, `Error` types

```go
// vault-main/sdk/database/helper/connutil/connutil.go:12 — sentinel value.
var ErrNotInitialized = errors.New("connection has not been initialized")
```

```go
// kubernetes/cmd/kubeadm/app/util/errors/errors.go:41 — error *type* (suffix).
type errorWithStack struct { ... }
func (e errorWithStack) Error() string { ... }
```

So: sentinel **values** get an `Err` prefix; error **types** get an `Error`
suffix. (Message-string style — lowercase, no punctuation — is in
`go-error-handling`.)

## Receivers — short, consistent, never `this`/`self`

```go
// chi-master/mux.go — every method on *Mux uses the same 2-letter receiver.
func (mx *Mux) Handle(pattern string, handler http.Handler)
func (mx *Mux) Use(middlewares ...func(http.Handler) http.Handler)
```

Pick a 1–3 letter receiver derived from the type and use it on *every* method of
that type. Mixing `mx` and `m` and `mux` across methods is a smell.

## Getters and acronyms

- Getters omit `Get`: a field `name` is exposed as `Name()`, not `GetName()`.
  Setters keep `Set`: `SetName(...)`.
- Acronyms keep a single case throughout the word: `ID`, `URL`, `HTTP`, `API`.
  So `userID` (not `userId`), `ServeHTTP` (not `ServeHttp`), `parseURL` (not
  `parseUrl`). chi's `ServeHTTP` (`mux.go:63`) is the canonical example.
