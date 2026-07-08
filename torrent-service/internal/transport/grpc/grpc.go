package grpc

import (
	"context"
	"errors"

	"github.com/google/uuid"
	pb "github.com/monstrong/gracker2/proto/gen/go/torrent/v1" // при смене версии изменить
	"github.com/monstrong/gracker2/torrent-service/internal/models"
	"github.com/monstrong/gracker2/torrent-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedTorrentServiceServer
	s service.Service
}

func NewServer(s service.Service) *Server{
	return &Server{s: s}
}

func (s *Server) GetTorrent(ctx context.Context, in *pb.GetTorrentRequest) (*pb.GetTorrentResponse, error) {
	const op = "transport.grpc.getTorrent"
	id, err := uuid.Parse(in.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid format for torrent id")
	}

	m, err := s.s.GetByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotFound):
			return nil, status.Error(codes.NotFound, models.ErrNotFound.Error())
		case errors.Is(err, models.ErrInvalidData):
			return nil, status.Error(codes.InvalidArgument, models.ErrInvalidData.Error())
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}
	return &pb.GetTorrentResponse{Torrent: &pb.Torrent{
		Id: m.ID.String(), 
		Name: m.Name,
		Description: m.Description,
		InfoHash: m.InfoHash,
		AuthorId: m.AuthorID.String(),
		CategoryId: m.CategoryID.String(),
		Status: int32(m.Status),
		Downloads: int64(m.Downloads),
		CreatedAt: timestamppb.New(m.CreatedAt),
		UpdatedAt: timestamppb.New(m.UpdatedAt),
		}}, nil
	
}

// func (s *Server) ListTorrents(ctx context.Context, in *pb.ListTorrentsRequest) (*pb.ListTorrentsResponse, error) {
// 	const op = "transport.grpc.listTorrents"

// }

// func (s *Server) CreateTorrent(ctx context.Context, in *pb.CreateTorrentRequest) (*pb.CreateTorrentResponse, error) {
// 	const op = "transport.grpc.createTorrent"

// }

// func (s *Server) DeleteTorrent(ctx context.Context, in *pb.DeleteTorrentRequest) (*pb.DeleteTorrentResponse, error) {
// 	const op = "transport.grpc.deleteTorrent"

// }
