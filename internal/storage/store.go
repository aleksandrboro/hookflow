package storage

import (
	"context"
	"errors"

	"github.com/aleksandrboro/hookflow/internal/domain"
)

type Storage struct {
	jobs     []*domain.Job
	attempts []*domain.Attempt
}

func NewStore() *Storage {
	return &Storage{
		jobs:     make([]*domain.Job, 0, 20),
		attempts: make([]*domain.Attempt, 0, 20),
	}
}

func (s *Storage) SaveJob(ctx context.Context, job *domain.Job) (string, error) {
	if job == nil {
		return "", errors.New("provide job")
	}
	s.jobs = append(s.jobs, job)
	return job.ID, nil
}

func (s *Storage) GetJob(ctx context.Context, jobID string) (*domain.Job, error) {
	for _, job := range s.jobs {
		if job.ID == jobID {
			return job, nil
		}
	}

	return nil, errors.New("job not find")
}

func (s *Storage) GetAttempts(ctx context.Context, jobID string) ([]*domain.Attempt, error) {
	attempts := []*domain.Attempt{}

	for _, attempt := range s.attempts {
		if attempt.JobID == jobID {
			attempts = append(attempts, attempt)
		}
	}

	return attempts, nil
}

func (s *Storage) CancelJob(ctx context.Context, jobID string) error {
	for _, job := range s.jobs {
		if job.ID == jobID {
			job.Status = "cancelled"
			return nil
		}
	}

	return nil
}
