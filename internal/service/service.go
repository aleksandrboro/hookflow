package service

import (
	"context"
	"errors"

	"github.com/aleksandrboro/hookflow/internal/domain"
	"github.com/aleksandrboro/hookflow/internal/handler"
)

type JobService struct {
	store Store
}

type Store interface {
	SaveJob(context.Context, domain.Job) (string, error)
	GetJob(context.Context, string) (*domain.Job, error)
	GetAttempts(context.Context, string) ([]domain.Attempt, error)
	CancelJob(context.Context, string) error
}

func NewJobService(store Store) *JobService {
	return &JobService{
		store: store,
	}
}

func (s *JobService) CreateJob(ctx context.Context, req handler.CreateJobRequest) (string, error) {
	return "", errors.New("Unimplemented")
}

func (s *JobService) GetJob(ctx context.Context, jobID string) (*domain.Job, error) {
	return nil, errors.New("Unimplemented")
}

func (s *JobService) GetAttempts(ctx context.Context, jobID string) ([]domain.Attempt, error) {
	return nil, errors.New("Unimplemented")
}

func (s *JobService) CancelJob(ctx context.Context, jobID string) error {
	return errors.New("Unimplemented")
}
