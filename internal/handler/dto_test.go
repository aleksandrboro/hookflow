package handler_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/aleksandrboro/hookflow/internal/handler"
)

func TestCreateJobRequestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		target  string
		wantErr bool
	}{
		{name: "HTTPS URL", target: "https://example.com/webhook"},
		{name: "HTTP URL with port", target: "http://mock-receiver:8090/fail-once"},
		{name: "case-insensitive scheme", target: "HTTPS://example.com/webhook"},
		{name: "maximum port", target: "https://example.com:65535/webhook"},
		{name: "empty URL", target: "", wantErr: true},
		{name: "relative path", target: "/webhook", wantErr: true},
		{name: "unsupported scheme", target: "ftp://example.com/webhook", wantErr: true},
		{name: "missing host", target: "http:///webhook", wantErr: true},
		{name: "credentials", target: "https://user:secret@example.com/webhook", wantErr: true},
		{name: "fragment", target: "https://example.com/webhook#section", wantErr: true},
		{name: "empty port", target: "https://example.com:/webhook", wantErr: true},
		{name: "zero port", target: "https://example.com:0/webhook", wantErr: true},
		{name: "port too large", target: "https://example.com:65536/webhook", wantErr: true},
		{name: "leading whitespace", target: " https://example.com/webhook", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validCreateJobRequest(tt.target)
			err := req.Validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateJobRequestValidatePayloadMustBeObject(t *testing.T) {
	tests := []struct {
		name    string
		payload json.RawMessage
		wantErr bool
	}{
		{name: "empty object is allowed", payload: json.RawMessage(`{}`)},
		{name: "object is allowed", payload: json.RawMessage(`{"event":"invoice.paid"}`)},
		{name: "missing payload", payload: nil, wantErr: true},
		{name: "null is not an object", payload: json.RawMessage(`null`), wantErr: true},
		{name: "array is not an object", payload: json.RawMessage(`[]`), wantErr: true},
		{name: "string is not an object", payload: json.RawMessage(`"value"`), wantErr: true},
		{name: "malformed JSON", payload: json.RawMessage(`{"event":`), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validCreateJobRequest("https://example.com/webhook")
			req.Payload = tt.payload
			err := req.Validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateJobRequestValidateDefaults(t *testing.T) {
	req := handler.CreateJobRequest{
		URL:     "https://example.com/webhook",
		Payload: json.RawMessage(`{}`),
	}

	if err := req.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if req.MaxAttempts != 5 {
		t.Errorf("MaxAttempts = %d, want 5", req.MaxAttempts)
	}
	if req.RunAt.IsZero() {
		t.Error("RunAt should be defaulted")
	}
	if req.RunAt.Location() != time.UTC {
		t.Errorf("RunAt location = %v, want UTC", req.RunAt.Location())
	}
}

func TestCreateJobRequestValidateMaxAttempts(t *testing.T) {
	tests := []struct {
		name     string
		attempts int
	}{
		{name: "negative", attempts: -1},
		{name: "above maximum", attempts: 11},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := validCreateJobRequest("https://example.com/webhook")
			req.MaxAttempts = tt.attempts
			if err := req.Validate(); err == nil {
				t.Fatalf("Validate() with MaxAttempts=%d succeeded, want error", tt.attempts)
			}
		})
	}
}

func validCreateJobRequest(targetURL string) handler.CreateJobRequest {
	return handler.CreateJobRequest{
		URL:         targetURL,
		Payload:     json.RawMessage(`{"event":"test"}`),
		RunAt:       time.Now().UTC(),
		MaxAttempts: 1,
	}
}
