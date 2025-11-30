package http

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter wires up all routes and returns a *gin.Engine
func SetupRouter(jobHandler *JobHandler) *gin.Engine {
	r := gin.Default() // includes Logger + Recovery middleware
	r.Use(ErrorHandler())

	api := r.Group("/api/v1")
	{
		// Job routes

		api.GET("/jobs/:id", jobHandler.GetJob)
		api.GET("/jobs", jobHandler.ListJobs)

		api.POST("/jobs", jobHandler.CreateJob)

		api.POST("/jobs/:id/pause", jobHandler.PauseJob)
		api.POST("/jobs/:id/resume", jobHandler.ResumeJob)

		api.DELETE("/jobs/:id", jobHandler.DeleteJob)
	}

	// You can also add middleware here
	//r.Use(ErrorHandler())

	return r
}
