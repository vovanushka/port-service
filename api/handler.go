package api

import (
	"encoding/json"

	"github.com/vovanushka/port-service/model"
	"github.com/vovanushka/port-service/repo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"golang.org/x/net/context"
)

// Server represents the gRPC server
type Server struct {
	repo *repo.PortRepo
}

// NewServer is constructor for Server
func NewServer(repo *repo.PortRepo) *Server {
	return &Server{repo}
}

// SavePort calls repo create method and creates grpc error basing on the repo response
func (s Server) SavePort(ctx context.Context, in *PortMessage) (*PortMessage, error) {
	port := model.Port{}

	err := json.Unmarshal(in.Data, &port)
	if err != nil {
		err := status.Error(codes.Internal, "internal port service error")
		return nil, err
	}

	err = s.repo.Create(&port)
	if err != nil {
		err := status.Error(codes.Internal, "internal port service error")
		return nil, err
	}
	return &PortMessage{}, nil
}

// GetPort calls repo create method and creates grpc error basing on repo response
func (s Server) GetPort(ctx context.Context, in *PortIDMessage) (*PortMessage, error) {
	port, err := s.repo.Get(in.Id)
	if err != nil {
		// If port was not found, return not found code error
		if err.Error() == "not found" {
			err := status.Error(codes.NotFound, "port was not found")
			return nil, err
		}
		err = status.Error(codes.Internal, "internal port service error")
		return nil, err
	}
	portJSON, err := json.Marshal(port)
	if err != nil {
		err := status.Error(codes.Internal, "internal port service error")
		return nil, err
	}

	return &PortMessage{Data: portJSON}, nil
}
