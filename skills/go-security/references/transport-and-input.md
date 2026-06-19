# Transport Security & Input Validation

## TLS configuration — always set a version floor

A zero-value `tls.Config` historically allowed obsolete protocol versions. Always
pin a `MinVersion`. Vault sets TLS 1.2 as the floor everywhere it builds a config:

```go
// vault-main/sdk/helper/certutil/types.go:614 — explicit MinVersion floor.
tlsConfig := &tls.Config{
	MinVersion: tls.VersionTLS12,
}
```

Rules:
- **`MinVersion: tls.VersionTLS12`** minimum; prefer `tls.VersionTLS13` for new
  internal services.
- **Never `InsecureSkipVerify: true` in production.** It disables certificate
  verification, defeating TLS entirely. For private CAs, set `RootCAs` to a pool
  containing your CA instead of skipping verification.
- Let Go pick cipher suites for TLS 1.3 (they're not configurable, by design);
  only constrain suites for TLS 1.2 if you have a specific requirement.

A ready-to-use hardened server config is in `assets/secure-http-server.go`.

## Input validation — never trust external data

Validate at the system boundary, before the data reaches any logic. Check type,
length, and range; reject early with a clear (non-leaking) error.

```go
// Validate shape and bounds up front.
func parseLimit(raw string) (int, error) {
	n, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("limit must be an integer")
	}
	if n < 1 || n > 1000 {
		return 0, fmt.Errorf("limit out of range [1,1000]")
	}
	return n, nil
}
```

- Use schema validation for structured payloads; set explicit size limits on
  request bodies (`http.MaxBytesReader`) to bound memory.
- Pair with request deadlines (see `go-context`) so slow/oversized input can't
  exhaust resources.

## Injection prevention

**SQL — parameterize, never concatenate:**

```go
// Bad — string concatenation is an injection vector.
db.Query("SELECT * FROM users WHERE id = '" + id + "'")

// Good — the driver binds the value safely.
db.QueryContext(ctx, "SELECT * FROM users WHERE id = $1", id)
```

**OS commands — pass args as a slice, never a shell string:**

```go
// Good — no shell, arguments are not re-parsed.
exec.CommandContext(ctx, "convert", in, out)
// Bad — invoking a shell with interpolated input enables command injection.
exec.CommandContext(ctx, "sh", "-c", "convert "+in+" "+out)
```

**File paths — clean and confine to a root:**

```go
clean := filepath.Clean(userPath)
full := filepath.Join(root, clean)
if !strings.HasPrefix(full, root+string(os.PathSeparator)) {
	return errors.New("path escapes root")
}
```

This blocks `../../etc/passwd`-style path traversal.
