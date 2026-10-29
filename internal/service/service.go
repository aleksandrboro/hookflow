package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/aleksandrboro/hookflow/internal/domain"
	"github.com/aleksandrboro/hookflow/internal/handler"
	"github.com/google/uuid"
)

type JobService struct {
	store Store
}

type Store interface {
	SaveJob(context.Context, *domain.Job) (string, error)
	GetJob(context.Context, string) (domain.Job, error)
	CancelJob(context.Context, string) error
}

func NewJobService(store Store) *JobService {
	return &JobService{
		store: store,
	}
}

func (s *JobService) CreateJob(ctx context.Context, req handler.CreateJobRequest) (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", errors.New("failed to generate id")
	}

	job := &domain.Job{
		ID:            id.String(),
		TargetURL:     req.URL,
		Payload:       req.Payload,
		Status:        domain.JobStatusPending,
		RunAt:         req.RunAt,
		AttemptsCount: 0,
		MaxAttempts:   req.MaxAttempts,
	}

	jobID, err := s.store.SaveJob(ctx, job)
	if err != nil {
		return "", fmt.Errorf("failed to save job: %w", err)
	}

	return jobID, nil
}

func (s *JobService) GetJob(ctx context.Context, jobID string) (domain.Job, error) {
	job, err := s.store.GetJob(ctx, jobID)
	if err != nil {
		return domain.Job{}, fmt.Errorf("failed to get job: %w", err)
	}

	return job, nil
}

func (s *JobService) CancelJob(ctx context.Context, jobID string) error {
	if err := s.store.CancelJob(ctx, jobID); err != nil {
		return fmt.Errorf("failed to cancel job: %w", err)
	}

	return nil
}
