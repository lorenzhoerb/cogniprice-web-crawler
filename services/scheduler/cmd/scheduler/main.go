package main

import (
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/config"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/db"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/handler/http"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/repository/postgres"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/service"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/validator"
)

func main() {
	db, err := db.Connect(&config.DBConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "1234",
		DBName:   "cp_scheduler",
		SSLMode:  "disable",
	})
	if err != nil {
		panic(err)
	}

	repo := postgres.New(db)
	jobSvc := service.NewJobService(repo)
	jobHandler := http.NewJobHandler(jobSvc)

	r := http.SetupRouter(jobHandler)
	validator.RegisterValidators()
	// register application middleware

	r.Run(":8080")
}
