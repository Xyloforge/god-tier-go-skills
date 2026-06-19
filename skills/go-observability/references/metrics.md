# Prometheus Metrics in Go

Prometheus instruments itself with the `client_golang` library. Its notifier is a
clean, citable template for how to define metrics correctly.

## Counters — cumulative, suffixed `_total`

A counter only goes up (resets to 0 on restart). Use it for events: requests,
errors, bytes processed. The name **ends in `_total`**, and it always has `Help`.

```go
// prometheus-main/notifier/metric.go:57 — CounterVec with namespace/subsystem,
// a _total suffix, Help text, and a single bounded label.
errors: prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: namespace,
	Subsystem: subsystem,
	Name:      "errors_total",
	Help:      "Total number of sent alerts affected by errors.",
}, []string{alertmanagerLabel}),

sent: prometheus.NewCounterVec(prometheus.CounterOpts{
	Namespace: namespace, Subsystem: subsystem,
	Name: "sent_total", Help: "Total number of alerts sent.",
}, []string{alertmanagerLabel}),
```

## Histograms — distributions, base-unit suffix

Latencies and sizes are *distributions*, so they're histograms, never gauges.
The name carries the **base unit** (`_seconds`, `_bytes`), and buckets are chosen
to bracket the expected range:

```go
// prometheus-main/notifier/metric.go:44 — latency as a histogram in seconds,
// with explicit buckets (and native-histogram config for newer Prometheus).
latencyHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Namespace: namespace, Subsystem: subsystem,
	Name: "latency_histogram_seconds",
	Help: "Latency histogram for sending alert notifications.",
	Buckets: []float64{.01, .1, 1, 10},
}, []string{alertmanagerLabel}),
```

A histogram lets you compute p50/p95/p99 with `histogram_quantile` in PromQL — a
gauge of "current latency" cannot.

## Gauges — instantaneous up/down values

```go
// prometheus-main/notifier/metric.go:79 — queue depth: goes up and down.
queueLength: prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: namespace, Subsystem: subsystem,
	Name: "queue_length",
	Help: "The number of alert notifications in the queue.",
}, []string{alertmanagerLabel}),
```

## Labels — bounded dimensions only

Every metric above is labeled by `alertmanagerLabel` (which Alertmanager) — a
**small, bounded** set. The cardinal rule of labels:

- **Good labels:** `method`, `status_code`, `path_template`, `queue`,
  `instance` — values come from a known, small set.
- **Catastrophic labels:** `user_id`, `request_id`, `email`, raw `url`,
  `timestamp`. Each distinct value creates a new time series; unbounded labels
  blow up memory and the TSDB. Put high-cardinality identifiers in **logs and
  traces**, never in metric labels.

## The RED method

For any request-handling service, instrument three things per route:

- **R**ate — `requests_total` (counter)
- **E**rrors — `requests_total{status="5xx"}` or a separate `errors_total`
- **D**uration — `request_duration_seconds` (histogram)

That trio answers "is it up, is it failing, is it slow?" — see
`assets/recording-and-alert-rules.yml` for ready alert rules over them.

## Registration

Define metrics once (package-level or in a constructor), register them with the
registry, and reuse the handles — never create a metric per request.
`promauto` registers on creation; otherwise call `reg.MustRegister(...)`.
