---
name: go-naming
description: "Go naming conventions — packages, exported/unexported identifiers, interfaces (-er), constructors (New), errors (Err/Error), booleans, receivers, getters, and acronyms. Use when naming anything new, reviewing names, or choosing between alternatives (New vs NewClient, ErrNotFound vs NotFoundError, isReady vs ready). If you see ALL_CAPS consts, Get-prefixed getters, stutter (http.HTTPServer), or a utils/ package → apply this skill. For the prose around names → See `go-clean-code`. Do NOT use for non-naming implementation questions."
origin: god_code
---

# Go Naming

Part of the **God-Tier Go** set. Names are the API. Go's conventions are tight
and near-universal across the stdlib and the codebases in this repo. This is a
thin router; the cited specifics live in `references/conventions.md`.

## When to Activate

- Naming a package, type, function, method, constant, variable, or error.
- Reviewing or refactoring identifiers; resolving a naming debate.
- Designing an interface, constructor, or sentinel error.

## Decision Guide

| You're naming… | Convention | Example (cited) |
|----------------|-----------|-----------------|
| A package | short, lower-case, no plural, no `util` | `io`, `tsdb`, `middleware` |
| An exported identifier | `MixedCaps`, no underscores | `ServeHTTP`, `WithSegmentSize` |
| A single-method interface | verb + `-er` | `Reader`, `Writer`, `Closer` |
| A constructor | `New` or `NewТype` | `NewWriter`, `New` |
| A sentinel error | `Err` prefix | `ErrNotInitialized` |
| A custom error type | `Error` suffix | `errorWithStack` |
| A boolean | adjective/predicate, no `is`/`Get` noise | `useUncachedIO` |
| A receiver | 1–3 letters, consistent per type | `mx *Mux` |

## Core Rules

1. **Packages:** short, lowercase, singular, evocative. Never `utils`, `helpers`, `common`, `base`. The package name is part of every call (`tsdb.NewWriter`), so avoid stutter (`http.Server`, not `http.HTTPServer`).
2. **MixedCaps, not snake_case or ALL_CAPS** — including constants. Exported = capitalized; unexported = lowercase. Underscores appear only in `_test.go` test names.
3. **Interfaces:** one-method interfaces are `verb`+`er` (`Reader`). Name interfaces for behavior, not implementation.
4. **Constructors:** `New` when the package returns one main type; `NewТype` when several. Return the concrete struct.
5. **Errors:** sentinels are `var ErrX = errors.New(...)`; error *types* are `XError`/`...Error`. Error strings are lowercase, no punctuation (see `go-error-handling`).
6. **Booleans & getters:** predicates read as assertions (`ok`, `useUncachedIO`); getters drop `Get` (`obj.Name()`, not `obj.GetName()`). Setters keep `Set`.
7. **Receivers:** short and identical across the whole method set; never `this`/`self`.
8. **Acronyms keep one case:** `URL`, `ID`, `HTTP` — `userID`, `ServeHTTP`, `parseURL`, never `userId` or `ServeHttp`.

## Anti-Patterns

| Smell | Fix |
|-------|-----|
| `package utils` / `helpers` / `common` | Name by domain; split by responsibility (Rule 1) |
| `http.HTTPServer`, `tsdb.TSDBWriter` (stutter) | Drop the redundant prefix (Rule 1) |
| `const MAX_SIZE = 100` | `const MaxSize = 100` (Rule 2) |
| `type UserInterface interface` | Name for behavior, e.g. `UserStore` / `-er` (Rule 3) |
| `func (this *Mux) ...` | `func (mx *Mux) ...` (Rule 7) |
| `obj.GetName()` | `obj.Name()` (Rule 6) |
| `userId`, `ServeHttp`, `apiUrl` | `userID`, `ServeHTTP`, `apiURL` (Rule 8) |
| `NotFoundError` as a sentinel value | `ErrNotFound` (Rule 5) |

## Checklist

- [ ] Package names are short, lowercase, singular, non-generic (no `utils`).
- [ ] No stutter between package and identifier names.
- [ ] All identifiers use `MixedCaps`; no `ALL_CAPS`, no `snake_case` (outside test names).
- [ ] One-method interfaces use the `-er` form; interfaces named for behavior.
- [ ] Constructors are `New`/`NewТype` and return concrete types.
- [ ] Sentinels are `ErrX`; error types end in `Error`.
- [ ] Getters omit `Get`; receivers are short and consistent per type.
- [ ] Acronyms keep consistent case (`ID`, `URL`, `HTTP`).

## Deep Dives

- `references/conventions.md` — each rule with real identifiers from the corpus (cited: stdlib, chi, Prometheus, Vault, k8s).

## Related

- → See `go-clean-code` — naming as part of overall readability.
- → See `go-error-handling` — error value vs type naming and message style.
- → See `go-structs-interfaces` — naming interfaces for behavior.
- → See `go-project-layout` — package naming and boundaries.
