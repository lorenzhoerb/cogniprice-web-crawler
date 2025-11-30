package main

import (
	"fmt"

	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/config"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/db"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/handler/http"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/repository/postgres"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/service"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/validator"
)

func main() {
	cfg, err := config.Load("local")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Loaded config: %+v\n", cfg)

	db, err := db.Connect(&cfg.DB)
	if err != nil {
		panic(err)
	}

	repo := postgres.New(db)
	jobSvc := service.NewJobService(repo)
	jobHandler := http.NewJobHandler(jobSvc)

	r := http.SetupRouter(jobHandler)
	validator.RegisterValidators()
	// register application middleware

	r.Run(fmt.Sprintf(":%d", cfg.Server.Port))
}
