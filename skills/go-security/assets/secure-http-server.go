// Package main shows a hardened net/http server configuration.
//
// Drop-in reference: a production http.Server should always set timeouts (to
// resist slowloris and resource exhaustion) and a TLS floor. Adapt the handler
// and addresses to your service.
package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"
)

func newSecureServer(handler http.Handler, addr string) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: handler,

		// Timeouts bound how long a single connection can tie up resources.
		// Without these, a slow client can hold a connection open indefinitely.
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,

		// Cap header size to bound memory per request.
		MaxHeaderBytes: 1 << 20, // 1 MiB

		// TLS floor: never serve below TLS 1.2; prefer 1.3 for internal services.
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			// For TLS 1.3, Go selects cipher suites automatically (not configurable).
			// Do NOT set InsecureSkipVerify here or on any client used in prod.
		},
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	srv := newSecureServer(mux, ":8443")
	// ListenAndServeTLS enforces the TLSConfig above.
	log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))
}
