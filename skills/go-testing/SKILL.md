---
name: go-testing
description: "God-tier Go testing — table-driven subtests, t.Parallel isolation, benchmarks with ReportAllocs, fuzzing, golden files, t.Helper, and right-sized mocking. Use when writing any _test.go, benchmark, or fuzz target. If you see copy-pasted test funcs, a parser tested on a few hand-picked inputs, helpers without t.Helper(), or over-mocking → apply this skill. For what to benchmark → See `go-performance`; for the -race detector → See `go-concurrency`. Do NOT use for non-test code."
origin: god_code
---

# Go Testing & Benchmarking

Part of the **God-Tier Go** set. This goes deeper than "write a test" — it's how
the stdlib, Prometheus, and Kubernetes get *coverage that catches real bugs*:
exhaustive cases, fuzzing, and benchmarks that gate performance. Every example is
cited from real code in this repo.

## When to Activate

- Writing or reviewing any `_test.go` file.
- Adding a benchmark for a hot path (pairs with [[go-performance]]).
- Testing a parser, encoder, or anything taking untrusted bytes (fuzz it).
- A bug slipped through — the fix needs a reproducing test.

## Principles

### 1. Table-driven tests with `t.Run` subtests

Don't copy-paste near-identical test functions. Express cases as a slice of
structs and loop, running each as a named subtest so failures point at the exact
case. The stdlib's `netip` tests are the canonical form:

```go
// gostd/net/netip/netip_test.go:138 — one row per case; t.Run names the subtest
// after its input so a failure says exactly which case broke.
for _, test := range validIPs {
	t.Run(test.in, func(t *testing.T) {
		got, err := ParseAddr(test.in)
		...
	})
}
```

Why god-tier: adding a case is one line; failures are isolated and named; the
test reads as a spec of behavior. This is the default shape for Go tests.

### 2. Isolate and parallelize with `t.Parallel()`

Subtests that don't share mutable state should declare `t.Parallel()` — it
surfaces hidden coupling and speeds the suite. Each subtest must own its data
(no shared loop variable capture, no shared globals).

```go
for _, tc := range cases {
	t.Run(tc.name, func(t *testing.T) {
		t.Parallel()        // runs concurrently with sibling subtests
		got := Process(tc.in)
		if got != tc.want {
			t.Errorf("Process(%q) = %v, want %v", tc.in, got, tc.want)
		}
	})
}
```

Why god-tier: parallel subtests prove isolation and run faster. If they flake
under `-race`, you've found a real bug — see [[go-concurrency]].

### 3. Benchmarks gate the hot path — and report allocations

Performance work without a benchmark is unverifiable (see [[go-performance]]).
Use nested `b.Run` for variants, `b.ReportAllocs()` to track allocations, and
`b.ResetTimer()` after setup. Prometheus benchmarks its query path this way:

```go
// prometheus-main/tsdb/querier_bench_test.go:179
for _, c := range cases {
	b.Run(c.name, func(b *testing.B) {
		b.ReportAllocs()   // allocation counts in the output
		b.ResetTimer()     // exclude setup above from timing
		for b.Loop() {     // modern loop; b.N still works on older Go
			p, err := PostingsForMatchers(ctx, ir, c.matchers...)
			require.NoError(b, err)
		}
	})
}
```

Why god-tier: `ReportAllocs` turns "feels faster" into numbers you compare with
`benchstat`; `ResetTimer` keeps setup out of the measurement.

### 4. Fuzz anything that parses untrusted input

For parsers, decoders, and byte-handling code, table tests can't cover the input
space — fuzzing can. Seed the corpus with `f.Add`, then assert invariants in
`f.Fuzz`. The stdlib fuzzes JSON unmarshaling:

```go
// gostd/json/fuzz_test.go:15 — seed a realistic corpus, then let the fuzzer
// mutate it; the body asserts "never panic / round-trips" over all inputs.
func FuzzUnmarshalJSON(f *testing.F) {
	f.Add([]byte(`{ "object": { "slice": [1, 2.0, "3"] } }`))
	f.Fuzz(func(t *testing.T, b []byte) {
		var v any
		if err := Unmarshal(b, &v); err != nil {
			return // invalid input is fine; a *panic* would not be
		}
	})
}
```

Why god-tier: fuzzing finds the malformed input you'd never think to type. Run
`go test -fuzz=FuzzUnmarshalJSON`; commit any crasher to `testdata/`.

### 5. Golden files for large/structured output

When the expected output is big (rendered config, serialized data), store it in
`testdata/` and compare against the golden file, with a `-update` flag to
regenerate. `testdata/` is special: the Go tool ignores it for builds.

```go
// Pattern: compare against golden, regenerate with -update.
got := render(input)
golden := filepath.Join("testdata", tc.name+".golden")
if *update {
	os.WriteFile(golden, got, 0o644)
}
want, _ := os.ReadFile(golden)
if !bytes.Equal(got, want) {
	t.Errorf("output differs from %s", golden)
}
```

Why god-tier: golden files keep large expectations out of the test source and
make intentional changes a reviewable diff.

### 6. Helpers call `t.Helper()`; mock only what's expensive

A test helper that does assertions must call `t.Helper()` so failures report the
*caller's* line, not the helper's. The stdlib does this consistently:

```go
// gostd/json/scanner_test.go:213 — t.Helper() so a failure points at the
// test that called this helper, not at the helper's own line.
func ... (t *testing.T, ...) {
	t.Helper()
	...
}
```

And don't mock what's cheap and real. Prefer a real `bytes.Buffer`, a real temp
dir (`t.TempDir()`), an in-memory store. Mock only slow/external dependencies
(network, clock, paid APIs). Over-mocking tests the mock, not the code.

Why god-tier: accurate failure locations and real dependencies mean the test
actually exercises production behavior.

## Anti-Patterns

| Smell | Fix |
|-------|-----|
| Copy-pasted near-identical `TestX1`/`TestX2` | Table-driven + `t.Run` (Principle 1) |
| `t` shared mutable state across subtests | Own data per subtest; `t.Parallel()` (Principle 2) |
| "It's faster now" with no benchmark | `testing.B` + `ReportAllocs` + `benchstat` (Principle 3) |
| Parser tested only on hand-picked inputs | Add a `FuzzXxx` (Principle 4) |
| 200-line expected string literal in the test | Golden file in `testdata/` (Principle 5) |
| Assertion helper without `t.Helper()` | Add `t.Helper()` (Principle 6) |
| Mocking a `bytes.Buffer` / filesystem | Use the real thing / `t.TempDir()` (Principle 6) |
| Testing one happy case only | Cover error + boundary rows in the table (Principle 1) |

## Checklist

- [ ] New logic is covered by a table-driven test with named subtests.
- [ ] Error and boundary cases are rows in the table, not just the happy path.
- [ ] Independent subtests call `t.Parallel()` and own their data.
- [ ] Hot paths have a `testing.B` with `b.ReportAllocs()` and `b.ResetTimer()`.
- [ ] Code parsing untrusted input has a `FuzzXxx`; crashers saved to `testdata/`.
- [ ] Large expected outputs use golden files with an `-update` flag.
- [ ] Assertion helpers call `t.Helper()`.
- [ ] Only slow/external deps are mocked; cheap deps are real.
- [ ] Suite passes under `go test -race ./...`.

## Related

- [[go-performance]] — what to benchmark and how to read `benchstat`.
- [[go-concurrency]] — the `-race` detector and concurrent test isolation.
- [[go-error-handling]] — assert error branches with `errors.Is`/`As`.
- [[go-clean-code]] — tests are code; keep them readable too.
