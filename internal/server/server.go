// Package server ...
package server

import (
	"context"
	"net"
	"net/http"

	server "github.com/SerjRamone/metrius/internal/grpc"
	"github.com/SerjRamone/metrius/internal/storage"
	"github.com/SerjRamone/metrius/pkg/logger"
	pb "github.com/SerjRamone/metrius/pkg/metrius_v1"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Server ...
type Server interface {
	Up() error
	Down(context.Context) error
}

// HTTPServer ...
type HTTPServer struct {
	server *http.Server
}

// NewHTTPServer ...
func NewHTTPServer(a string, r chi.Router) *HTTPServer {
	s := &http.Server{
		Addr:    a,
		Handler: r,
	}
	return &HTTPServer{
		server: s,
	}
}

// Up ...
func (s *HTTPServer) Up() error {
	return s.server.ListenAndServe()
}

// Down ...
func (s *HTTPServer) Down(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// GRPCServer ...
type GRPCServer struct {
	address string
	server  *grpc.Server
}

// NewGRPCServer ...
func NewGRPCServer(a string, store storage.Storage) *GRPCServer {
	s := grpc.NewServer()
	server := server.NewMetricsServer(store)
	pb.RegisterMetricsServiceServer(s, server)

	return &GRPCServer{
		address: a,
		server:  s,
	}
}

// Up ...
func (s *GRPCServer) Up() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		logger.Error("can't listen", zap.Error(err))
		return err
	}

	if err := s.server.Serve(listener); err != nil {
		logger.Error("can't serve", zap.Error(err))
		return err
	}

	return nil
}

// Down ...
func (s *GRPCServer) Down(ctx context.Context) error {
	_ = ctx
	s.server.GracefulStop()
	return nil
}
