package domain

type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusQueued    JobStatus = "queued"
	JobStatusRunning   JobStatus = "running"
	JobStatusRetrying  JobStatus = "retrying"
	JobStatusSucceeded JobStatus = "succeeded"
	JobStatusFailed    JobStatus = "failed"
	JobStatusDead      JobStatus = "dead"
	JobStatusCancelled JobStatus = "cancelled"
)

var (
	AllowedToCancel = map[JobStatus]struct{}{
		JobStatusPending:  {},
		JobStatusRetrying: {},
		JobStatusQueued:   {},
	}
)
