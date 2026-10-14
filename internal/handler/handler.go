package handler

import (
	"net/http"

	"github.com/aleksandrboro/hookflow/internal/models"
)

type JobHandler struct {
	jobService JobService
}

type JobService interface {
	CreateJob(string) (int, error)
	GetJob(int) (*models.Job, error)
	GetAttempts(int) ([]models.Attempt, error)
	CancelJob(int) error
}

func NewJobHandler(service JobService) *JobHandler {
	return &JobHandler{
		jobService: service,
	}
}

func (h *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {}

func (h *JobHandler) GetJob(w http.ResponseWriter, r *http.Request) {}

func (h *JobHandler) GetAttempts(w http.ResponseWriter, r *http.Request) {}

func (h *JobHandler) CancelJob(w http.ResponseWriter, r *http.Request) {}
