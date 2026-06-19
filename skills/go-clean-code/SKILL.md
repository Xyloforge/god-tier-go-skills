---
name: go-clean-code
description: "Write Go that reads obviously — clarity over cleverness, early returns/guard clauses, precise naming, minimal state, honest doc comments. Load this FIRST for any Go work; it routes to the specialists. If you see deep nesting, clever one-liners, or undocumented exports → apply this skill. For naming depth → See `go-naming`; for function shape → See `go-function-design`. Do NOT use for performance micro-tuning (→ See `go-performance`)."
origin: god_code
---

# Go Clean Code

The foundational skill of the **God-Tier Go** set. Every claim below is cited
from real, production code in this repo (chi, the Go stdlib, Kubernetes,
Prometheus, Vault). Load this first, then drill into a specialist via the
cross-links in **Related**.

The bar: *a senior engineer reading the diff cold should never have to pause and
decode it.* If they pause, the code is not done.

## When to Activate

- Writing any new Go function, type, or package.
- Reviewing or refactoring Go for readability.
- A function "works" but feels hard to read, deeply nested, or cleverly terse.
- Before declaring any Go change complete (run the Checklist).

## Principles

### 1. Clarity over cleverness

Go optimizes for the reader, not the writer. Prefer the obvious form even when a
terser one exists. Look at how chi's request entry point reads top-to-bottom
with no trickery — each block does one thing and says why:

```go
// chi-master/mux.go:63 — the hottest path in the router, still written plainly.
func (mx *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Ensure the mux has some routes defined on the mux
	if mx.handler == nil {
		mx.NotFoundHandler().ServeHTTP(w, r)
		return
	}
	// Check if a routing context already exists from a parent router.
	rctx, _ := r.Context().Value(RouteCtxKey).(*Context)
	if rctx != nil {
		mx.handler.ServeHTTP(w, r)
		return
	}
	...
}
```

Why god-tier: the busiest function in a high-performance router has *zero*
cleverness. Performance lives in the data structures, not in dense syntax.

### 2. Guard clauses and early returns — keep the happy path flat

Handle the exceptional case, `return`, and let the main logic live unindented at
the left margin. The two `return`s above (`mux.go:65` and `mux.go:72`) are guard
clauses: each peels off one edge case so the normal flow never nests.

```go
// Bad — happy path buried under nesting
func process(r *Request) error {
	if r != nil {
		if r.Valid() {
			if r.Body != nil {
				return handle(r) // 3 levels deep to do the real work
			}
		}
	}
	return errInvalid
}

// Good — guard clauses, flat happy path
func process(r *Request) error {
	if r == nil || !r.Valid() {
		return errInvalid
	}
	if r.Body == nil {
		return errNoBody
	}
	return handle(r)
}
```

Why god-tier: nesting depth is a proxy for cognitive load. Flat code is scannable.

### 3. Naming — short where the scope is short, descriptive where it's wide

Go convention: receivers and short-lived locals get 1–3 letters; exported
identifiers get clear, caller-facing names. chi names its receiver `mx *Mux`
(used hundreds of times, so it stays terse) while exposing descriptive method
names like `Handle`, `Method`, `NotFound`:

```go
// chi-master/mux.go:109 — terse receiver, descriptive exported name.
func (mx *Mux) Handle(pattern string, handler http.Handler) { ... }
```

Why god-tier: name length should be proportional to scope. `i` in a 3-line loop
is clearer than `index`; a package-level export deserves a full, honest name.
For the *shape* of names and signatures, see [[go-function-design]].

### 4. Comment why, not what — and document every export

Don't narrate code the reader can already read. Comment intent, invariants, and
the non-obvious. The Go stdlib sets the standard: exported identifiers get a doc
comment that *starts with the identifier's name* and states the contract:

```go
// gostd/io/io.go:86 — the doc comment defines the contract, including the
// non-obvious invariant a reader could never guess from the signature.
//
// Reader is the interface that wraps the basic Read method.
// ...
// Implementations must not retain p.
type Reader interface {
	Read(p []byte) (n int, err error)
}
```

Why god-tier: the comment carries the *contract* (`must not retain p`) that the
type signature cannot. That is the only kind of comment worth writing.

### 5. One responsibility per function; keep them small

A function should do one thing at one level of abstraction. chi's `Handle`
(`mux.go:109`) does exactly one thing — register a route — and delegates the
method-specific work to `mx.handle`. When a function mixes parsing, validation,
and I/O, split it. *When* to extract and what the helper's signature should be
is [[go-function-design]]; *that* it should be small and single-purpose is here.

### 6. Lean on the zero value; avoid needless state

Idiomatic Go types are useful at their zero value, which removes constructors,
nil checks, and init ceremony. `sync.Mutex`, `bytes.Buffer`, and `sync.Pool` all
work with `var x T`. Design your own types the same way:

```go
// Good — usable immediately, no constructor required.
var mu sync.Mutex
mu.Lock()

// Avoid — forcing a constructor for what could be a useful zero value.
mu := NewMutex() // ceremony with no payoff
```

Why god-tier: fewer states means fewer bugs. The best initializer is none.

## Anti-Patterns

| Smell | Fix |
|-------|-----|
| Pyramid of nested `if`s | Guard clauses + early return (Principle 2) |
| `// increment i by one` | Delete it; comment *why* or nothing (Principle 4) |
| Exported func with no doc comment | Add `// Name …` contract comment (Principle 4) |
| `data`, `data2`, `tmp`, `obj` | Name for role/scope (Principle 3) |
| 120-line function doing 4 jobs | Split by responsibility (Principle 5) |
| `NewFoo()` that only zeroes fields | Make the zero value work (Principle 6) |
| Clever one-liner needing a comment to decode | Write the boring three lines |

## Checklist

Before declaring Go code complete:

- [ ] No function nests more than ~3 levels — guard clauses applied.
- [ ] Every exported identifier has a doc comment starting with its name.
- [ ] Comments explain *why*/contracts, not *what* the code literally does.
- [ ] Receiver and short-lived local names are terse; exports are descriptive.
- [ ] Each function does one thing at one level of abstraction.
- [ ] No constructor exists solely to zero fields a zero value would handle.
- [ ] A cold reader can follow the happy path straight down the left margin.

## Related

- [[go-function-design]] — the *shape* of a function (this skill is the *prose inside* it).
- [[go-error-handling]] — clarity in the error path.
- [[go-concurrency]] — clear, leak-free goroutine code.
- [[go-performance]] — when (and only when) to trade clarity for speed.
- [[go-testing]] — proving the clarity holds under change.
