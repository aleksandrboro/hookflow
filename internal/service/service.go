package service

import "github.com/aleksandrboro/hookflow/internal/models"

type JobService struct {
	store Store
}

type Store interface {
	SaveJob(models.Job) (int, error)
	GetJob(int) (*models.Job, error)
	GetAttempts(int) ([]models.Attempt, error)
	CancelJob(int) error
}

func NewJobService(store Store) *JobService {
	return &JobService{
		store: store,
	}
}

func (s *JobService) CreateJob(url string) (int, error) {
	return 0, nil
}

func (s *JobService) GetJob(jobID int) (*models.Job, error) {
	return nil, nil
}

func (s *JobService) GetAttempts(jobID int) ([]models.Attempt, error) {
	return nil, nil
}

func (s *JobService) CancelJob(jobID int) error {
	return nil
}
