package storage

import (
	"context"
	"errors"
	"net"
	"net/url"
	"time"

	"github.com/aleksandrboro/hookflow/internal/config"
	"github.com/aleksandrboro/hookflow/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrEmptyJob            = errors.New("provide job")
	ErrJobNotExists        = errors.New("job not exists")
	ErrEmptyDBConfig       = errors.New("db config cannot be empty")
	ErrUncancellableStatus = errors.New("job with this status cannot be canceled")
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
	if job == nil {
		return "", ErrEmptyJob
	}

	var uuid pgtype.UUID

	if err := uuid.Scan(job.ID); err != nil {
		return "", err
	}

	if _, err := s.db.Exec(ctx,
		"INSERT INTO jobs (id, target_url, payload, run_at, max_attempts, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		uuid, job.TargetURL, job.Payload, job.RunAt, job.MaxAttempts, time.Now(), time.Now()); err != nil {
		return "", err
	}

	return job.ID, nil

}

func (s *Storage) GetJob(ctx context.Context, jobID string) (domain.Job, error) {
	var job domain.Job

	if err := s.db.QueryRow(ctx,
		"SELECT id, target_url, payload, status, run_at, attempts_count, max_attempts, created_at, updated_at FROM jobs WHERE id = $1",
		jobID).Scan(&job.ID, &job.TargetURL, &job.Payload, &job.Status, &job.RunAt, &job.AttemptsCount, &job.MaxAttempts, &job.CreatedAt, &job.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Job{}, ErrJobNotExists
		}

		return domain.Job{}, err
	}

	return job, nil
}

func (s *Storage) CancelJob(ctx context.Context, jobID string) error {
	var job domain.Job

	if _, ok := domain.AllowedToCancel[job.Status]; !ok {
		return ErrUncancellableStatus
	}

	if _, err := s.db.Exec(ctx, "UPDATE jobs SET status = $1 updated_at = $2 WHERE id = $3", domain.JobStatusCancelled, time.Now(), jobID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrJobNotExists
		}

		return err
	}

	return nil
}

func (s *Storage) Close() {
	s.db.Close()
}
