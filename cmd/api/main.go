// Command api is the composition root for the appointment booking service.
// Phase 0 scope: boot a Gin server with a health check route, wired to config
// and a structured logger, with graceful shutdown. No business logic, no DB
// connection, no routes beyond /health yet — those arrive in later phases.
package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/yourorg/appointment-booking/internal/config"
	"github.com/yourorg/appointment-booking/internal/logger"
)

func main() {
	// .env is only loaded in local dev. In Docker/production, env vars are
	// injected by docker-compose / the orchestrator, and this call is a no-op.
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		// Logger isn't built yet (it depends on cfg.LogLevel), so this is the
		// one place in the whole codebase allowed to print directly — startup
		// failed before we had any logger to use.
		println("FATAL: failed to load config:", err.Error())
		os.Exit(1)
	}

	log := logger.New(cfg.LogLevel, cfg.Environment)
	log.Info().Str("environment", cfg.Environment).Msg("starting appointment-booking API")

	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Recovery())

	registerHealthRoutes(router)

	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Run the server in a goroutine so main can block on the shutdown signal.
	go func() {
		log.Info().Str("port", cfg.ServerPort).Msg("server listening")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server failed to start")
		}
	}()

	// Wait for SIGINT/SIGTERM (Ctrl+C locally, or a Docker/Kubernetes stop
	// signal in production) and shut down cleanly instead of dropping
	// in-flight requests.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutdown signal received, draining connections")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("graceful shutdown failed")
		os.Exit(1)
	}

	log.Info().Msg("server exited cleanly")
}

// registerHealthRoutes wires up the one route Phase 0 cares about: a liveness
// check that confirms the process is up and responding. This intentionally
// does NOT check the database yet — that arrives once Phase 1 wires up
// Postgres, at which point /health becomes a real readiness probe.
func registerHealthRoutes(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"status": "ok",
			},
		})
	})
}