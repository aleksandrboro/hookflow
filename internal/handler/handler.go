package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/aleksandrboro/hookflow/internal/domain"
	requestid "github.com/aleksandrboro/hookflow/internal/request_id"
	"github.com/aleksandrboro/hookflow/internal/storage"
	"github.com/gorilla/mux"
)

const (
	maxRequestBodyBytes = 80 * 1024
	maxPayloadBytes     = 64 * 1024
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
	if log == nil {
		log = slog.Default()
	}

	return &JobHandler{
		jobService: service,
		log:        log,
	}
}

func (h *JobHandler) HandleCreateJob(w http.ResponseWriter, r *http.Request) {
	op := "handler.HandleCreateJob"

	ctx := r.Context()
	requestID, err := requestid.RequestIDFromContext(ctx)
	if err != nil {
		h.log.Error("failed to get request_id from context",
			slog.String("op", op),
			slog.Any("error", err),
		)
		h.writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	log := h.log.With(slog.String("op", op), slog.String("request_id", requestID))

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)

	var req CreateJobRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	w.Header().Set("Content-Type", "application/json")

	if err := decoder.Decode(&req); err != nil {
		h.writeDecodeError(w, err)
		return
	}

	var extra json.RawMessage
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			h.writeError(w, http.StatusBadRequest, "request must contain exactly one JSON value")
		} else {
			h.writeDecodeError(w, err)
		}
		return
	}

	if len(req.Payload) > maxPayloadBytes {
		h.writeError(w, http.StatusRequestEntityTooLarge, "payload is too large")
		return
	}

	if err := req.Validate(); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.jobService.CreateJob(ctx, req)
	if err != nil {
		log.Error("create job failed", slog.Any("error", err))
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
		log.Error("marshal job response failed", slog.Any("error", err))
		h.writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.WriteHeader(http.StatusAccepted)
	if _, err := w.Write(body); err != nil {
		log.Error("failed to send response",
			slog.String("job_id", id),
			slog.String("error", err.Error()),
		)
	}

	log.Info("job created",
		slog.String("job_id", id),
		slog.String("status", string(resp.Status)),
	)
}

func (h *JobHandler) HandleGetJob(w http.ResponseWriter, r *http.Request) {
	op := "handler.HandleGetJob"

	ctx := r.Context()
	requestID, err := requestid.RequestIDFromContext(ctx)
	if err != nil {
		h.log.Error("failed to get request_id from context",
			slog.String("op", op),
			slog.Any("error", err),
		)
		h.writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	log := h.log.With(slog.String("op", op), slog.String("request_id", requestID))

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

		log.Error("failed to get job", slog.Any("error", err))
		h.writeError(w, http.StatusInternalServerError, "internal error")
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
		log.Error("failed to marshal response", slog.Any("error", err))
		h.writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(body); err != nil {
		log.Error("failed to send response",
			slog.String("job_id", jobID),
			slog.String("error", err.Error()),
		)
	}

	log.Info("job received",
		slog.String("job_id", job.ID),
		slog.String("status", string(job.Status)),
	)
}

func (h *JobHandler) HandleGetAttempts(w http.ResponseWriter, r *http.Request) {
	op := "handler.HandleGetAttempts"

	ctx := r.Context()
	requestID, err := requestid.RequestIDFromContext(ctx)
	if err != nil {
		h.log.Error("failed to get request_id from context",
			slog.String("op", op),
			slog.Any("error", err),
		)
		h.writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	log := h.log.With(slog.String("op", op), slog.String("request_id", requestID))

	w.Header().Set("Content-Type", "application/json")

	jobID, ok := mux.Vars(r)["job_id"]
	if !ok {
		h.writeError(w, http.StatusBadRequest, "provide job_id")
		return
	}

	attempts, err := h.jobService.GetAttempts(r.Context(), jobID)
	if err != nil {
		if errors.Is(err, storage.ErrJobNotExists) {
			h.writeError(w, http.StatusNotFound, "job not found")
			return
		}
		log.Error("failed to get attempts", slog.Any("error", err))
		h.writeError(w, http.StatusInternalServerError, "internal error")
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
		log.Error("failed to marshal response", slog.Any("error", err))
		h.writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(body); err != nil {
		log.Error("failed to send response",
			slog.String("job_id", jobID),
			slog.String("error", err.Error()),
		)
	}

	log.Info("attempts received", slog.String("job_id", jobID))
}

func (h *JobHandler) HandleCancelJob(w http.ResponseWriter, r *http.Request) {
	op := "handler.HandleCancelJob"

	ctx := r.Context()
	requestID, err := requestid.RequestIDFromContext(ctx)
	if err != nil {
		h.log.Error("failed to get request_id from context",
			slog.String("op", op),
			slog.Any("error", err),
		)
		h.writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	log := h.log.With(slog.String("op", op), slog.String("request_id", requestID))

	jobID, ok := mux.Vars(r)["job_id"]
	if !ok {
		h.writeError(w, http.StatusBadRequest, "provide job_id")
		return
	}

	if err := h.jobService.CancelJob(r.Context(), jobID); err != nil {
		if errors.Is(err, storage.ErrJobNotExists) {
			h.writeError(w, http.StatusNotFound, "job not found")
			return
		}
		if errors.Is(err, storage.ErrFailedCanceletion) {
			h.writeError(w, http.StatusConflict, "failed to cancel the job")
			return
		}

		log.Error("failed to cancel job", slog.String("job_id", jobID), slog.Any("error", err))
		h.writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.WriteHeader(http.StatusNoContent)

	log.Info("job canceled", slog.String("job_id", jobID))
}

func (h *JobHandler) writeDecodeError(w http.ResponseWriter, err error) {
	var maxBytesErr *http.MaxBytesError
	if errors.As(err, &maxBytesErr) {
		h.writeError(w, http.StatusRequestEntityTooLarge, "request is too large")
		return
	}

	h.writeError(w, http.StatusBadRequest, "invalid json")
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
	}
}
