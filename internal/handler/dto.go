package handler

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/aleksandrboro/hookflow/internal/domain"
)

var (
	errInvalidMaxAttempts = errors.New("max attempts must be between 1 and 10")
)

type CreateJobRequest struct {
	URL         string          `json:"url"`
	Payload     json.RawMessage `json:"payload"`
	RunAt       time.Time       `json:"run_at,omitzero"`
	MaxAttempts int             `json:"max_attempts,omitempty"`
}

type JobResponse struct {
	ID            string           `json:"id"`
	Status        domain.JobStatus `json:"status"`
	RunAt         time.Time        `json:"run_at"`
	AttemptsCount int              `json:"attempts_count"`
}

type AttemptsResponse struct {
	Data []Attempt `json:"data"`
}

type Attempt struct {
	ID               string    `json:"id"`
	Number           int       `json:"number"`
	Status           string    `json:"status"`
	StartedAt        time.Time `json:"started_at"`
	FinishedAt       time.Time `json:"finished_at"`
	ErrorDescription string    `json:"error_description"`
}

type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (r *CreateJobRequest) Validate() error {
	if r.RunAt.IsZero() {
		r.RunAt = time.Now().UTC()
	}

	if r.MaxAttempts == 0 {
		r.MaxAttempts = 5
	}

	if r.MaxAttempts < 0 || r.MaxAttempts > 10 {
		return errInvalidMaxAttempts
	}

	return nil
}
