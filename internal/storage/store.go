package storage

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/aleksandrboro/hookflow/internal/domain"
)

var (
	ErrJobNotExists      = errors.New("job not exists")
	ErrFailedCanceletion = errors.New("job with this status cannot be canceled")
)

type Storage struct {
	jobs          []*domain.Job
	attempts      []*domain.Attempt
	jobsMutex     sync.RWMutex
	attemptsMutex sync.RWMutex
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
	s.jobsMutex.Lock()
	defer s.jobsMutex.Unlock()
	s.jobs = append(s.jobs, job)
	return job.ID, nil
}

func (s *Storage) GetJob(ctx context.Context, jobID string) (domain.Job, error) {
	s.jobsMutex.RLock()
	defer s.jobsMutex.RUnlock()
	for _, job := range s.jobs {
		if job.ID == jobID {
			return *job, nil
		}
	}

	return domain.Job{}, ErrJobNotExists
}

func (s *Storage) GetAttempts(ctx context.Context, jobID string) ([]domain.Attempt, error) {
	attempts := []domain.Attempt{}

	job, err := s.GetJob(ctx, jobID)
	if err != nil {
		return attempts, err
	}

	if job.ID == "" {
		return attempts, ErrJobNotExists
	}

	s.attemptsMutex.Lock()
	defer s.attemptsMutex.Unlock()
	for _, attempt := range s.attempts {
		if attempt.JobID == jobID {
			attempts = append(attempts, *attempt)
		}
	}

	return attempts, nil
}

func (s *Storage) CancelJob(ctx context.Context, jobID string) error {
	s.jobsMutex.Lock()
	defer s.jobsMutex.Unlock()
	for _, job := range s.jobs {
		if job.ID == jobID {
			if job.Status != domain.JobStatusPending && job.Status != domain.JobStatusRetrying && job.Status != domain.JobStatusQueued {
				return fmt.Errorf(ErrFailedCanceletion.Error()+"status: %s", job.Status)
			}
			job.Status = "cancelled"
			return nil
		}
	}

	return ErrJobNotExists
}
