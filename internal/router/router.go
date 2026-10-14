package router

import (
	"net/http"

	"github.com/aleksandrboro/hookflow/internal/handler"
	"github.com/gorilla/mux"
)

func NewRouter(h *handler.JobHandler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/v1/jobs", h.CreateJob).Methods(http.MethodPost)
	r.HandleFunc("/v1/jobs/{job_id}", h.GetJob).Methods(http.MethodGet)
	r.HandleFunc("/v1/jobs/{job_id}/attempts", h.GetAttempts).Methods(http.MethodGet)
	r.HandleFunc("/v1/jobs/{job_id}/cancel", h.CancelJob).Methods(http.MethodPost)

	return r
}
