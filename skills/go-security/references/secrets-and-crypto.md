# Secrets, Randomness & Constant-Time Comparison

The three highest-frequency Go security mistakes, and how Vault avoids them.

## 1. Cryptographically secure randomness — `crypto/rand`, never `math/rand`

`math/rand` is a deterministic PRNG: given the seed (often the time), its output
is predictable — fatal for tokens, session IDs, salts, or nonces. Use
`crypto/rand`, which reads from the OS CSPRNG. Vault generates one-time tokens
this way:

```go
// vault-main/sdk/helper/roottoken/otp.go:7,25 — crypto/rand for an OTP buffer.
import "crypto/rand"

readLen, err := rand.Read(buf) // OS CSPRNG, not a seeded PRNG
```

```go
// vault-main/sdk/helper/roottoken/otp.go:36 — or a vetted helper built on it.
otp, err := base62.Random(otpLength)
```

Rule of thumb: if a value protects something, it comes from `crypto/rand`. Reach
for `math/rand` only for non-security jitter, sampling, or test data.

## 2. Constant-time comparison — `subtle.ConstantTimeCompare`, never `==`

Comparing secrets with `==` or `bytes.Equal` short-circuits on the first
differing byte, so response time leaks how many leading bytes matched — a timing
oracle an attacker can use to recover the secret byte-by-byte.
`crypto/subtle.ConstantTimeCompare` always examines every byte. Vault guards its
request tokens with it:

```go
// vault-main/http/logical.go:382 — constant-time token check; returns 1 only on
// a full match, and takes the same time regardless of where bytes differ.
if reqToken == "" || expectedToken == "" ||
	subtle.ConstantTimeCompare([]byte(reqToken), []byte(expectedToken)) != 1 {
	// reject
}
```

Notes:
- It returns `1` for equal, `0` otherwise — compare against `1`.
- Inputs of different lengths return `0`; if length itself is sensitive, hash
  both sides to a fixed width first.
- Use it for tokens, passwords, HMAC/MAC tags, and signatures.

## 3. No hardcoded secrets — load and validate at startup

Secrets never belong in source (they leak via VCS history, logs, and binaries).
Read them from the environment or a secret manager, and **fail fast** if a
required one is missing, so a misconfigured deploy crashes loudly instead of
running insecurely:

```go
// Validate required secrets at startup — crash early, not mid-request.
apiKey := os.Getenv("API_KEY")
if apiKey == "" {
	log.Fatal("API_KEY not configured")
}
```

- Keep secrets out of error messages and structured logs (see
  `go-error-handling` and `go-observability`).
- Rotate anything that may have been exposed; treat a leaked secret as burned.
