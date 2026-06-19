# God-Tier Go Skill Set — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Write 6 god-tier Go skill files + a README index under `god_code/skills/`, every claim backed by a grep-verified citation from the reference repos.

**Architecture:** Each skill is a self-contained `SKILL.md` following one fixed template. `go-clean-code` is the foundational entry point; the other five are specialists it cross-links to via `[[name]]`. Authority comes from real snippets cited as `repo/path/file.go:line` — no invented code. Work proceeds skill-by-skill, README last, then a final boundary+citation audit.

**Tech Stack:** Markdown skill files (ECC SKILL.md format). Recon via `grep`/`rg` over `god_code/`. No build system; the "test" for each task is a verification pass (citations resolve + structure lints).

## Global Constraints

These apply to **every** task implicitly:

- **Location:** all files under `/Users/rong/Documents/np/god_code/skills/`.
- **Not a git repo:** there is no commit step. Each task ends with a **verification step** instead of a commit.
- **Template — every SKILL.md has exactly these sections, in order:**
  ```markdown
  ---
  name: <go-topic>
  description: <one-line trigger an agent matches on>
  origin: god_code
  ---

  # <Title>

  ## When to Activate
  ## Principles        (each: statement → cited snippet → why god-tier)
  ## Anti-Patterns     (wrong vs right, side by side)
  ## Checklist         (- [ ] pass/fail self-checks)
  ## Related           ([[other-skill]] links)
  ```
- **Citation format:** every Go snippet is immediately preceded by a comment
  `// <repo>/<path>.go:<line> — <what & why>`. The path must be real.
- **Authority rule:** before a snippet ships, grep-confirm it exists. If a
  claimed pattern is NOT found in `god_code/`, either remove it or mark it
  `// general Go idiom (not god_code-cited)`. No fabricated file paths.
- **Boundaries (no overlap):** `go-clean-code` = the *prose inside* a function;
  `go-function-design` = the *shape* of the function. Each specialist owns its
  domain; shared concerns stated once and cross-linked.
- **Size:** target 150–280 lines per skill. Many focused principles beat one
  giant essay.

## File Structure

```
god_code/skills/
├── README.md                    # Task 7 — "God-Tier Go" index + install note
├── go-clean-code/SKILL.md       # Task 1 — foundational entry point
├── go-function-design/SKILL.md  # Task 2
├── go-error-handling/SKILL.md   # Task 3
├── go-concurrency/SKILL.md      # Task 4
├── go-performance/SKILL.md      # Task 5
└── go-testing/SKILL.md          # Task 6
```

## Per-task verification procedure (referenced by every task as "VERIFY")

Run these from `/Users/rong/Documents/np/god_code`:

```bash
F=skills/<name>/SKILL.md
# 1. Structure: all required sections present
for s in "## When to Activate" "## Principles" "## Anti-Patterns" "## Checklist" "## Related"; do
  grep -qF "$s" "$F" && echo "OK  $s" || echo "MISSING  $s"
done
# 2. Frontmatter present
head -4 "$F" | grep -q 'origin: god_code' && echo "OK frontmatter" || echo "MISSING frontmatter"
# 3. Every cited path resolves to a real file (extract repo/path.go from citation comments)
grep -oE '(chi-master|gostd|kubernetes|moby|prometheus-main|vault-main)/[A-Za-z0-9_./-]+\.go' "$F" \
  | sort -u | while read p; do [ -f "$p" ] && echo "OK  $p" || echo "BAD CITATION  $p"; done
```
A task passes only when there are **zero** `MISSING`/`BAD CITATION` lines.

---

### Task 1: `go-clean-code` (foundational entry point)

**Files:**
- Create: `skills/go-clean-code/SKILL.md`

**Interfaces:**
- Produces: the entry-point skill. Its `## Related` section links to all five
  specialists: `[[go-function-design]]`, `[[go-error-handling]]`,
  `[[go-concurrency]]`, `[[go-performance]]`, `[[go-testing]]`.

**Principles to cover (each needs one grep-verified snippet):**
1. Clarity over cleverness — direct control flow.
2. Guard clauses / early returns to kill nesting.
3. Naming: short receiver/local names, descriptive exported names.
4. Comment discipline — comment *why*, exported-identifier doc comments.
5. Single responsibility / function size.
6. Zero-value usefulness & avoiding unnecessary state.

- [ ] **Step 1: Recon — pull real anchors**

```bash
cd /Users/rong/Documents/np/god_code
grep -n 'func (mx \*Mux) routeHTTP' chi-master/mux.go            # early-return style
grep -n 'func ' gostd/io/io.go | head -20                       # naming + small funcs
grep -rn 'Copy copies' gostd/io/io.go                           # doc-comment style
sed -n '1,40p' chi-master/middleware/logger.go                  # composition + clarity
```
Record the exact line numbers returned; use them in citations.

- [ ] **Step 2: Verify each chosen snippet exists**

For every file you intend to cite, run `test -f <path> && echo OK`. Expected: `OK` for each.

- [ ] **Step 3: Write `skills/go-clean-code/SKILL.md`**

Follow the Global-Constraints template. Frontmatter:
```yaml
name: go-clean-code
description: Write Go that reads obviously — clarity over cleverness, early returns, precise naming, minimal state. Load first; links to the specialists.
origin: god_code
```
Body: the 6 principles above, each as *statement → cited snippet (real line) →
why god-tier*. Anti-Patterns: nested-if pyramid vs guard clauses; clever
one-liner vs direct code. Checklist: 6–8 pass/fail items (e.g. "no function
nests >3 levels", "every exported identifier has a doc comment starting with
its name"). Related: link all five specialists + note "load me first."

- [ ] **Step 4: VERIFY** (run the per-task verification procedure with `<name>=go-clean-code`). Expected: zero MISSING / BAD CITATION lines.

---

### Task 2: `go-function-design`

**Files:**
- Create: `skills/go-function-design/SKILL.md`

**Interfaces:**
- Consumes: linked from `go-clean-code`.
- Produces: `[[go-function-design]]`. Related links back to `[[go-clean-code]]`
  and across to `[[go-error-handling]]` (error return position).

**Principles to cover:**
1. Accept interfaces, return structs — small interfaces (`io.Reader`/`Writer`).
2. Signature design: context first, error last, no boolean-trap params.
3. When to extract a helper vs inline; designing *shared/reusable* functions.
4. Functional-options pattern for extensible constructors.
5. Receiver choice: pointer vs value, consistency.
6. Variadic & default-friendly APIs.

- [ ] **Step 1: Recon**

```bash
cd /Users/rong/Documents/np/god_code
grep -n 'type Reader interface' gostd/io/io.go
grep -n 'type Writer interface' gostd/io/io.go
sed -n '1,60p' chi-master/middleware/value.go        # small helper design
grep -rn 'func With' chi-master/middleware/logger.go  # options-style ctor
grep -n 'type Handler interface' chi-master/chi.go
```

- [ ] **Step 2: Verify snippets exist** (`test -f` each path; grep the exact lines).

- [ ] **Step 3: Write `skills/go-function-design/SKILL.md`**

Frontmatter:
```yaml
name: go-function-design
description: Design the shape of a Go function — small interface params, context-first/error-last signatures, when to extract shared helpers, functional options, receiver choice.
origin: god_code
```
Body: 6 principles, each cited (chi + stdlib io anchors). Anti-Patterns:
boolean-trap param vs option/struct; giant 8-arg signature vs options;
concrete-type param vs `io.Reader`. Checklist: e.g. "params are the smallest
interface that works", "no boolean trap", "context.Context is first param".
Related: `[[go-clean-code]]`, `[[go-error-handling]]`.

- [ ] **Step 4: VERIFY** (`<name>=go-function-design`). Expected: zero failures.

---

### Task 3: `go-error-handling`

**Files:**
- Create: `skills/go-error-handling/SKILL.md`

**Interfaces:**
- Produces: `[[go-error-handling]]`. Related links to `[[go-clean-code]]`,
  `[[go-concurrency]]` (errgroup error propagation).

**Principles to cover:**
1. Wrap with context using `%w`; readable `fmt.Errorf("verb noun: %w", err)`.
2. Sentinel errors (`var ErrX = errors.New(...)`) + `errors.Is`.
3. Typed errors + `errors.As` for structured handling.
4. Never swallow; handle once at the right layer.
5. Retries/backoff & graceful degradation.
6. Error message style — lowercase, no trailing punctuation, no "failed to".

- [ ] **Step 1: Recon (all confirmed present)**

```bash
cd /Users/rong/Documents/np/god_code
grep -rn 'fmt.Errorf(".*: %w"' vault-main/sdk | head -3
grep -rn 'errors.New("' vault-main/sdk/database/dbplugin/client.go
sed -n '1,60p' kubernetes/cmd/kubeadm/app/util/errors/errors.go   # errors.As usage
grep -rn 'errors.Is(' kubernetes/staging/src/k8s.io/apimachinery | head -3
```

- [ ] **Step 2: Verify snippets exist** (paths confirmed in recon; re-`test -f`).

- [ ] **Step 3: Write `skills/go-error-handling/SKILL.md`**

Frontmatter:
```yaml
name: go-error-handling
description: God-tier Go error handling — wrap with %w, sentinel vs typed errors, errors.Is/As, never-swallow, retries/backoff, idiomatic message style.
origin: god_code
```
Body: 6 principles, cited from Vault + k8s. Anti-Patterns: `err != nil { return
err }` (context-free) vs wrapped; string-compare on `err.Error()` vs
`errors.Is`; swallowed error vs handled. Checklist: e.g. "every returned error
adds context or is a sentinel", "no `err.Error()` string matching", "messages
are lowercase, no punctuation". Related: `[[go-clean-code]]`, `[[go-concurrency]]`.

- [ ] **Step 4: VERIFY** (`<name>=go-error-handling`). Expected: zero failures.

---

### Task 4: `go-concurrency`

**Files:**
- Create: `skills/go-concurrency/SKILL.md`

**Interfaces:**
- Produces: `[[go-concurrency]]`. Related links to `[[go-error-handling]]`,
  `[[go-performance]]` (pool contention).

**Principles to cover:**
1. Goroutine lifecycle — every goroutine has a known exit; leak prevention.
2. `context.Context` for cancellation/deadlines, propagated not stored.
3. Channels vs mutexes — pick by ownership; "share memory by communicating".
4. `errgroup` for fan-out with first-error + cancellation.
5. `sync.Once` / `sync.WaitGroup` correct use.
6. Race-free design; the `-race` discipline.

- [ ] **Step 1: Recon**

```bash
cd /Users/rong/Documents/np/god_code
grep -n 'errgroup' kubernetes/staging/src/k8s.io/cli-runtime/pkg/resource/visitor.go
grep -rn 'sync.Once' moby/daemon 2>/dev/null | head -3
grep -rn 'context.Context' gostd/net/http 2>/dev/null | head -3
grep -n 'type Pool struct' gostd/sync/pool.go
grep -rn 'sync.WaitGroup' moby/daemon 2>/dev/null | head -2
```

- [ ] **Step 2: Verify snippets exist** (`test -f` each; capture lines).

- [ ] **Step 3: Write `skills/go-concurrency/SKILL.md`**

Frontmatter:
```yaml
name: go-concurrency
description: God-tier Go concurrency — goroutine lifecycle and leak prevention, context cancellation, channels vs mutexes, errgroup fan-out, sync primitives, race-free design.
origin: god_code
```
Body: 6 principles cited from k8s/moby/stdlib. Anti-Patterns: fire-and-forget
goroutine (leak) vs lifecycle-bound; mutex-guarded map where a channel fits;
ignoring ctx cancellation. Checklist: e.g. "every goroutine has a defined
stop", "context passed as first arg, never stored in struct", "tested with
-race". Related: `[[go-error-handling]]`, `[[go-performance]]`.

- [ ] **Step 4: VERIFY** (`<name>=go-concurrency`). Expected: zero failures.

---

### Task 5: `go-performance`

**Files:**
- Create: `skills/go-performance/SKILL.md`

**Interfaces:**
- Produces: `[[go-performance]]`. Related links to `[[go-concurrency]]`,
  `[[go-testing]]` (benchmarks).

**Principles to cover:**
1. Measure first — benchmark before optimizing; optimize the hot path only.
2. `sync.Pool` to reuse allocations on hot paths.
3. Preallocate slices/maps with known capacity (`make([]T, 0, n)`).
4. Escape analysis — keep values on the stack; avoid pointer-to-local leaks.
5. Zero-copy / reuse buffers; avoid `[]byte`↔`string` churn.
6. Avoid premature interface boxing & reflection in hot loops.

- [ ] **Step 1: Recon (Prometheus pool code confirmed present)**

```bash
cd /Users/rong/Documents/np/god_code
sed -n '1,60p' prometheus-main/util/zeropool/pool.go
sed -n '1,50p' prometheus-main/util/pool/pool.go
grep -n 'sync.Pool' prometheus-main/tsdb/head_read.go
grep -rn 'make(\[\].*, 0, ' prometheus-main/tsdb 2>/dev/null | head -3
```

- [ ] **Step 2: Verify snippets exist** (`test -f` each path).

- [ ] **Step 3: Write `skills/go-performance/SKILL.md`**

Frontmatter:
```yaml
name: go-performance
description: God-tier Go performance — benchmark-first, sync.Pool reuse, slice/map preallocation, escape analysis, zero-copy buffers, avoiding hot-path allocations and reflection.
origin: god_code
```
Body: 6 principles cited from Prometheus + stdlib. Anti-Patterns: `append` to
nil slice in a loop vs preallocated; allocating per-request vs pooled; `+`
string concat in loop vs `strings.Builder`. Checklist: e.g. "hot path has a
benchmark", "slices preallocated when length is known", "no allocation in the
inner loop without a measured reason". Related: `[[go-concurrency]]`, `[[go-testing]]`.

- [ ] **Step 4: VERIFY** (`<name>=go-performance`). Expected: zero failures.

---

### Task 6: `go-testing`

**Files:**
- Create: `skills/go-testing/SKILL.md`

**Interfaces:**
- Produces: `[[go-testing]]`. Related links to `[[go-performance]]` (benchmarks),
  `[[go-clean-code]]`.

**Principles to cover (god-tier depth beyond existing golang-testing):**
1. Table-driven tests with `t.Run` subtests.
2. `t.Parallel()` and isolation.
3. Benchmarks (`testing.B`, `b.ReportAllocs`, `b.ResetTimer`).
4. Fuzzing (`func FuzzXxx(f *testing.F)`).
5. Golden files & `testdata/`.
6. Helpers with `t.Helper()`; avoiding over-mocking (real deps where cheap).

- [ ] **Step 1: Recon (confirmed present)**

```bash
cd /Users/rong/Documents/np/god_code
sed -n '1,50p' gostd/json/fuzz_test.go
grep -n 'func Benchmark' prometheus-main/tsdb/querier_bench_test.go | head
grep -rn 'for _, tc := range' kubernetes/staging 2>/dev/null | head -3   # table-driven
grep -rn 't.Helper()' kubernetes/staging 2>/dev/null | head -3
```

- [ ] **Step 2: Verify snippets exist** (`test -f` each path).

- [ ] **Step 3: Write `skills/go-testing/SKILL.md`**

Frontmatter:
```yaml
name: go-testing
description: God-tier Go testing — table-driven subtests, t.Parallel isolation, benchmarks with ReportAllocs, fuzzing, golden files, t.Helper, and right-sized mocking.
origin: god_code
```
Body: 6 principles cited from stdlib + Prometheus + k8s. Anti-Patterns: copy-
pasted near-identical test funcs vs table-driven; asserting on a single case vs
fuzz; mocking everything vs real cheap deps. Checklist: e.g. "new logic has a
table-driven test", "hot path has a benchmark", "test helpers call t.Helper()".
Related: `[[go-performance]]`, `[[go-clean-code]]`.

- [ ] **Step 4: VERIFY** (`<name>=go-testing`). Expected: zero failures.

---

### Task 7: README index

**Files:**
- Create: `skills/README.md`

**Interfaces:**
- Consumes: all six skills must exist before this task.

- [ ] **Step 1: Write `skills/README.md`**

Contents:
- Title **"God-Tier Go"** + one-paragraph philosophy (every claim cite-backed by
  chi / k8s / moby / Prometheus / Vault / stdlib).
- A table: skill name → one-line description → "load when…".
- Explicit **load order**: start with `go-clean-code`, drill into specialists.
- **Install note:** "These live in `god_code/` as the source of truth. To make an
  agent auto-trigger them, symlink each folder into `~/.claude/skills/`:
  `ln -s \"$PWD/skills/go-clean-code\" ~/.claude/skills/go-clean-code` (repeat per
  skill), or wrap the folder as a plugin."

- [ ] **Step 2: Verify all six skill links in the README resolve**

```bash
cd /Users/rong/Documents/np/god_code
for s in go-clean-code go-function-design go-error-handling go-concurrency go-performance go-testing; do
  test -f "skills/$s/SKILL.md" && echo "OK $s" || echo "MISSING $s"
done
grep -c 'go-clean-code\|go-function-design\|go-error-handling\|go-concurrency\|go-performance\|go-testing' skills/README.md
```
Expected: six `OK` lines.

---

### Task 8: Final boundary + citation audit

**Files:**
- Modify: any SKILL.md that fails a check.

- [ ] **Step 1: Run citation audit across all skills**

```bash
cd /Users/rong/Documents/np/god_code
grep -rohE '(chi-master|gostd|kubernetes|moby|prometheus-main|vault-main)/[A-Za-z0-9_./-]+\.go' skills \
  | sort -u | while read p; do [ -f "$p" ] && echo "OK  $p" || echo "BAD  $p"; done
```
Expected: zero `BAD` lines. Fix any bad citation inline, then re-run.

- [ ] **Step 2: Cross-link integrity — every `[[link]]` names a real skill**

```bash
cd /Users/rong/Documents/np/god_code
grep -rohE '\[\[[a-z-]+\]\]' skills | tr -d '[]' | sort -u | while read s; do
  test -d "skills/$s" && echo "OK  $s" || echo "DANGLING  $s"
done
```
Expected: zero `DANGLING` lines.

- [ ] **Step 3: Boundary spot-check (manual)**

Read the `## Principles` headings of `go-clean-code` and `go-function-design`
side by side. Confirm clean-code talks about *prose inside* functions and
function-design talks about *function shape*. If any principle is in the wrong
skill, move it. Repeat the eyeball check for any other pair that felt close
(e.g. performance vs concurrency on pools — pool *reuse* lives in performance,
pool *contention/safety* lives in concurrency).

- [ ] **Step 4: Structure lint across all six skills**

```bash
cd /Users/rong/Documents/np/god_code
for s in go-clean-code go-function-design go-error-handling go-concurrency go-performance go-testing; do
  F="skills/$s/SKILL.md"
  for sec in "## When to Activate" "## Principles" "## Anti-Patterns" "## Checklist" "## Related"; do
    grep -qF "$sec" "$F" || echo "$s MISSING $sec"
  done
done
```
Expected: no output (all sections present in all skills).

---

## Self-Review (completed by author)

**Spec coverage:** All 6 spec skills → Tasks 1–6. README + install note → Task 7.
Authority rule → per-task VERIFY + Task 8 citation audit. Boundaries → Task 8
Step 3. Success criteria 1–5 all map to a task. ✅

**Placeholder scan:** No "TBD/implement later". Each task names exact files,
exact recon commands, exact frontmatter, the specific principle list, concrete
anti-patterns, and concrete checklist examples. Snippet prose is written at
execution from grep-verified anchors (the plan supplies the anchors + commands).

**Type/name consistency:** Skill folder names, `name:` frontmatter, and
`[[links]]` use the identical 6 slugs throughout (`go-clean-code`,
`go-function-design`, `go-error-handling`, `go-concurrency`, `go-performance`,
`go-testing`). Verified consistent across all tasks and the audit. ✅
