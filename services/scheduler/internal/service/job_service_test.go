package service

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/model"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/repository"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/mocks"
)

func TestJobService_GetJobByID(t *testing.T) {
	tests := []struct {
		name        string
		jobID       int
		mock        func(repo *mocks.MockJobRepository)
		wantErr     bool
		wantErrCode string
		assert      func(t *testing.T, job *model.JobResponse)
	}{
		{
			name:  "ok",
			jobID: 1,
			mock: func(repo *mocks.MockJobRepository) {
				repo.EXPECT().
					GetByID(1).
					Return(&model.Job{
						ID:       1,
						URL:      "http://example.com",
						Interval: 10 * time.Minute,
					}, nil)
			},
			wantErr: false,
			assert: func(t *testing.T, job *model.JobResponse) {
				require.Equal(t, uint(1), job.ID)
				require.Equal(t, "http://example.com", job.URL)
				require.Equal(t, "10m0s", job.Interval)
			},
		},
		{
			name:  "not found",
			jobID: 42,
			mock: func(repo *mocks.MockJobRepository) {
				repo.EXPECT().
					GetByID(42).
					Return(nil, repository.ErrNotFound)
			},
			wantErr:     true,
			wantErrCode: "NOT_FOUND",
		},
		{
			name:  "repository error",
			jobID: 99,
			mock: func(repo *mocks.MockJobRepository) {
				repo.EXPECT().
					GetByID(99).
					Return(nil, errors.New("db down"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockJobRepository(ctrl)
			tt.mock(repo)

			service := NewJobService(repo)

			job, err := service.GetJobByID(tt.jobID)

			if tt.wantErr {
				require.Error(t, err)

				if tt.wantErrCode != "" {
					appErr, ok := err.(*AppError)
					require.True(t, ok, "expected AppError")
					require.Equal(t, tt.wantErrCode, appErr.Code)
				}

				require.Nil(t, job)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, job)

			if tt.assert != nil {
				tt.assert(t, job)
			}
		})
	}
}

func TestJobService_PauseJob(t *testing.T) {
	tests := []struct {
		name        string
		jobID       int
		mock        func(repo *mocks.MockJobRepository)
		wantErr     bool
		wantErrCode string
		assert      func(t *testing.T, job *model.JobResponse)
	}{
		{
			name:  "pause from scheduled - expect pause job",
			jobID: 1,
			mock: func(repo *mocks.MockJobRepository) {
				repo.EXPECT().
					GetByID(1).
					Return(&model.Job{
						ID:       1,
						URL:      "http://example.com",
						Interval: 10 * time.Minute,
						Status:   model.JobStatusScheduled,
					}, nil)

				repo.EXPECT().
					Save(gomock.Any()).
					Return(nil)
			},
			wantErr: false,
			assert: func(t *testing.T, job *model.JobResponse) {
				require.Equal(t, uint(1), job.ID)
				require.Equal(t, "http://example.com", job.URL)
				require.Equal(t, "10m0s", job.Interval)
				require.Equal(t, false, job.PauseRequested)
				require.Equal(t, model.JobStatusPaused, job.Status)
			},
		},
		{
			name:  "pause from in progress - expect pause requested",
			jobID: 1,
			mock: func(repo *mocks.MockJobRepository) {
				repo.EXPECT().
					GetByID(1).
					Return(&model.Job{
						ID:             1,
						URL:            "http://example.com",
						Interval:       10 * time.Minute,
						Status:         model.JobStatusInProgress,
						PauseRequested: false,
					}, nil)

				repo.EXPECT().
					Save(gomock.Any()).
					Return(nil)
			},
			wantErr: false,
			assert: func(t *testing.T, job *model.JobResponse) {
				require.Equal(t, uint(1), job.ID)
				require.Equal(t, "http://example.com", job.URL)
				require.Equal(t, "10m0s", job.Interval)
				require.Equal(t, true, job.PauseRequested)
				require.Equal(t, model.JobStatusInProgress, job.Status)
			},
		},
		{
			name:  "pause from in failed state - expect ErrCannotPause",
			jobID: 1,
			mock: func(repo *mocks.MockJobRepository) {
				repo.EXPECT().
					GetByID(1).
					Return(&model.Job{
						ID:             1,
						URL:            "http://example.com",
						Interval:       10 * time.Minute,
						Status:         model.JobStatusFailed,
						PauseRequested: false,
					}, nil)

				repo.EXPECT().
					Save(gomock.Any()).
					Times(0)
			},
			wantErr:     true,
			wantErrCode: "CANNOT_PAUSE_JOB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockJobRepository(ctrl)
			tt.mock(repo)

			service := NewJobService(repo)

			job, err := service.PauseJob(tt.jobID)

			if tt.wantErr {
				require.Error(t, err)

				if tt.wantErrCode != "" {
					appErr, ok := err.(*AppError)
					require.True(t, ok, "expected AppError")
					require.Equal(t, tt.wantErrCode, appErr.Code)
				}

				require.Nil(t, job)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, job)

			if tt.assert != nil {
				tt.assert(t, job)
			}
		})
	}
}

func TestDeleteJob(t *testing.T) {
	tests := []struct {
		name            string
		jobID           int
		setupMock       func(m *mocks.MockJobRepository)
		wantErr         bool
		wantedErrStatus int
	}{
		{
			name:  "job does not exist",
			jobID: 5,
			setupMock: func(m *mocks.MockJobRepository) {
				m.EXPECT().
					GetByID(5).
					Return((*model.Job)(nil), repository.ErrNotFound)
			},
			wantErr:         true,
			wantedErrStatus: 404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockJobRepository(ctrl)
			tt.setupMock(repo)

			svc := NewJobService(repo)
			err := svc.DeleteJob(tt.jobID)
			if err == nil {
				assert.False(t, tt.wantErr)
			} else {
				assert.True(t, tt.wantErr)
				appErr, _ := err.(*AppError)
				assert.Equal(t, appErr.Status, tt.wantedErrStatus)
			}
		})
	}
}
