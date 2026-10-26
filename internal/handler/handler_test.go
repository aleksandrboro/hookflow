package handler_test

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aleksandrboro/hookflow/internal/domain"
	"github.com/aleksandrboro/hookflow/internal/handler"
)

type fakeJobService struct {
	createID       string
	createErr      error
	createCalled   bool
	createdRequest handler.CreateJobRequest
}

func (f *fakeJobService) CreateJob(_ context.Context, req handler.CreateJobRequest) (string, error) {
	f.createCalled = true
	f.createdRequest = req
	return f.createID, f.createErr
}

func (f *fakeJobService) GetJob(context.Context, string) (domain.Job, error) {
	return domain.Job{}, nil
}

func (f *fakeJobService) GetAttempts(context.Context, string) ([]domain.Attempt, error) {
	return nil, nil
}

func (f *fakeJobService) CancelJob(context.Context, string) error {
	return nil
}

func TestHandleCreateJobAcceptsValidRequest(t *testing.T) {
	service := &fakeJobService{createID: "job-123"}
	testHandler := newTestJobHandler(service)
	body := `{"url":"https://example.com/webhook","payload":{}}` + " \n\t"
	request := httptest.NewRequest(http.MethodPost, "/v1/jobs", strings.NewReader(body))
	recorder := httptest.NewRecorder()

	testHandler.HandleCreateJob(recorder, request)

	if recorder.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d; body: %s", recorder.Code, http.StatusAccepted, recorder.Body.String())
	}
	if !service.createCalled {
		t.Fatal("CreateJob was not called")
	}
	if service.createdRequest.MaxAttempts != 5 {
		t.Errorf("MaxAttempts = %d, want 5", service.createdRequest.MaxAttempts)
	}
	if service.createdRequest.RunAt.IsZero() {
		t.Error("RunAt was not defaulted")
	}

	var response handler.JobResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	if response.ID != "job-123" {
		t.Errorf("response ID = %q, want %q", response.ID, "job-123")
	}
	if response.Status != domain.JobStatusPending {
		t.Errorf("response status = %q, want %q", response.Status, domain.JobStatusPending)
	}
	if response.RunAt.IsZero() {
		t.Error("response RunAt is zero")
	}
}

func TestHandleCreateJobRejectsSecondJSONValue(t *testing.T) {
	service := &fakeJobService{createID: "job-123"}
	handler := newTestJobHandler(service)
	body := `{"url":"https://example.com/webhook","payload":{}} {"extra":true}`
	request := httptest.NewRequest(http.MethodPost, "/v1/jobs", strings.NewReader(body))
	recorder := httptest.NewRecorder()

	handler.HandleCreateJob(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
	if service.createCalled {
		t.Error("CreateJob was called for a request with a second JSON value")
	}
}

func TestHandleCreateJobRejectsUnknownField(t *testing.T) {
	service := &fakeJobService{createID: "job-123"}
	handler := newTestJobHandler(service)
	body := `{"url":"https://example.com/webhook","payload":{},"max_attemps":3}`
	request := httptest.NewRequest(http.MethodPost, "/v1/jobs", strings.NewReader(body))
	recorder := httptest.NewRecorder()

	handler.HandleCreateJob(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
	if service.createCalled {
		t.Error("CreateJob was called for a request with an unknown field")
	}
}

func TestHandleCreateJobRejectsOversizedPayload(t *testing.T) {
	service := &fakeJobService{createID: "job-123"}
	handler := newTestJobHandler(service)
	largePayload := `{"data":"` + strings.Repeat("x", 64*1024) + `"}`
	body := `{"url":"https://example.com/webhook","payload":` + largePayload + `}`
	request := httptest.NewRequest(http.MethodPost, "/v1/jobs", strings.NewReader(body))
	recorder := httptest.NewRecorder()

	handler.HandleCreateJob(recorder, request)

	if recorder.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusRequestEntityTooLarge)
	}
	if service.createCalled {
		t.Error("CreateJob was called for an oversized payload")
	}
}

func TestHandleCreateJobMapsBodyLimitErrorTo413(t *testing.T) {
	service := &fakeJobService{createID: "job-123"}
	handler := newTestJobHandler(service)
	validJSON := `{"url":"https://example.com/webhook","payload":{}}`
	body := validJSON + strings.Repeat(" ", 80*1024)
	request := httptest.NewRequest(http.MethodPost, "/v1/jobs", strings.NewReader(body))
	recorder := httptest.NewRecorder()

	handler.HandleCreateJob(recorder, request)

	if recorder.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusRequestEntityTooLarge)
	}
	if service.createCalled {
		t.Error("CreateJob was called for a request over the body limit")
	}
}

func newTestJobHandler(service *fakeJobService) *handler.JobHandler {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return handler.NewJobHandler(service, logger)
}
