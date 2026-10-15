package domain

import (
	"encoding/json"
	"time"
)

type Job struct {
	ID            string
	TargetURL     string
	Payload       json.RawMessage
	Status        JobStatus
	RunAt         time.Time
	AvailableAt   time.Time
	AttemptsCount int
	MaxAttempts   int
}
