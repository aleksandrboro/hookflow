package domain

import "time"

type Attempt struct {
	ID               string
	JobID            string
	Number           int
	Status           string
	StartedAt        time.Time
	FinishedAt       time.Time
	ErrorDescription string
}
