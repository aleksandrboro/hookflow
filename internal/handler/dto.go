package handler

import (
	"encoding/json"
	"time"

	"github.com/aleksandrboro/hookflow/internal/domain"
)

type CreateJobRequest struct {
	URL         string          `json:""`
	Payload     json.RawMessage `json:"payload"`
	RunAt       time.Time       `json:"run_at,omitzero"`
	MaxAttempts string          `json:"max_attempts,omitempty"`
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
