package http

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/model"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/service"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/validator"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/mocks"
	"github.com/stretchr/testify/assert"
)

func setupRouter(handler *JobHandler) *gin.Engine {
	validator.RegisterValidators()
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(ErrorHandler())

	r.POST("/jobs", handler.CreateJob)
	r.GET("/jobs/:id", handler.GetJob)
	r.DELETE("/jobs/:id", handler.DeleteJob)
	return r
}

func TestHandleDeleteJob(t *testing.T) {
	tests := []struct {
		name           string
		jobID          string
		setupMock      func(svc *mocks.MockJobService)
		expectedStatus int
	}{
		{
			name:  "job deleted successfully",
			jobID: "10",
			setupMock: func(svc *mocks.MockJobService) {
				svc.EXPECT().
					DeleteJob(10).
					Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:  "job not found",
			jobID: "10",
			setupMock: func(svc *mocks.MockJobService) {
				svc.EXPECT().
					DeleteJob(10).
					Return(service.ErrNotFound(10))
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "invalid job id",
			jobID:          "abc",
			setupMock:      nil, // service should not be called
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:  "internal error",
			jobID: "10",
			setupMock: func(svc *mocks.MockJobService) {
				svc.EXPECT().
					DeleteJob(10).
					Return(errors.New("db down"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := mocks.NewMockJobService(ctrl)

			if tt.setupMock != nil {
				tt.setupMock(svc)
			}

			handler := NewJobHandler(svc)
			router := setupRouter(handler)

			req := httptest.NewRequest(
				http.MethodDelete,
				"/jobs/"+tt.jobID,
				nil,
			)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestCreateJob(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		setupMock      func(svc *mocks.MockJobService)
		expectedStatus int
	}{
		{
			name: "success",
			body: `{"url":"https://example.com","interval":"1h"}`,
			setupMock: func(svc *mocks.MockJobService) {
				svc.EXPECT().
					CreateJob(gomock.Any()).
					Return(&model.JobResponse{ID: 1}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid url",
			body:           `{"url":"not-a-url","interval":"1h"}`,
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid interval",
			body:           `{"url":"https://example.com","interval":"1tsdf"}`,
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing required url",
			body:           `{"interval":"1tsdf"}`,
			setupMock:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service error",
			body: `{"url":"https://example.com","interval":"1h"}`,
			setupMock: func(svc *mocks.MockJobService) {
				svc.EXPECT().
					CreateJob(gomock.Any()).
					Return(nil, errors.New("db down"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := mocks.NewMockJobService(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(svc)
			}

			handler := NewJobHandler(svc)
			router := setupRouter(handler)

			req := httptest.NewRequest(
				http.MethodPost,
				"/jobs",
				bytes.NewBufferString(tt.body),
			)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
