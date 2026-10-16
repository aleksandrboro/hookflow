package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aleksandrboro/hookflow/internal/domain"
)

type JobHandler struct {
	jobService JobService
}

type JobService interface {
	CreateJob(context.Context, CreateJobRequest) (string, error)
	GetJob(context.Context, string) (*domain.Job, error)
	GetAttempts(context.Context, string) ([]domain.Attempt, error)
	CancelJob(context.Context, string) error
}

func NewJobHandler(service JobService) *JobHandler {
	return &JobHandler{
		jobService: service,
	}
}

func (h *JobHandler) HandleCreateJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateJobRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}

	id, err := h.jobService.CreateJob(ctx, req)
	if err != nil {
		http.Error(w, "Failed to create job", http.StatusInternalServerError)
		return
	}

	resp := JobResponse{
		ID:     id,
		Status: domain.JobStatusPending,
		RunAt:  req.RunAt,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to send response", http.StatusInternalServerError)
		return
	}
}

func (h *JobHandler) HandleGetJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	jobID := r.PathValue("job_id")
	if jobID == "" {
		http.Error(w, "Provide job_id", http.StatusBadRequest)
		return
	}

	job, err := h.jobService.GetJob(ctx, jobID)
	if err != nil {
		http.Error(w, "Failed to get job", http.StatusInternalServerError)
		return
	}

	resp := JobResponse{
		ID:            job.ID,
		Status:        job.Status,
		RunAt:         job.RunAt,
		AttemptsCount: job.AttemptsCount,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to send response", http.StatusInternalServerError)
		return
	}
}

func (h *JobHandler) HandleGetAttempts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	jobID := r.PathValue("job_id")
	if jobID == "" {
		http.Error(w, "Provide job_id", http.StatusBadRequest)
		return
	}

	attempts, err := h.jobService.GetAttempts(ctx, jobID)
	if err != nil {
		http.Error(w, "Failed to get attempts", http.StatusInternalServerError)
		return
	}

	resp := AttemptsResponse{
		Data: make([]Attempt, 0, len(attempts)),
	}

	for _, attempt := range attempts {
		resp.Data = append(resp.Data, Attempt{
			ID:               attempt.ID,
			Number:           attempt.Number,
			Status:           attempt.Status,
			StartedAt:        attempt.StartedAt,
			FinishedAt:       attempt.FinishedAt,
			ErrorDescription: attempt.ErrorDescription,
		})
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to send response", http.StatusInternalServerError)
		return
	}
}

func (h *JobHandler) HandleCancelJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	jobID := r.PathValue("job_id")
	if jobID == "" {
		http.Error(w, "Provide job_id", http.StatusBadRequest)
		return
	}

	if err := h.jobService.CancelJob(ctx, jobID); err != nil {
		http.Error(w, "Failed to cancel the job", http.StatusInternalServerError)
		return
	}
}
