---
name: go-performance
description: "God-tier Go performance — benchmark first, reuse allocations with sync.Pool, preallocate slices/maps, mind escape analysis, zero-copy buffers, avoid hot-path allocation and reflection. Use when a benchmark/profile shows a bottleneck or you're reviewing a hot path. If you see optimization without a benchmark, append-to-nil in a known-size loop, per-request allocations, or string concat in a loop → apply this skill. For benchmark mechanics → See `go-testing`; for pool safety → See `go-concurrency`. Do NOT micro-optimize before profiling."
origin: god_code
---

# Go Performance & Allocation

Part of the **God-Tier Go** set. Go is fast by default; the wins come from *not
allocating* on hot paths and from measuring before you touch anything. This
skill distills how Prometheus — a system whose hot path runs millions of times a
second — squeezes performance. Every example is cited from real code in this
repo.

> The Prometheus maintainers' own rule (`prometheus-main/AGENTS.md`): *"Reuse
> allocations in hot paths where possible (slices, buffers)"* and *"Performance
> improvements require a benchmark that demonstrates the improvement."* This
> skill operationalizes exactly that.

## When to Activate

- Writing or reviewing a hot path (per-request, per-sample, tight loop).
- A benchmark/profile shows allocations or CPU you want to cut.
- Building a buffer/object that's created and discarded at high frequency.
- Tempted to optimize — start here so you measure first.

## Principles

### 1. Measure first — never optimize without a benchmark

Optimization without a benchmark is guessing, and guesses make code worse. Write
a `testing.B`, capture allocations, optimize, compare with `benchstat`. The
discipline is non-negotiable (see [[go-testing]] for benchmark mechanics).

```go
func BenchmarkEncode(b *testing.B) {
	b.ReportAllocs()   // make allocation regressions visible
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Encode(sample)
	}
}
```

Why god-tier: you optimize the *measured* hot path, prove the win with numbers,
and catch regressions. Everything below is justified only by a benchmark.

### 2. Reuse allocations on hot paths with `sync.Pool`

The single biggest GC win is not allocating per-operation. Pool and reuse
objects across calls. Prometheus pools chunk objects on its read path:

```go
// prometheus-main/tsdb/head_read.go:643 — reuse a *memChunk from the pool
// instead of allocating one per read.
mc := memChunkPool.Get().(*memChunk)
...
// prometheus-main/tsdb/head_read.go:573 — return it for the next caller.
h.memChunkPool.Put(c)
```

Why god-tier: `Get`/`Put` turns N allocations into ~0 steady-state. **Contract:**
a pooled object's contents belong to the current holder — callers must reset
before reuse and must not retain references after `Put`.

### 3. A type-safe pool avoids the `interface{}` allocation (SA6002)

Putting a non-pointer into `sync.Pool` boxes it and *allocates* — the staticcheck
SA6002 warning. Prometheus solves it with a generic, zero-allocation wrapper:

```go
// prometheus-main/util/zeropool/pool.go:21 — "a type-safe pool of items that
// does not allocate pointers to items." Zero value is valid and usable.
type Pool[T any] struct {
	items    sync.Pool
	pointers sync.Pool
}
```

And for variably-sized buffers, a bucketed pool returns a slice that *fits*:

```go
// prometheus-main/util/pool/pool.go:22 — bucketed pool for variably sized
// byte slices; Get(sz) returns the smallest bucket >= sz.
type Pool struct {
	buckets []sync.Pool
	sizes   []int
	make    func(int) any
}
```

Why god-tier: these patterns eliminate the hidden boxing allocation that a naive
`sync.Pool` of values incurs, while staying type-safe and zero-value-usable.

### 4. Preallocate slices and maps when the size is known

`append` to a `nil` slice repeatedly causes log-many reallocations and copies.
If you know (or can bound) the length, allocate the capacity up front. Prometheus
does this everywhere on the query path:

```go
// prometheus-main/tsdb/querier.go:526 — capacity known from the input, so
// allocate once, append without regrowing.
values := make([]string, 0, len(indexes))
// prometheus-main/tsdb/querier.go:610
chks := make([]chunks.Meta, 0, nChks)
```

Why god-tier: one allocation instead of O(log n) regrowths, and zero copies. Same
for maps: `make(map[K]V, n)`.

### 5. Mind allocation cost at call boundaries — know what escapes

Some innocuous-looking stdlib calls allocate; a god-tier author knows which and
comments the cost. chi annotates the exact allocation count of its per-request
context wiring:

```go
// chi-master/mux.go:86 — the author measured and documented the cost.
// NOTE: r.WithContext() causes 2 allocations and context.WithValue() causes 1 allocation
r = r.WithContext(context.WithValue(r.Context(), RouteCtxKey, rctx))
```

Why god-tier: knowing what escapes to the heap (`go build -gcflags=-m`) lets you
keep short-lived values on the stack and avoid pointer-to-local leaks. Returning
a pointer to a local forces a heap allocation — return values, or reuse buffers.

### 6. Zero-copy: reuse buffers, avoid `[]byte`↔`string` churn

In tight loops, repeated string concatenation and `[]byte("...")`/`string(b)`
conversions each allocate. Use `strings.Builder`/`bytes.Buffer`, write into a
caller-provided buffer, and avoid converting back and forth.

```go
// Bad — each += allocates a new backing array.
out := ""
for _, s := range parts { out += s }

// Good — one growing buffer, amortized allocation.
var b strings.Builder
b.Grow(estimatedLen) // preallocate when you can estimate
for _, s := range parts { b.WriteString(s) }
out := b.String()
```

Why god-tier: buffer reuse and avoiding conversions cut both allocations and
copies. **Contract:** when you hand a reused buffer to an interface, document
that callers must copy and must not retain it (the Prometheus AGENTS.md rule).

## Anti-Patterns

| Smell | Fix |
|-------|-----|
| Optimizing with no benchmark | Write `testing.B` + `benchstat` first (Principle 1) |
| Allocating a fresh object per request | Pool and reuse with `sync.Pool` (Principle 2) |
| `sync.Pool` of non-pointer values | Pool pointers / use a zeropool wrapper (Principle 3) |
| `append` to `nil` in a known-size loop | `make([]T, 0, n)` up front (Principle 4) |
| Returning `&local` on a hot path | Return the value or write into a buffer (Principle 5) |
| `s += x` in a loop | `strings.Builder` with `Grow` (Principle 6) |
| `string(b)` / `[]byte(s)` churn in a loop | Keep one representation; reuse buffers (Principle 6) |
| Reflection/`interface{}` boxing in the inner loop | Concrete types / generics (Principle 5) |

## Checklist

- [ ] The optimized path has a `testing.B` benchmark with `b.ReportAllocs()`.
- [ ] Improvements are proven with before/after numbers (`benchstat`).
- [ ] High-frequency objects are pooled and reset on reuse.
- [ ] No non-pointer values are put into a raw `sync.Pool`.
- [ ] Slices/maps with known size are preallocated with capacity.
- [ ] No pointer-to-local escapes on the hot path (`-gcflags=-m` checked).
- [ ] No string concatenation or `[]byte`↔`string` churn in tight loops.
- [ ] Reused buffers handed across interfaces document the no-retain contract.

## Related

- [[go-testing]] — benchmarks, `ReportAllocs`, and `benchstat` mechanics.
- [[go-concurrency]] — `sync.Pool` safety and lock contention.
- [[go-clean-code]] — only trade clarity for speed where a benchmark demands it.
- [[go-function-design]] — avoiding interface boxing in hot signatures.
