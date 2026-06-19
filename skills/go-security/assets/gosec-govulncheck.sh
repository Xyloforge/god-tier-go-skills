#!/usr/bin/env bash
# Security scan step for Go projects. Run locally or in CI.
# - gosec:       static analysis for insecure code patterns (G-rules).
# - govulncheck: checks your dependencies (and call graph) against the Go
#                vulnerability database; only reports vulns you actually reach.
set -euo pipefail

# Install once (or pin versions in your toolchain / CI image):
#   go install github.com/securego/gosec/v2/cmd/gosec@latest
#   go install golang.org/x/vuln/cmd/govulncheck@latest

echo "==> gosec (static security analysis)"
gosec -severity medium -confidence medium ./...

echo "==> govulncheck (known vulnerabilities in reachable code)"
govulncheck ./...

echo "==> go vet"
go vet ./...

echo "All security checks passed."
