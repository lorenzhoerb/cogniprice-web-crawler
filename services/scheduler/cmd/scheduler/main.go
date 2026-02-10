package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/lorenzhoerb/cogniprice/services/scheduler/app"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/config"
	"github.com/lorenzhoerb/cogniprice/shared/logger"
	"go.uber.org/zap"
)

func main() {
	logger.Init(true)

	// Context that cancels on SIGINT or SIGTERM
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	env := "local"
	cfg, err := config.Load(env)
	if err != nil {
		logger.Log.Fatal("failed to load config", zap.Error(err))
	}

	// initialize the app
	a, err := app.New(cfg)
	if err != nil {
		logger.Log.Fatal("failed to initialize app", zap.Error(err))
	}

	a.Start(ctx)

	// wait for signal
	<-ctx.Done()

	a.Shutdown(cfg.Server.ShutdownTimeoutSeconds)
}
