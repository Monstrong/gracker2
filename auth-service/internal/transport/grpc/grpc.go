package grpc

import (
	"context"

	pb "github.com/monstrong/gracker2/proto/gen/go/auth/v1" // <-- при смене версии изменить
)
func (s *Server) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	return nil, nil
}

