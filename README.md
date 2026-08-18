# chi — Middleware Chain Execution Fix

This repository demonstrates a common Go `net/http` pitfall: middleware continuing execution after writing a response.

## The Problem

Calling `http.Error(w, ...)` or `w.WriteHeader(...)` inside middleware does **not** automatically halt the handler chain. Without an explicit `return`, the next handler executes — causing `http: superfluous response.WriteHeader call` warnings.

## The Fix

Always add `return` after writing a response:

```go
if token == "" {
    http.Error(w, "Unauthorized", http.StatusUnauthorized)
    return // Prevents the chain from continuing
}
```

## Files

- `main.go` — Contains both the broken and correct patterns for comparison.
