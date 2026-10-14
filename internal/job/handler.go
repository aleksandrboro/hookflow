package job

import (
	"net/http"
)

type JobHandler struct {
	// jobService JobService
}

// type JobService interface {
// 	CreateJob() (int, error)
// 	GetJob(int) (models.Job, error)
// 	GetAttempts(int) (int, error)
// 	CancelJob(int) error
// }

func NewJobHandler() *JobHandler {
	return &JobHandler{}
}

func (h *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {}

func (h *JobHandler) GetJob(w http.ResponseWriter, r *http.Request) {}

func (h *JobHandler) GetAttempts(w http.ResponseWriter, r *http.Request) {}

func (h *JobHandler) CancelJob(w http.ResponseWriter, r *http.Request) {}
