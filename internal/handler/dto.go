package handler

import (
	"encoding/json"
	"errors"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/aleksandrboro/hookflow/internal/domain"
)

var (
	errInvalidMaxAttempts = errors.New("max attempts must be between 1 and 10")
	errEmptyURL           = errors.New("url cannot be empty")
	errInvalidURL         = errors.New("invalid url")
	errPayloadToLarge     = errors.New("payload to large")
	errInvalidPayload     = errors.New("invalid payload")
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

	if r.URL == "" {
		return errEmptyURL
	}

	if err := validateTargetURL(r.URL); err != nil {
		return err
	}

	if !json.Valid(r.Payload) {
		return errInvalidPayload
	}

	var data any

	if err := json.Unmarshal(r.Payload, &data); err != nil {
		return err
	}

	if _, ok := data.(map[string]any); !ok {
		return errInvalidPayload
	}

	return nil
}

func validateTargetURL(raw string) error {
	if raw == "" {
		return errEmptyURL
	}
	if strings.TrimSpace(raw) != raw {
		return errInvalidURL
	}

	parsedURL, err := url.Parse(raw)
	if err != nil {
		return errInvalidURL
	}

	scheme := strings.ToLower(parsedURL.Scheme)
	if scheme != "http" && scheme != "https" {
		return errInvalidURL
	}

	if !parsedURL.IsAbs() || parsedURL.Opaque != "" || parsedURL.Hostname() == "" {
		return errInvalidURL
	}

	if parsedURL.User != nil || parsedURL.Fragment != "" {
		return errInvalidURL
	}

	if portText := parsedURL.Port(); portText != "" {
		port, err := strconv.Atoi(portText)
		if err != nil || port < 1 || port > 65535 {
			return errInvalidURL
		}
	}

	return nil
}
