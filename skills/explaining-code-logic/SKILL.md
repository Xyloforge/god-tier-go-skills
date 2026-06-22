---
name: explaining-code-logic
description: Use when explaining how code, an algorithm, a function, or stateful business logic (reward/eligibility engines, state machines, pricing, cycles, period accruals, anti-cheat checks) works to a human — especially complex or unfamiliar logic where a wall of text would lose the reader and a visual, input-driven, step-by-step walkthrough would land better. Covers onboarding new devs to gnarly code.
---

# Explaining Code Logic

## Overview

Humans don't natively process abstract logic — they process visuals, concrete
examples, and stories. The most effective way to explain complex code is NOT a
wall of prose. It is a single self-contained HTML page that applies four
cognitive-science principles so the reader actually builds a mental model.

**Core principle:** Dual-code (visual + verbal), chunk to fit working memory,
anchor with a familiar analogy, and reveal detail progressively.

## Pick the template first

| If the logic is… | Use | Why |
|---|---|---|
| Linear / conceptual — an algorithm, a flow you trace once, no input-dependent state (e.g. how debounce works, how a parser tokenizes) | `template.html` (static walkthrough) | Layered prose + diagram is enough; nothing to simulate |
| **Stateful business logic whose output depends on inputs** — reward/eligibility engines, state machines, pricing, cycles, period accruals, anti-cheat checks (e.g. the campaign services) | `simulator-template.html` (interactive step-tracer) | The reader must *change inputs and watch state evolve* to actually get it |

When in doubt for real product code, prefer the **simulator** — business logic is
almost always input-driven, and stepping through a concrete scenario is what makes
it click for a new dev.

## When to Use

- User asks "explain how X works", "walk me through this function/algorithm", "help me understand this code".
- Onboarding someone to an unfamiliar codebase or a gnarly piece of logic.
- The logic has multiple steps, branches, or a loop that's easy to lose track of.

**When NOT to use:** a one-line answer suffices, or the user explicitly wants
plain text / a quick verbal answer. Don't generate HTML for trivial questions.

## The Four Principles (apply ALL)

| Principle | Research | How it shows up in the HTML |
|---|---|---|
| Cognitive Load Theory (Sweller) | Working memory holds ~4 chunks | Break logic into discrete numbered steps; strip extraneous detail |
| Dual Coding (Paivio) | Visual + verbal = two memory traces | Diagram beside text; **color-link** each diagram node to its explanation |
| Analogy (Hofstadter) | Learning = mapping new onto familiar | One concrete real-world analogy that mirrors the actual mechanism |
| Progressive Disclosure (UX / Feynman) | Reveal detail on demand | Plain↔Technical toggle, collapsible steps, jargon tooltips, layered depth |

## Workflow A — static walkthrough (`template.html`)

1. **Understand the logic fully first.** Identify four things:
   - The **what** — one sentence, plain language.
   - The **why** — the motivation / problem it solves.
   - The **steps** — the chunks (aim for 3–7).
   - The **hard part** — where humans actually get lost.
2. **Pick ONE analogy** that maps the core mechanism onto something familiar. It
   must mirror the real mechanism, not just be vaguely thematic.
3. **Copy `template.html`** and replace its content with the target logic.
   Preserve its structure and interactions — only swap the content.
4. **Deliver the file** (write it to the working dir / hand the path to the
   user). Then offer to go deeper or shallower. Ask the reader to predict or
   explain a step back — never "does that make sense?"

## Workflow B — interactive simulator (`simulator-template.html`)

For input-driven business logic. The reader feeds a scenario, a consistency
check validates the config, then they step through the real pipeline watching
each input flow into each tracked variable. Port the target logic like this:

1. **Read the target function** and extract its **four parts**:
   - **Inputs** — every argument that changes the outcome → become editable controls.
   - **Config consistency rules** — what makes inputs valid together (e.g. "today must fall inside the date range", "cycle length ≥ 1") → become `checkConsistency()`.
   - **Pipeline stages** — the 3–6 phases the data passes through (e.g. *fetch → chunk into cycles → check conditions → decide state → done*) → become `STAGES[]`.
   - **Tracked state** — the variables/flags that matter and change (e.g. `cycleIndex`, `fulfilled`, `eligible`, `isExpired`) → become keys on `STATE`.
2. **Re-implement the logic in JS inside `runEngine()`**, keeping the real
   structure and the **real variable/flag names from the source** (new devs must
   map your trace back to the actual code). The ONLY addition vs the real code:
   keep tracked vars on `STATE` and call `snap(stage, plainText, codeLine)` at
   every meaningful operation — `snap` records a replayable snapshot.
3. **Write the consistency rules** in `checkConsistency()` — return the list of
   problems so the banner can show exactly what's inconsistent.
4. **Keep PART B (the replay UI) untouched.** Only adjust `STAGES`, `NUM_DAYS`,
   and the input controls' markup to match your logic.
5. **Seed a realistic default scenario** so the page is instantly useful, then
   deliver the file and invite the reader to break it (set a day below the
   minimum, move "today", mark a cycle claimed) and predict the result.

> **Audience default = new devs.** Keep the real TypeScript, real function and
> flag names, and the anti-cheat / verification path visible (behind the
> Technical layer). The goal is they can jump from the trace straight into the
> source file and recognise everything.

## Required HTML components

`template.html` already contains all of these — keep every one:

- One-line "what it does" + a "why it exists" callout at the top.
- **Plain-English ↔ Technical toggle** (progressive disclosure).
- A **diagram** (flow/steps) whose parts are color-linked to the matching text.
- **Collapsible, numbered steps** — one chunk each, default to high-level.
- An **analogy card**.
- **Jargon tooltips** (hover to define).
- A **"predict the next step"** check-for-understanding prompt (active learning).

Self-contained: inline CSS + JS, no external dependencies, opens in any browser.

## Common Mistakes

| Mistake | Fix |
|---|---|
| Wall of text with a diagram bolted on | Interleave — align each visual with the words that describe it |
| Explaining everything at once | Default to the high-level layer; hide depth behind the toggle/accordion |
| "Does that make sense?" | Ask the reader to predict or explain a step back (built into template) |
| Analogy that's thematic but doesn't map | The analogy must mirror the actual mechanism step-for-step |
| Jargon left undefined | Wrap unfamiliar terms in tooltips |

## Reference

- `template.html` — static walkthrough. Complete runnable example (binary
  search) demonstrating every component. Use for linear/conceptual logic.
- `simulator-template.html` — interactive step-tracer. Complete runnable example
  (a simplified repeat-recharge reward engine) with editable inputs, a
  consistency check, a live state tracker, and step/run controls. Use for
  input-driven business logic. Replace PART A (the instrumented port); keep
  PART B (the replay UI).

Adapt the content; preserve each template's structure and interactions.
