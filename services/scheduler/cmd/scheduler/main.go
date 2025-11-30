package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/config"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/db"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/handler/http"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/model"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/repository/postgres"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/scheduler"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/service"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/validator"
)

type logDispatcher struct{}

func (d *logDispatcher) DispatchJobs(jobs []model.JobDispatched) error {
	for _, job := range jobs {
		fmt.Printf("dispatching job: id=%d, url=%s\n", uint64(job.ID), job.URL)
	}
	return nil
}

var shutDownWg sync.WaitGroup

func main() {
	// Context that cancels on SIGINT or SIGTERM
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.Load("local")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Loaded config: %+v\n", cfg)

	gormDB, err := db.Connect(&cfg.DB)
	if err != nil {
		panic(err)
	}

	// reset db
	if err := db.Reset(gormDB); err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(gormDB); err != nil {
		panic(err)
	}

	repo := postgres.New(gormDB)
	jobSvc := service.NewJobService(repo)
	jobHandler := http.NewJobHandler(jobSvc)

	r := http.SetupRouter(jobHandler)
	validator.RegisterValidators()
	// register application middleware

	scheduler := scheduler.NewScheduler(&cfg.Scheduler, repo, &logDispatcher{})

	StartScheduler(ctx, scheduler)

	// start api server
	StartAPI(ctx, r, cfg.Server.Port)

	// Wait for SIGINT, SIGTERM or cancel signa l
	<-ctx.Done()
	GracefulShutdown(cfg.Server.ShutdownTimeoutSeconds)
}

func StartScheduler(ctx context.Context, scheduler *scheduler.Scheduler) {
	shutDownWg.Add(1)
	go func() {
		defer shutDownWg.Done()
		scheduler.Run(ctx)
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
	log.Printf("[INFO] shutting down with timeout period of %s ...\n", timeout)

	done := make(chan struct{})
	go func() {
		shutDownWg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("[INFO] all goroutines finished, exiting")
	case <-time.After(timeout):
		log.Println("[INFO] scheduler and server stopped, exiting")
	}
}
