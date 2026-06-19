---
name: go-security
description: "Secure-coding patterns for Go services — secret handling, cryptographically secure randomness, constant-time comparison, TLS configuration, input validation, and injection prevention. Use when handling secrets/tokens, comparing credentials, generating random values, configuring TLS, or accepting untrusted input. If you see math/rand for tokens, == on secrets, missing tls.MinVersion, or string-built SQL → apply this skill. For error messages that must not leak internals → See `go-error-handling`. Do NOT use for auth business logic design."
origin: god_code
---

# Go Security

Part of the **God-Tier Go** set. Security bugs in Go are rarely exotic — they're
`math/rand` where you needed `crypto/rand`, `==` where you needed constant time,
a missing TLS floor. This is a thin router; cited depth lives in `references/`.

## When to Activate

- Handling secrets, API keys, tokens, or passwords.
- Comparing credentials or MACs; generating random IDs/tokens.
- Configuring TLS for a client or server.
- Accepting untrusted input (HTTP bodies, query params, file paths, SQL args).

## Decision Guide

| Situation | Do this | Depth |
|-----------|---------|-------|
| Need a random token/ID/salt | `crypto/rand`, never `math/rand` | `references/secrets-and-crypto.md` |
| Compare a secret/token/MAC | `subtle.ConstantTimeCompare`, never `==` | `references/secrets-and-crypto.md` |
| A secret value in code | Load from env/secret manager; fail fast if absent | `references/secrets-and-crypto.md` |
| Configuring TLS | Set `MinVersion: tls.VersionTLS12`+ | `references/transport-and-input.md` |
| SQL with user input | Parameterized queries, never string concat | `references/transport-and-input.md` |
| File path from user | `filepath.Clean` + validate within a root | `references/transport-and-input.md` |

## Core Rules

1. **Secure randomness:** `crypto/rand` for anything security-relevant; `math/rand` is predictable.
2. **Constant-time compare:** `subtle.ConstantTimeCompare` for secrets — `==` and string compare leak timing.
3. **No hardcoded secrets:** read from environment/secret manager; validate presence at startup.
4. **TLS floor:** always set `MinVersion` (TLS 1.2 minimum, 1.3 preferred); never disable verification in prod.
5. **Validate at the boundary:** never trust external data; validate type, length, and range before use.
6. **Parameterize queries:** pass values as query args; never build SQL/shell strings from input.
7. **Don't leak internals in errors:** see `go-error-handling`; return generic messages to clients, log detail server-side.

## Anti-Patterns

| Smell | Fix |
|-------|-----|
| `math/rand` for a token/session ID | `crypto/rand` (Rule 1) |
| `if token == expected` | `subtle.ConstantTimeCompare` (Rule 2) |
| `apiKey := "sk-..."` in source | Env/secret manager + startup check (Rule 3) |
| `&tls.Config{}` with no `MinVersion` | Set `MinVersion: tls.VersionTLS12` (Rule 4) |
| `InsecureSkipVerify: true` in prod | Verify certs; pin a CA if needed (Rule 4) |
| `"SELECT ... WHERE id='" + id + "'"` | Parameterized query (Rule 6) |
| `os.Open(filepath.Join(root, userPath))` unchecked | Clean + confirm inside root (Rule 5) |

## Checklist

- [ ] All security-relevant randomness uses `crypto/rand`.
- [ ] Secret/token/MAC comparisons use `subtle.ConstantTimeCompare`.
- [ ] No secrets in source; required secrets validated at startup.
- [ ] Every `tls.Config` sets `MinVersion` ≥ TLS 1.2; no `InsecureSkipVerify` in prod.
- [ ] All external input is validated for type/length/range at the boundary.
- [ ] All SQL uses parameterized queries; no string-built queries or shell.
- [ ] Client-facing errors are generic; details are logged server-side only.
- [ ] `gosec ./...` and `govulncheck ./...` run clean (see `assets/`).

## Deep Dives

- `references/secrets-and-crypto.md` — randomness, constant-time compare, secret loading (cited: Vault).
- `references/transport-and-input.md` — TLS config, input validation, injection (cited: Vault).

## Assets

- `assets/secure-http-server.go` — a drop-in hardened `http.Server` + `tls.Config`.
- `assets/gosec-govulncheck.sh` — copyable scan step for CI/local.

## Related

- → See `go-error-handling` — not leaking internals in error messages.
- → See `go-context` — request deadlines that bound resource exhaustion.
- → See `go-observability` — auditing/alerting on security-relevant events.
