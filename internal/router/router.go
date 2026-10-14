package router

import (
	"net/http"

	"github.com/aleksandrboro/hookflow/internal/job"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	h := job.NewJobHandler()

	r.HandleFunc("/", h.CreateJob).Methods("POST")
	r.HandleFunc("/{job_id}", h.GetJob).Methods("GET")
	r.HandleFunc("/{job_id}/attempts", h.GetAttempts).Methods("GET")
	r.HandleFunc("/{job_id}/cancel", h.CancelJob).Methods("POST")

	http.Handle("/v1/jobs", r)

	return r
}
