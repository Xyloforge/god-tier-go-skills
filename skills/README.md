# God-Tier Go

A set of **fourteen** Go skills for AI agents (and humans) who want to write
**quality, high-performance, idiomatic Go** — code no reviewer can question.

What makes this set different from generic "Go best practices": **every
non-trivial claim is backed by a traceable citation** to real, production code
sitting in this repo — the Go standard library, [chi](chi-master),
[Kubernetes](kubernetes), [Moby/Docker](moby), [Prometheus](prometheus-main),
and [HashiCorp Vault](vault-main). When a skill shows you a pattern, it points at
the exact file and line where a world-class team does it. No invented examples.

## The skills

**Foundations** — load `go-clean-code` first; it routes to the rest.

| Skill | One-liner | Load when… |
|-------|-----------|-----------|
| **[go-clean-code](go-clean-code/SKILL.md)** | Clarity over cleverness — early returns, naming, minimal state, doc comments. | Writing/reviewing *any* Go. **Start here.** |
| **[go-naming](go-naming/SKILL.md)** | Packages, interfaces (`-er`), `New`, `Err`, receivers, acronyms. | Naming anything or settling a naming debate. |
| **[go-function-design](go-function-design/SKILL.md)** | The *shape* of a function — small interface params, context-first/error-last, options, receivers. | Designing a signature, constructor, or helper. |
| **[go-structs-interfaces](go-structs-interfaces/SKILL.md)** | *Type* design — tiny consumer-side interfaces, embedding, useful zero values. | Declaring an interface or struct. |
| **[go-project-layout](go-project-layout/SKILL.md)** | `cmd/`, `internal/`, domain packages, module paths, no `utils/`. | Starting a repo or placing a package/binary. |

**Correctness & runtime**

| Skill | One-liner | Load when… |
|-------|-----------|-----------|
| **[go-error-handling](go-error-handling/SKILL.md)** | Wrap with `%w`, sentinel vs typed errors, `errors.Is/As`, never swallow, retries. | Returning, wrapping, or matching errors. |
| **[go-context](go-context/SKILL.md)** | Cancellation, timeouts, `defer cancel()`, typed value keys. | Threading `ctx`, timeouts, or request values. |
| **[go-concurrency](go-concurrency/SKILL.md)** | Goroutine lifecycle, leak prevention, bounded `errgroup`, sync primitives. | Starting a goroutine, channel, or pool. |
| **[go-security](go-security/SKILL.md)** | `crypto/rand`, constant-time compare, TLS floor, input validation, injection. | Handling secrets, tokens, TLS, or untrusted input. |

**Performance & operability**

| Skill | One-liner | Load when… |
|-------|-----------|-----------|
| **[go-performance](go-performance/SKILL.md)** | Benchmark first, `sync.Pool` reuse, preallocation, escape analysis, zero-copy. | Touching a hot path or chasing allocations. |
| **[go-observability](go-observability/SKILL.md)** | Prometheus metrics + naming, structured `slog`, contextual logging, tracing. | Adding metrics, logs, or traces. |
| **[go-testing](go-testing/SKILL.md)** | Table-driven subtests, `t.Parallel`, benchmarks, fuzzing, golden files. | Writing any `_test.go` or benchmark. |

**Review & understanding**

| Skill | One-liner | Load when… |
|-------|-----------|-----------|
| **[go-adversarial-qa](go-adversarial-qa/SKILL.md)** | Break a function before prod does — Toddler (bad/empty/wrong-type input), Hoarder (scale/limits/frequency), Inversion (failing deps). Locates breaks; never fixes. | Reviewing/hardening a function or asked "what breaks this?" |
| **[explaining-code-logic](explaining-code-logic/SKILL.md)** | Explain gnarly logic to a human with an interactive, input-driven step-tracer instead of a wall of text. | Onboarding someone to unfamiliar/complex logic. |

## How to use them

1. **Always load [go-clean-code](go-clean-code/SKILL.md) first** — it's the
   foundation and routes you to the right specialist.
2. **Drill into the specialist** for the task at hand via the `→ See` pointers in
   each skill's frontmatter `description` and the **Related** section.
3. **Run the Checklist** at the bottom of each skill before declaring code done.

The skills are intentionally non-overlapping. Sharpest boundaries:
`go-clean-code` = the *prose inside* a function; `go-function-design` = a
function's *signature shape*; `go-structs-interfaces` = the *type definitions*;
`go-naming` = the *names*.

## Each skill's structure

Skills use **progressive disclosure**: a thin `SKILL.md` router (decision guide +
core rules + checklist) with heavy detail in `references/`, and copyable files in
`assets/`.

```
<skill>/
├── SKILL.md          router: frontmatter (router-style triggers + → See pointers),
│                     When to Activate, Decision Guide, Core Rules/Principles,
│                     Anti-Patterns, Checklist, Deep Dives, Related
├── references/*.md   deep dives with the cited snippets (loaded on demand)
└── assets/*          drop-in files (e.g. secure-http-server.go, Makefile, alert rules)
```

The original 6 foundations use a `Principles` section (statement → cited snippet
→ why); the newer skills use a `Decision Guide` + `Core Rules` router that points
into `references/`. Both carry the same When-to-Activate / Anti-Patterns /
Checklist / Related contract.

## Installing for auto-trigger

This repo is the **source of truth**. To make a Claude Code agent auto-discover
and trigger the skills, symlink each folder into `~/.claude/skills/`:

```bash
cd god-tier-go-skills/skills
for s in go-clean-code go-naming go-function-design go-structs-interfaces \
         go-project-layout go-error-handling go-context go-concurrency \
         go-security go-performance go-observability go-testing; do
  ln -s "$PWD/$s" "$HOME/.claude/skills/$s"
done
```

(Or copy instead of symlink, or wrap the `skills/` folder as a plugin.) Editing a
file here then updates the live skill automatically when symlinked.

## Provenance

Distilled from, and cite-traceable to, the codebases in this directory:
`gostd` (Go stdlib), `chi-master`, `kubernetes`, `moby`, `prometheus-main`,
`vault-main`. Several of these ship their own agent guides
(`prometheus-main/AGENTS.md`, `moby/CLAUDE.md`) whose maintainer rules
reinforce these skills.
