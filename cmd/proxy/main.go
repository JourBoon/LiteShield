package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"cachedproxy/internal/config"
	"cachedproxy/internal/proxy"
	"cachedproxy/internal/utils"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	logger := utils.NewLogger()

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	reg := proxy.NewRegistry()
	for _, c := range cfg.Clients {
		reg.Add(&proxy.Client{
			ID:         c.ID,
			BackendURL: c.BackendURL,
			TTL:        c.TTL(),
			RateRPS:    c.RateRPS,
		})
		logger.Infof("Loaded client: %s (%s)", c.ID, c.BackendURL)
	}

	handler := proxy.NewProxyHandler(logger, reg)

	// ADMIN KEYS TO DEFINE IN ENVIRONMENT
	adminKey := os.Getenv("ADMIN_KEY")
	if adminKey == "" {
		os.Setenv("ADMIN_KEY", "dev-secret")
	}

	mux := http.NewServeMux()
	mux.Handle("/", handler)
	mux.Handle("/admin/", handler.(*proxy.ProxyHandler).AdminHandler())
	mux.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:           ":" + cfg.Server.Port,
		Handler:        mux,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	logger.Infof("LiteShield proxy started on :%s", cfg.Server.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
