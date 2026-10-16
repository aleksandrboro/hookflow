package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/aleksandrboro/hookflow/internal/domain"
	"github.com/gorilla/mux"
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

	w.Header().Set("Content-Type", "application/json")

	if err := decoder.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid json")
		return
	}

	if req.RunAt.IsZero() {
		req.RunAt = time.Now().UTC()
	}

	id, err := h.jobService.CreateJob(ctx, req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create job")
		return
	}

	resp := JobResponse{
		ID:     id,
		Status: domain.JobStatusPending,
		RunAt:  req.RunAt,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to send response")
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *JobHandler) HandleGetJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	jobID, ok := mux.Vars(r)["job_id"]
	if !ok {
		writeError(w, http.StatusBadRequest, "Provide job_id")
		return
	}

	job, err := h.jobService.GetJob(ctx, jobID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get job")
		return
	}

	if job == nil {
		writeError(w, http.StatusNotFound, "Task not fount")
		return
	}

	resp := JobResponse{
		ID:            job.ID,
		Status:        job.Status,
		RunAt:         job.RunAt,
		AttemptsCount: job.AttemptsCount,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to send response")
		return
	}
}

func (h *JobHandler) HandleGetAttempts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	jobID, ok := mux.Vars(r)["job_id"]
	if !ok {
		writeError(w, http.StatusBadRequest, "Provide job_id")
		return
	}

	attempts, err := h.jobService.GetAttempts(ctx, jobID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get attempts")
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
		writeError(w, http.StatusInternalServerError, "Failed to send response")
		return
	}
}

func (h *JobHandler) HandleCancelJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	jobID, ok := mux.Vars(r)["job_id"]
	if !ok {
		writeError(w, http.StatusBadRequest, "Provide job_id")
		return
	}

	if err := h.jobService.CancelJob(ctx, jobID); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to cancel the job")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := ErrorMessage{
		Code:    code,
		Message: msg,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to send error message", http.StatusInternalServerError)
		return
	}
}
