// Command hello is the smallest lab app: it greets the user the proxy signed in and counts the visits.
//
// Runtime contract: HTTP on 8080, GET /healthz without auth, the user from the Remote-User header set by the
// proxy, data under /data (empty on first start), secrets as files under /secrets, logs on stdout, a clean
// stop on SIGTERM.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// version is set at build time (-ldflags "-X main.version=...").
var version = "dev"

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	dataDir := env("DATA_DIR", "/data")
	addr := ":" + env("PORT", "8080")

	srv := &http.Server{
		Addr:              addr,
		Handler:           newHandler(dataDir, log),
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer stop()

	go func() {
		log.Info("listening", "addr", addr, "version", version, "data", dataDir)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server failed", "err", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	log.Info("stopping")
	shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdown); err != nil {
		log.Error("shutdown", "err", err)
		os.Exit(1)
	}
	log.Info("stopped")
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
