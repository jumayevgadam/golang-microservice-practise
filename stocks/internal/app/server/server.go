package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"stocks/internal/config"
	"stocks/pkg/connection"
	"stocks/pkg/constants"
	"syscall"
	"time"
)

// Server represent server configurations for this stocks service.
type Server struct {
	server *http.Server
	cfg    config.Config
	psqlDB connection.DB
}

// NewServer creates and returns a new instance of Server.
func NewServer(
	cfg config.Config,
	psqlDB connection.DB,
) *Server {
	return &Server{
		server: nil,
		cfg:    cfg,
		psqlDB: psqlDB,
	}
}

// RunHTTPServer starts http server in goroutines and gracefully shutdown if signal catches.
func (s *Server) RunHTTPServer() error {
	// setup routes.
	mux := s.setupRoutes()

	s.server = &http.Server{
		Addr:         s.cfg.Address(),
		Handler:      mux,
		ReadTimeout:  s.cfg.SrvConfig().ReadTimeOut,
		WriteTimeout: s.cfg.SrvConfig().WriteTimeOut,
	}

	go func() {
		log.Printf("server starting on %+v\n", s.cfg.Address())

		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("s.server.ListenAndServe: %v", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("shutting down server...")

	ctxTimeOut, cancel := context.WithTimeout(context.Background(), constants.SrvTimeOut*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctxTimeOut); err != nil {
		log.Printf("s.server.Shutdown: %+v", err.Error())
	}

	log.Println("stock service successfully shutdown")

	return nil
}
