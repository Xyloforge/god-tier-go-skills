---
name: go-structs-interfaces
description: "Designing Go types — defining interfaces at the point of use, keeping them small, struct composition via embedding, useful zero values, and field layout. Use when declaring an interface or struct, deciding embedding vs a named field, or where an interface should live. If you see a big interface defined next to its only implementation, an interface accepted then never abstracted, or a type needing a constructor just to be usable → apply this skill. For function signatures → See `go-function-design`. Do NOT use for naming (→ See `go-naming`)."
origin: god_code
---

# Go Structs & Interfaces

Part of the **God-Tier Go** set. This governs **type definitions** — the shape of
interfaces and structs themselves — where `go-function-design` governs the
functions that use them. Thin router; cited depth in `references/`.

## When to Activate

- Declaring an `interface` or a `struct` type.
- Deciding embedding vs a named field; composition vs inheritance instincts.
- Choosing where an interface should be defined (producer vs consumer side).
- Designing a type to be usable at its zero value.

## Decision Guide

| Question | Answer | Depth |
|----------|--------|-------|
| How big should an interface be? | As small as the consumer needs (often 1 method) | `references/interfaces.md` |
| Where do I define the interface? | In the **consumer** package, not the implementer | `references/interfaces.md` |
| Should I return an interface? | No — return the concrete struct | `references/interfaces.md` |
| Embed or named field? | Embed for "is-a behavior" / promotion; field for "has-a" | `references/structs.md` |
| Does my type need a constructor? | Only if the zero value isn't usable | `references/structs.md` |

## Core Rules

1. **Interfaces belong to the consumer.** Define the interface where it's *used*, listing only the methods that caller needs. The implementer just has concrete methods.
2. **Keep interfaces tiny** — one to three methods. `io.Reader` (one method) composes with everything; a ten-method interface composes with nothing.
3. **Accept interfaces, return concrete structs.** Abstraction is for inputs; callers want the full output type. (Mirrors `go-function-design`.)
4. **Compose with embedding.** Embed a type to promote its methods/fields; embed `sync.Mutex` to get `Lock`/`Unlock` on your type. Prefer composition over deep type hierarchies (Go has none).
5. **Make the zero value useful.** Design structs so `var x T` works without init where possible (`sync.Mutex`, `bytes.Buffer` do). A constructor is needed only when the zero value can't be valid.
6. **Group struct fields by meaning**; place embedded types first. Don't micro-optimize field order for padding unless a benchmark demands it.
7. **Satisfy interfaces implicitly** — no `implements` keyword. Optionally assert with `var _ Iface = (*T)(nil)` to catch drift at compile time.

## Anti-Patterns

| Smell | Fix |
|-------|-----|
| Big interface beside its single implementation | Define a small interface at the consumer (Rule 1) |
| 8-method interface "for flexibility" | Split into small, composable ones (Rule 2) |
| Function returns an interface to "be generic" | Return the concrete struct (Rule 3) |
| Simulating inheritance with deep types | Compose via embedding (Rule 4) |
| `NewT()` that only zeroes fields | Make the zero value usable (Rule 5) |
| Interface with `IFoo`/`FooInterface` naming | Name for behavior (→ `go-naming`) |

## Checklist

- [ ] Interfaces are defined in the consumer and list only methods that caller uses.
- [ ] Interfaces are small (1–3 methods); large ones are decomposed.
- [ ] Functions accept interfaces and return concrete structs.
- [ ] Composition uses embedding; no inheritance-shaped type hierarchies.
- [ ] The zero value is usable, or a constructor exists for a real reason.
- [ ] Fields are grouped by meaning; embedded types come first.
- [ ] Compile-time interface assertions guard important implementations.

## Deep Dives

- `references/interfaces.md` — size, consumer-side definition, return-concrete (cited: stdlib `io`, chi).
- `references/structs.md` — embedding, zero values, field layout (cited: Moby, chi).

## Related

- → See `go-function-design` — signatures that consume these types.
- → See `go-naming` — naming interfaces (`-er`) and types.
- → See `go-concurrency` — embedding `sync` types and copy semantics.
