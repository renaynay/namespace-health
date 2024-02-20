package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/renaynay/namespace-health/feeder"
	"github.com/renaynay/namespace-health/reader"
)

type Server struct {
	router *mux.Router
	reader *reader.Reader
	feeder *feeder.Feeder
}

func New(read *reader.Reader, feed *feeder.Feeder) *Server {
	return &Server{
		router: mux.NewRouter(),
		reader: read,
		feeder: feed,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.router.HandleFunc("/health", s.getHealthStatus).Methods("GET")
	s.router.HandleFunc("/feed", s.feed).Methods("GET")
	s.router.HandleFunc("/", s.nothing).Methods("GET")

	// Start listening and serving requests
	fmt.Println("Listening on port 8000...")
	return http.ListenAndServe(":8000", s.router)
}

type Health struct {
	IsHealthier bool `json:"is_healthier"`
	Scale       int  `json:"scale"`
}

// getHealthStatus responds with the current health of the namespace
func (s *Server) getHealthStatus(w http.ResponseWriter, r *http.Request) {
	_ = s.reader.Health()

	w.Header().Set("Content-Type", "application/json")
	health := s.reader.Health()
	json.NewEncoder(w).Encode(Health{IsHealthier: health.IsHealthier, Scale: health.Scale}) // TODO @renaynay:
}

func (s *Server) feed(w http.ResponseWriter, r *http.Request) {
	err := s.feeder.Feed(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("fed stuie pooie :--------)"))
}

func (s *Server) nothing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("nothing to see here"))
}
