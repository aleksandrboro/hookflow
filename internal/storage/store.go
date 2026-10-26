package storage

import (
	"context"
	"errors"
	"net"
	"net/url"

	"github.com/aleksandrboro/hookflow/internal/config"
	"github.com/aleksandrboro/hookflow/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrJobNotExists      = errors.New("job not exists")
	ErrFailedCanceletion = errors.New("job with this status cannot be canceled")
	ErrEmptyDBConfig     = errors.New("db config cannot be empty")
)

type Storage struct {
	db *pgxpool.Pool
}

func NewStore(ctx context.Context, dbCfg *config.DBCfg) (*Storage, error) {
	if dbCfg == nil {
		return nil, ErrEmptyDBConfig
	}

	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(dbCfg.User, dbCfg.Password),
		Host:   net.JoinHostPort(dbCfg.Host, dbCfg.Port),
		Path:   "/" + dbCfg.Name,
	}

	pool, err := pgxpool.New(ctx, u.String())
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return &Storage{
		db: pool,
	}, nil
}

func (s *Storage) SaveJob(ctx context.Context, job *domain.Job) (string, error) {
	// if job == nil {
	// 	return "", errors.New("provide job")
	// }
	// s.jobsMutex.Lock()
	// defer s.jobsMutex.Unlock()
	// s.jobs = append(s.jobs, job)
	// return job.ID, nil

	return "", ErrEmptyDBConfig
}

func (s *Storage) GetJob(ctx context.Context, jobID string) (domain.Job, error) {
	// s.jobsMutex.RLock()
	// defer s.jobsMutex.RUnlock()
	// for _, job := range s.jobs {
	// 	if job.ID == jobID {
	// 		return *job, nil
	// 	}
	// }

	// return domain.Job{}, ErrJobNotExists

	return domain.Job{}, ErrEmptyDBConfig
}

func (s *Storage) GetAttempts(ctx context.Context, jobID string) ([]domain.Attempt, error) {
	// attempts := []domain.Attempt{}

	// job, err := s.GetJob(ctx, jobID)
	// if err != nil {
	// 	return attempts, err
	// }

	// if job.ID == "" {
	// 	return attempts, ErrJobNotExists
	// }

	// s.attemptsMutex.Lock()
	// defer s.attemptsMutex.Unlock()
	// for _, attempt := range s.attempts {
	// 	if attempt.JobID == jobID {
	// 		attempts = append(attempts, *attempt)
	// 	}
	// }

	// return attempts, nil
	return nil, ErrEmptyDBConfig
}

func (s *Storage) CancelJob(ctx context.Context, jobID string) error {
	// s.jobsMutex.Lock()
	// defer s.jobsMutex.Unlock()
	// for _, job := range s.jobs {
	// 	if job.ID == jobID {
	// 		if job.Status != domain.JobStatusPending && job.Status != domain.JobStatusRetrying && job.Status != domain.JobStatusQueued {
	// 			return fmt.Errorf("%w status: %s", ErrFailedCanceletion, job.Status)
	// 		}
	// 		job.Status = "cancelled"
	// 		return nil
	// 	}
	// }

	// return ErrJobNotExists
	return ErrEmptyDBConfig
}

func (s *Storage) Close() {
	s.db.Close()
}
