package server

import (
	"errors"
	"fmt"
	"net"
	"net/http"
)

type Server struct {
	listener net.Listener
}

func New(port int) (*Server, error) {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("net.Listen(%s): %w", addr, err)
	}

	return &Server{
		listener: listener,
	}, nil
}

func (s *Server) ServeHTTP(srv *http.Server) error {
	if err := srv.Serve(s.listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("srv.Serve: %w", err)
	}

	return nil
}
