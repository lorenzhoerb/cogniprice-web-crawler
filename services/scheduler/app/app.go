package app

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/api/http"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/config"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/db"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/dispatcher"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/messaging"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/repository/postgres"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/scheduler"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/service"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/validator"
	"github.com/lorenzhoerb/cogniprice/shared/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	// immutable
	cfg *config.Config

	// infrastructure
	db     *gorm.DB
	rabbit *messaging.Rabbit

	// root components
	scheduler *scheduler.Scheduler
	router    *gin.Engine

	// lifecycle
	wg sync.WaitGroup
}

func New(cfg *config.Config) (*App, error) {
	logger.Log.Info("config loaded", zap.String("env", cfg.Env))

	// ---- DB ----
	gormDB, err := db.Connect(&cfg.DB)
	if err != nil {
		return nil, err
	}

	if err := db.Reset(gormDB); err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(gormDB); err != nil {
		return nil, err
	}

	// ---- RabbitMQ ----
	rabbit, err := messaging.NewRabbit(cfg.Rabbit.Url)
	if err != nil {
		return nil, err
	}

	ch, err := rabbit.Conn.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	if err := messaging.DeclareTopology(ch); err != nil {
		return nil, err
	}

	// ---- Wiring ----
	repo := postgres.New(gormDB)
	jobSvc := service.NewJobService(repo)
	jobHandler := http.NewJobHandler(jobSvc)

	dispatcher, err := dispatcher.NewRabbitDispatcher(rabbit)
	if err != nil {
		return nil, err
	}

	router := http.SetupRouter(jobHandler)
	validator.RegisterValidators()

	sched := scheduler.NewScheduler(&cfg.Scheduler, repo, dispatcher)

	return &App{
		cfg:       cfg,
		db:        gormDB,
		rabbit:    rabbit,
		scheduler: sched,
		router:    router,
	}, nil
}

func (a *App) Start(ctx context.Context) {
	a.startScheduler(ctx)
	a.startAPI()
}

func (a *App) startScheduler(ctx context.Context) {
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		logger.Log.Info("scheduler started")
		a.scheduler.Run(ctx)
		logger.Log.Info("scheduler stopped")
	}()
}

func (a *App) startAPI() {
	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		addr := fmt.Sprintf(":%d", a.cfg.Server.Port)
		logger.Log.Info("api started", zap.String("addr", addr))
		_ = a.router.Run(addr)
	}()
}

func (a *App) Shutdown(timeoutSeconds int) {
	timeout := time.Duration(timeoutSeconds) * time.Second
	logger.Log.Info("shutting down", zap.Duration("timeout", timeout))

	done := make(chan struct{})
	go func() {
		a.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.Log.Info("all goroutines finished, exiting")
	case <-time.After(timeout):
		logger.Log.Warn("shutdown timeout reached, exiting")
	}

	if a.rabbit != nil {
		_ = a.rabbit.Close()
	}
}
