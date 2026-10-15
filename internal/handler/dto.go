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
