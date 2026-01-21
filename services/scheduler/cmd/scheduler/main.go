package main

import (
	"context"
	"fmt"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/api/http"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/config"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/db"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/dispatcher"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/repository/postgres"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/scheduler"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/service"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/validator"
	"github.com/lorenzhoerb/cogniprice/shared/logger"
	"go.uber.org/zap"
)

var shutDownWg sync.WaitGroup

func main() {
	logger.Init(true)

	// Context that cancels on SIGINT or SIGTERM
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.Load("local")
	if err != nil {
		logger.Log.Fatal("failed to load config", zap.Error(err))
	}

	logger.Log.Info("config loaded", zap.String("env", "local"))

	gormDB, err := db.Connect(&cfg.DB)
	if err != nil {
		logger.Log.Fatal("failed to connect to DB", zap.Error(err))
	}

	// reset db
	if err := db.Reset(gormDB); err != nil {
		logger.Log.Fatal("failed to reset DB", zap.Error(err))
	}
	if err := db.AutoMigrate(gormDB); err != nil {
		logger.Log.Fatal("failed to auto-migrate DB", zap.Error(err))
	}

	repo := postgres.New(gormDB)
	jobSvc := service.NewJobService(repo)
	jobHandler := http.NewJobHandler(jobSvc)

	r := http.SetupRouter(jobHandler)
	validator.RegisterValidators()
	// register application middleware

	scheduler := scheduler.NewScheduler(&cfg.Scheduler, repo, dispatcher.NewLogDispatcher())

	StartScheduler(ctx, scheduler)

	// start api server
	StartAPI(ctx, r, cfg.Server.Port)

	// Wait for SIGINT, SIGTERM or cancel signal
	<-ctx.Done()
	GracefulShutdown(cfg.Server.ShutdownTimeoutSeconds)
}

func StartScheduler(ctx context.Context, scheduler *scheduler.Scheduler) {
	shutDownWg.Add(1)
	go func() {
		defer shutDownWg.Done()
		logger.Log.Info("scheduler started")
		scheduler.Run(ctx)
		logger.Log.Info("scheduler stopped")
	}()
}

func StartAPI(ctx context.Context, ginEngine *gin.Engine, port int) {
	//shutDownWg.Add(1)
	go func() {
		defer shutDownWg.Done()
		ginEngine.Run(fmt.Sprintf(":%d", port))
	}()
}

func GracefulShutdown(timeoutSeconds int) {
	timeout := time.Duration(timeoutSeconds) * time.Second
	logger.Log.Info("shutting down", zap.Duration("timeout", timeout))

	done := make(chan struct{})
	go func() {
		shutDownWg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.Log.Info("all goroutines finished, exiting")
	case <-time.After(timeout):
		logger.Log.Warn("shutdown timeout reached, exiting")

	}
}
