---
name: go-project-layout
description: "Structuring a Go module/repository — cmd/ for binaries, internal/ for private code, package boundaries by domain, go.mod/module path, and avoiding god-packages. Use when starting a repo, adding a package or binary, or deciding where code belongs. If you see a utils/ or common/ package, business logic in main(), or import cycles → apply this skill. For package *names* → See `go-naming`. Do NOT use for in-package type design (→ See `go-structs-interfaces`)."
origin: god_code
---

# Go Project Layout

Part of the **God-Tier Go** set. How a module is organized determines what can
import what, what's hideable, and how it grows. This distills the layouts of the
big repos in this corpus (Kubernetes, Moby, Prometheus, chi). Thin router; cited
depth in `references/structure.md`.

## When to Activate

- Starting a new Go module or restructuring an existing one.
- Adding a package or a binary; deciding which directory it belongs in.
- Resolving import cycles or a package that's become a junk drawer.

## Decision Guide

| You're adding… | Put it in… | Why |
|----------------|-----------|-----|
| A runnable binary | `cmd/<name>/main.go` | thin entrypoint, one per binary |
| Code no one outside should import | `internal/...` | compiler-enforced privacy |
| Domain logic | a package named for the domain | high cohesion |
| Code meant for external import | a top-level package (library) | part of your public API |
| A "misc" helper | the package that uses it | there is no `utils/` |

## Core Rules

1. **`cmd/<binary>/main.go` is thin.** `main` parses flags/env, wires dependencies, and starts the app — no business logic. Every real repo here does this.
2. **`internal/` for everything private.** The Go toolchain *forbids* imports of `internal/` from outside its parent module subtree — use it aggressively to keep your public surface tiny.
3. **Organize by domain, not by layer.** Group `order/`, `billing/`, `tsdb/` — not `models/`, `controllers/`, `services/`. Code that changes together lives together.
4. **No `utils`/`common`/`helpers`/`shared` packages.** They have no cohesion and become cycle magnets. A helper lives in the package that uses it (see `go-naming`).
5. **One module per repo by default** (`go.mod` at root). Reach for multi-module only with a concrete reason (separately versioned public SDK).
6. **Keep the dependency graph acyclic and one-directional.** Inner/domain packages don't import outer/transport packages. If you hit an import cycle, the boundary is wrong — extract a small interface at the consumer (see `go-structs-interfaces`).
7. **`pkg/` is optional, not mandatory.** Use a top-level package directly; only add a `pkg/` umbrella if it genuinely clarifies a large repo. Don't cargo-cult it.

## Anti-Patterns

| Smell | Fix |
|-------|-----|
| Business logic inside `main()` | Move to a package; keep `cmd/` thin (Rule 1) |
| Everything exported / no `internal/` | Hide private code in `internal/` (Rule 2) |
| `models/`, `services/`, `handlers/` by layer | Reorganize by domain (Rule 3) |
| `utils/`, `common/`, `helpers/` | Co-locate with the user; name by domain (Rule 4) |
| Import cycle between packages | Invert with a consumer-side interface (Rule 6) |
| `pkg/` added reflexively to a tiny repo | Use top-level packages (Rule 7) |

## Checklist

- [ ] Each binary is `cmd/<name>/main.go` and `main` only wires + starts.
- [ ] Private code is under `internal/`; the public surface is intentional.
- [ ] Packages are organized by domain, not technical layer.
- [ ] No `utils`/`common`/`helpers`/`shared` packages.
- [ ] Single module at repo root unless a real need justifies more.
- [ ] The import graph is acyclic and points one direction (domain ← transport).
- [ ] `pkg/` exists only if it earns its keep.

## Deep Dives

- `references/structure.md` — annotated layouts of Kubernetes, Moby, Prometheus, chi, with module paths (cited).

## Assets

- `assets/Makefile` — a starter Makefile (build/test/lint/vet/cover) for a `cmd/`-based module.

## Related

- → See `go-naming` — package naming rules this layout depends on.
- → See `go-structs-interfaces` — consumer-side interfaces to break cycles.
- → See `go-function-design` — wiring dependencies in `main`.
