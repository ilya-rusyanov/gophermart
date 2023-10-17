package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	server    *http.Server
	errorChan chan error
}

func New(listenAddr string, handler http.Handler) *Server {
	res := &Server{
		server: &http.Server{
			Addr:    listenAddr,
			Handler: handler,
		},
		errorChan: make(chan error, 1),
	}

	go func() {
		res.errorChan <- res.server.ListenAndServe()
		close(res.errorChan)
	}()

	return res
}

func (s *Server) Error() <-chan error {
	return s.errorChan
}

func (s *Server) Shutdown(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := s.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("failed to shutdown http server: %w", err)
	}
	return nil
}
