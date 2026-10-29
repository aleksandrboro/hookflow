package router

import (
	"net/http"

	"github.com/aleksandrboro/hookflow/internal/handler"
	"github.com/aleksandrboro/hookflow/internal/middleware"
	"github.com/gorilla/mux"
)

func NewRouter(h *handler.JobHandler) *mux.Router {
	r := mux.NewRouter()

	r.Use(middleware.RequestID)

	r.HandleFunc("/v1/jobs", h.HandleCreateJob).Methods(http.MethodPost)
	r.HandleFunc("/v1/jobs/{job_id}", h.HandleGetJob).Methods(http.MethodGet)
	// r.HandleFunc("/v1/jobs/{job_id}/attempts", h.HandleGetAttempts).Methods(http.MethodGet)
	r.HandleFunc("/v1/jobs/{job_id}/cancel", h.HandleCancelJob).Methods(http.MethodPost)

	return r
}
