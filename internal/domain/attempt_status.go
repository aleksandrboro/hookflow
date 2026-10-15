package domain

type AttemptStatus string

const (
	AttemptStatusRunning          AttemptStatus = "running"
	AttemptStatusSucceeded        AttemptStatus = "succeeded"
	AttemptStatusRetryableFailure AttemptStatus = "retryable_failure"
	AttemptStatusPermanentFailure AttemptStatus = "permanent_failure"
	AttemptStatusAbandoned        AttemptStatus = "abandoned"
)
