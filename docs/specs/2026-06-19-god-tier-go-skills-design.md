# God-Tier Go — Skill Set Design

**Date:** 2026-06-19
**Status:** Approved (design), pending implementation plan
**Author:** Claude (with rong)

## Goal

A set of **6 Go skills**, each distilled from real, traceable patterns in the
world-class codebases under `god_code/`, that an AI agent loads to generate
quality, high-performance, idiomatic Go — code that "no human can question."

This set is sharper and more authoritative than the existing generic
`golang-patterns` / `golang-testing` skills: every non-trivial claim is backed
by a cited snippet from production code (chi, Kubernetes, Moby, Prometheus,
Vault, or the curated Go stdlib slices).

## Reference corpus (`god_code/`)

| Folder | Project | Primary lessons mined |
|--------|---------|----------------------|
| `chi-master` | go-chi/chi | small interfaces, composition, router/radix tree, middleware design |
| `gostd` | Go stdlib (io, json, net, os, sync) | the canonical idioms — `io.Reader/Writer`, `sync` primitives |
| `kubernetes` | Kubernetes | concurrency at scale, error handling, context, testing breadth |
| `moby` | Moby / Docker Engine | concurrency, resource lifecycle, performance |
| `prometheus-main` | Prometheus | allocation discipline, hot-path performance, benchmarking |
| `vault-main` | HashiCorp Vault | error wrapping/resilience, defensive design |

## The 6 skills

| # | Skill (folder) | Distilled from | Domain |
|---|----------------|----------------|--------|
| 1 | `go-clean-code` | stdlib + all repos | **Foundational entry point.** Clarity over cleverness, naming, guard clauses / early returns, killing nesting, comment discipline, function/file size, single responsibility — the *prose* of the code. |
| 2 | `go-function-design` | chi, stdlib `io`/`net` | The *shape* of a function: signature design (params, returns, error position), when to extract a helper, designing shared/reusable functions, small interfaces as params, options pattern, receiver choice. |
| 3 | `go-error-handling` | Vault, Kubernetes | Wrapping / `errors.Is` / `errors.As`, sentinel vs typed errors, error position, retries/backoff, graceful degradation, never-swallow discipline. |
| 4 | `go-concurrency` | Kubernetes, Moby, stdlib `sync` | Goroutine lifecycle & leak prevention, context cancellation, channels vs mutexes, `errgroup`, worker pools, race-free patterns. |
| 5 | `go-performance` | Prometheus tsdb, stdlib, Moby | Hot-path allocation avoidance, `sync.Pool`, slice/map preallocation, escape analysis, zero-copy, benchmark-driven optimization. |
| 6 | `go-testing` | Kubernetes, stdlib | Table-driven tests, subtests, fuzzing, benchmarks, golden files, parallelism — god-tier depth beyond the existing `golang-testing`. |

### Boundary rules (no overlap)

- **`go-clean-code` = the prose inside a function** (is it readable/simple?).
- **`go-function-design` = the shape of the function** (should it exist, what
  goes in/out, is it reusable?).
- Each specialist owns its domain exclusively; shared concerns are stated once
  in the most relevant skill and cross-linked via `[[skill-name]]`.

## File structure

```
god_code/
├── skills/
│   ├── README.md                       # "God-Tier Go" index: each skill + when to load
│   ├── go-clean-code/SKILL.md
│   ├── go-function-design/SKILL.md
│   ├── go-error-handling/SKILL.md
│   ├── go-concurrency/SKILL.md
│   ├── go-performance/SKILL.md
│   └── go-testing/SKILL.md
└── (chi-master, gostd, kubernetes, moby, prometheus-main, vault-main)
```

## SKILL.md template (every skill follows this)

```markdown
---
name: <go-topic>
description: <one-line trigger an agent matches on>
origin: god_code
---

# <Title>

## When to Activate
- concrete signals (writing X, reviewing Y, optimizing Z)

## Principles
For each principle:
  - State the principle plainly.
  - Show a REAL cited snippet:  // <repo>: <path/file.go> — <what & why>
  - Explain why it is god-tier (the reasoning a senior would give).

## Anti-Patterns
- Wrong way vs right way, side by side.

## Checklist
- [ ] short pass/fail items the agent self-checks before declaring code done.

## Related
- [[other-skill]] links for drill-down.
```

## Authority rule (the thing that makes this "unquestionable")

Every non-trivial pattern is backed by a **traceable snippet** from one of the 6
repos, cited as `repo/path/file.go`. No invented examples. During writing, the
source is `grep`'d to pull genuine patterns; if a claimed pattern cannot be
found in the corpus, it is either removed or explicitly marked as a general
idiom (not a god_code citation).

## Discovery / installation note

Sitting in `god_code/skills/` these are **reference material**, not auto-loaded
Claude Code skills. The README documents that they can be symlinked or copied
into `~/.claude/skills/` (or wrapped as a plugin) to make an agent auto-trigger
them. The source of truth stays in `god_code/` per the owner's request.

## Out of scope (YAGNI)

- Package/module architecture (intentionally dropped — too abstract to action).
- Observability/metrics, generics deep-dives, memory-layout — may become future
  skills but are not in this set.
- Rewriting the existing `golang-patterns` / `golang-testing` skills.

## Success criteria

1. All 6 `SKILL.md` files + README exist and follow the template.
2. Every code example is real and cite-traceable to a file in `god_code/`.
3. Boundaries hold — no two skills cover the same ground.
4. Each skill ends with an actionable self-check checklist.
5. A cold agent loading `go-clean-code` can navigate to the right specialist.
```
