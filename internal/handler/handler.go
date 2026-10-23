package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/aleksandrboro/hookflow/internal/domain"
	"github.com/aleksandrboro/hookflow/internal/storage"
	"github.com/gorilla/mux"
)

type JobHandler struct {
	jobService JobService
	log        *slog.Logger
}

type JobService interface {
	CreateJob(context.Context, CreateJobRequest) (string, error)
	GetJob(context.Context, string) (domain.Job, error)
	GetAttempts(context.Context, string) ([]domain.Attempt, error)
	CancelJob(context.Context, string) error
}

func NewJobHandler(service JobService, log *slog.Logger) *JobHandler {
	return &JobHandler{
		jobService: service,
		log:        log,
	}
}

func (h *JobHandler) HandleCreateJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	r.Body = http.MaxBytesReader(w, r.Body, 80*1024)

	var req CreateJobRequest

	decoder := json.NewDecoder(r.Body)

	w.Header().Set("Content-Type", "application/json")

	if err := decoder.Decode(&req); err != nil {
		var maxBytesErr *http.MaxBytesError

		if errors.As(err, &maxBytesErr) {
			h.writeError(w, http.StatusRequestEntityTooLarge, "request is too large")
			return
		}

		h.writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	var extra any
	if err := decoder.Decode(extra); err != io.EOF {
		h.writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if len(req.Payload) > 64*1024 {
		h.writeError(w, http.StatusRequestEntityTooLarge, "payload is too large")
		return
	}

	if err := req.Validate(); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.jobService.CreateJob(ctx, req)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "failed to create job")
		return
	}

	resp := JobResponse{
		ID:     id,
		Status: domain.JobStatusPending,
		RunAt:  req.RunAt,
	}

	body, err := json.Marshal(resp)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "failed to marshal response")
		return
	}

	w.WriteHeader(http.StatusAccepted)
	if _, err := w.Write(body); err != nil {
		h.log.Error("failed to send response", slog.String("error", err.Error()))
		return
	}
}

func (h *JobHandler) HandleGetJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	jobID, ok := mux.Vars(r)["job_id"]
	if !ok {
		h.writeError(w, http.StatusBadRequest, "provide job_id")
		return
	}

	job, err := h.jobService.GetJob(ctx, jobID)
	if err != nil {
		if errors.Is(err, storage.ErrJobNotExists) {
			h.writeError(w, http.StatusNotFound, "job not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "failed to get job")
		return
	}

	resp := JobResponse{
		ID:            job.ID,
		Status:        job.Status,
		RunAt:         job.RunAt,
		AttemptsCount: job.AttemptsCount,
	}

	body, err := json.Marshal(resp)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "failed to marshal response")
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(body); err != nil {
		h.log.Error("failed to send response", slog.String("error", err.Error()))
		return
	}
}

func (h *JobHandler) HandleGetAttempts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	jobID, ok := mux.Vars(r)["job_id"]
	if !ok {
		h.writeError(w, http.StatusBadRequest, "provide job_id")
		return
	}

	attempts, err := h.jobService.GetAttempts(ctx, jobID)
	if err != nil {
		if errors.Is(err, storage.ErrJobNotExists) {
			h.writeError(w, http.StatusNotFound, "job not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "failed to get attempts")
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

	body, err := json.Marshal(resp)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "failed to marshal response")
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(body); err != nil {
		h.log.Error("failed to send response", slog.String("error", err.Error()))
		return
	}
}

func (h *JobHandler) HandleCancelJob(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	jobID, ok := mux.Vars(r)["job_id"]
	if !ok {
		h.writeError(w, http.StatusBadRequest, "provide job_id")
		return
	}

	if err := h.jobService.CancelJob(ctx, jobID); err != nil {
		if errors.Is(err, storage.ErrJobNotExists) {
			h.writeError(w, http.StatusNotFound, "job not found")
			return
		}
		if errors.Is(err, storage.ErrFailedCanceletion) {
			h.writeError(w, http.StatusConflict, "failed to cancel the job")
			return
		}

		h.writeError(w, http.StatusInternalServerError, "failed to cancel the job")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *JobHandler) writeError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := ErrorMessage{
		Code:    code,
		Message: msg,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.log.Error("failed to send error message", slog.String("error", err.Error()))
		return
	}
}
