package storage

import (
	"context"

	"github.com/aleksandrboro/hookflow/internal/domain"
)

type Storage struct {
	data []domain.Job
}

func NewStore() *Storage {
	return &Storage{
		data: make([]domain.Job, 0, 20),
	}
}

func (s *Storage) SaveJob(ctx context.Context, job domain.Job) (string, error) {
	return "", nil
}

func (s *Storage) GetJob(ctx context.Context, jobID string) (*domain.Job, error) {
	return nil, nil
}

func (s *Storage) GetAttempts(ctx context.Context, jobID string) ([]domain.Attempt, error) {
	return nil, nil
}

func (s *Storage) CancelJob(ctx context.Context, jobID string) error {
	return nil
}
