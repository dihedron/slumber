package api

import (
	"context"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	UnimplementedVMControlServer
	osClient *OpenStackClient
}

func NewServer(osClient *OpenStackClient) *Server {
	return &Server{
		osClient: osClient,
	}
}

func (s *Server) Start(ctx context.Context, req *StartRequest) (*StartResponse, error) {
	slog.Info("starting VM", "userid", req.UserId)
	err := s.osClient.Start(ctx, req.UserId)
	if err != nil {
		slog.Error("failed to start VM", "userid", req.UserId, "error", err)
		return nil, status.Errorf(codes.Internal, "failed to start VM: %v", err)
	}
	return &StartResponse{
		Message: "VM start request submitted",
		Status:  "STARTING",
	}, nil
}

func (s *Server) Stop(ctx context.Context, req *StopRequest) (*StopResponse, error) {
	slog.Info("stopping VM", "userid", req.UserId)
	err := s.osClient.Stop(ctx, req.UserId)
	if err != nil {
		slog.Error("failed to stop VM", "userid", req.UserId, "error", err)
		return nil, status.Errorf(codes.Internal, "failed to stop VM: %v", err)
	}
	return &StopResponse{
		Message: "VM stop request submitted",
		Status:  "STOPPING",
	}, nil
}

func (s *Server) Pause(ctx context.Context, req *PauseRequest) (*PauseResponse, error) {
	slog.Info("pausing VM", "userid", req.UserId)
	err := s.osClient.Pause(ctx, req.UserId)
	if err != nil {
		slog.Error("failed to pause VM", "userid", req.UserId, "error", err)
		return nil, status.Errorf(codes.Internal, "failed to pause VM: %v", err)
	}
	return &PauseResponse{
		Message: "VM pause request submitted",
		Status:  "PAUSING",
	}, nil
}

func (s *Server) Unpause(ctx context.Context, req *UnpauseRequest) (*UnpauseResponse, error) {
	slog.Info("unpausing VM", "userid", req.UserId)
	err := s.osClient.Unpause(ctx, req.UserId)
	if err != nil {
		slog.Error("failed to unpause VM", "userid", req.UserId, "error", err)
		return nil, status.Errorf(codes.Internal, "failed to unpause VM: %v", err)
	}
	return &UnpauseResponse{
		Message: "VM unpause request submitted",
		Status:  "UNPAUSING",
	}, nil
}

func (s *Server) Status(ctx context.Context, req *StatusRequest) (*StatusResponse, error) {
	slog.Info("checking VM status", "userid", req.UserId)
	vmStatus, err := s.osClient.Status(ctx, req.UserId)
	if err != nil {
		slog.Error("failed to get VM status", "userid", req.UserId, "error", err)
		return nil, status.Errorf(codes.Internal, "failed to get VM status: %v", err)
	}
	return &StatusResponse{
		Status:  vmStatus,
		Message: "VM status retrieved",
	}, nil
}
