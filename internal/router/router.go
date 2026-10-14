package router

import (
	"net/http"

	"github.com/aleksandrboro/hookflow/internal/handler"
	"github.com/gorilla/mux"
)

func NewRouter(h *handler.JobHandler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", h.CreateJob).Methods("POST")
	r.HandleFunc("/{job_id}", h.GetJob).Methods("GET")
	r.HandleFunc("/{job_id}/attempts", h.GetAttempts).Methods("GET")
	r.HandleFunc("/{job_id}/cancel", h.CancelJob).Methods("POST")

	http.Handle("/v1/jobs", r)

	return r
}
