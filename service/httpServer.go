package service

import (
	"context"
	"net/http"
	"time"
)

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Service) initHTTPServer() {
	s.mux = http.NewServeMux()
}

// startHTTPServer initializes the server and starts a goroutine with a http server
// Service.mu is held when called
func (s *Service) startHTTPServer() {
	s.Log("Starting HTTP server")
	s.h = &http.Server{Addr: s.cfg.portString, Handler: s.mux}

	go func() {
		s.Logf("Listening on %s://%s%s", s.cfg.scheme, "0.0.0.0", s.cfg.portString)

		var err error
		if s.cfg.TLS {
			err = s.h.ListenAndServeTLS(s.cfg.CertFile, s.cfg.KeyFile)
		} else {
			err = s.h.ListenAndServe()
		}

		if err != nil {
			s.Log(err)
			s.Stop()
		}
	}()
}

// stopHTTPServer stops the http server
func (s *Service) stopHTTPServer() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.h == nil {
		return
	}

	s.Log("Stopping HTTP server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.h.Shutdown(ctx)
	s.h = nil

	if ctx.Err() == context.DeadlineExceeded {
		s.Log("HTTP server forcefully stopped after timeout")
	} else {
		s.Log("HTTP server gracefully stopped")
	}
}