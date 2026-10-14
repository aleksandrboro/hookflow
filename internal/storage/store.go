package storage

import "github.com/aleksandrboro/hookflow/internal/models"

type Storage struct {
	data []models.Job
}

func NewStore() *Storage {
	return &Storage{
		data: make([]models.Job, 20),
	}
}

func (s *Storage) SaveJob(job models.Job) (int, error) {
	return 0, nil
}

func (s *Storage) GetJob(jobID int) (*models.Job, error) {
	return nil, nil
}

func (s *Storage) GetAttempts(jobID int) ([]models.Attempt, error) {
	return nil, nil
}

func (s *Storage) CancelJob(jobID int) error {
	return nil
}
