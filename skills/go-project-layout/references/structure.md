# Go Repository Layout — learning from the corpus

The four repos in this corpus span the whole range, from a single flat package to
a giant monorepo. Their choices teach the rules.

## chi — a library, flat by design

```
chi-master/
├── go.mod          // module github.com/go-chi/chi/v5
├── chi.go          // package chi — public API at the root
├── mux.go
├── tree.go
├── context.go
└── middleware/     // a focused sub-package
```

Lesson: a library's public API lives in the **root package**, named for the
module (`chi`). No `pkg/`, no `internal/` ceremony — it's small and cohesive.
Note the `/v5` in the module path: that's the required suffix for v2+ modules
(semantic import versioning).

## Prometheus — an application with private internals

```
prometheus-main/
├── go.mod          // module github.com/prometheus/prometheus
├── cmd/            // binaries: cmd/prometheus, cmd/promtool
├── internal/       // private to this module
├── tsdb/           // domain package (the time-series database)
├── notifier/       // domain package
└── util/           // narrowly-scoped utilities (pool, zeropool, ...)
```

Lessons: binaries live in `cmd/<name>`; packages are named for their **domain**
(`tsdb`, `notifier`), not for layers. Even `util/` here is not a junk drawer —
it's split into cohesive sub-packages (`util/pool`, `util/zeropool`), each doing
one thing.

## Moby — a large app leaning hard on `internal/`

```
moby/
├── go.mod          // module github.com/moby/moby/v2
├── cmd/            // entrypoints
├── api/            // the public HTTP API types
├── daemon/         // the engine
│   └── internal/   // private to daemon: restartmanager, stream, ...
├── internal/       // private to the whole module
└── pkg/            // packages intended for external reuse
```

Lessons: `internal/` appears at **multiple levels** — `daemon/internal/` is
importable by `daemon/...` but nothing else. This scopes privacy precisely. The
split between `api/` (public types), `pkg/` (reusable), and `internal/` (private)
is a deliberate statement of what's supported for outside use.

## Kubernetes — monorepo scale

```
kubernetes/
├── go.mod          // module k8s.io/kubernetes
├── cmd/            // every binary (kube-apiserver, kubelet, kubeadm, ...)
├── pkg/            // core implementation packages
├── api/            // API definitions
└── staging/        // packages published as separate modules (k8s.io/client-go, ...)
```

Lesson: even at this scale the same primitives hold — `cmd/` for binaries, domain
packages under `pkg/`. The `staging/` mechanism is how Kubernetes publishes parts
of the monorepo (e.g. `client-go`) as independently importable modules; that's an
advanced pattern you almost certainly don't need.

## The portable rules

1. **`cmd/<name>/main.go`** per binary; `main` only wires and starts.
2. **`internal/`** — use it at whatever level scopes privacy correctly; the
   compiler enforces it.
3. **Domain packages**, named for what they do (`tsdb`, `daemon`, `notifier`).
4. **No `utils`/`common`** — note that even Prometheus's `util` is decomposed into
   single-purpose sub-packages, not a catch-all.
5. **Module path** = the import path; add `/vN` for v2+.
6. **`pkg/`/`api/`** are choices for *large* repos to signal public vs internal —
   a small service doesn't need them.
