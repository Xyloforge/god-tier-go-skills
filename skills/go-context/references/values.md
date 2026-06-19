# Context Values — the typed-key pattern

Request-scoped values (trace IDs, the authenticated user, a routing context) ride
on the `context.Context`. The danger is the *key*: a `string` key can collide
across packages and boxes into `interface{}` with an allocation. The god-tier
pattern — straight from the Go team and copied by chi — is an **unexported
pointer to a struct type**.

## The pattern (chi)

```go
// chi-master/context.go:157 — the canonical key type. Unexported, pointer-sized,
// self-documenting via String(), collision-proof across packages.
//
// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "chi context value " + k.name
}
```

```go
// chi-master/context.go:39 — one exported key value, created once.
RouteCtxKey = &contextKey{"RouteContext"}
```

```go
// chi-master/context.go:28 — read it back with a typed assertion.
val, _ := ctx.Value(RouteCtxKey).(*Context)
```

## Why this is god-tier

- **No collisions.** Two packages can both use a key named `"user"` as a string
  and silently clobber each other. Pointer identity makes every key globally
  unique — even two `contextKey{"user"}` values are different pointers.
- **Zero allocation.** A pointer already fits in an `interface{}` word, so
  `ctx.Value` lookups don't allocate (the comment says exactly this).
- **Self-describing.** `String()` makes the value print readably in logs/panics.

## How to apply it in your package

```go
// 1. Unexported key type — callers can't forge or collide with it.
type ctxKey struct{ name string }

// 2. One private key value per datum.
var userKey = &ctxKey{"user"}

// 3. Typed setter/getter wrap the untyped context API.
func WithUser(ctx context.Context, u *User) context.Context {
	return context.WithValue(ctx, userKey, u)
}

func UserFrom(ctx context.Context) (*User, bool) {
	u, ok := ctx.Value(userKey).(*User)
	return u, ok
}
```

Expose `WithUser`/`UserFrom`, keep `userKey` unexported. Callers never touch the
raw key, so the value's type and presence are controlled in one place.

## What NOT to put in a context value

Context values are for **request-scoped data that crosses API boundaries** —
trace IDs, auth principals, deadlines. They are *not* a side channel for optional
function parameters. If a function needs a value to do its job, that value
belongs in its signature (see `go-function-design`), not smuggled through
`ctx.Value`.

```go
// Bad — optional behavior hidden in the context.
ctx = context.WithValue(ctx, "retries", 3)

// Good — it's a parameter; the signature tells the truth.
func fetch(ctx context.Context, url string, retries int) error
```
