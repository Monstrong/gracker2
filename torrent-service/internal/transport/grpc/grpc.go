package grpc

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	pb "github.com/monstrong/gracker2/proto/gen/go/torrent/v1" // <-- при смене версии изменить
	"github.com/monstrong/gracker2/torrent-service/internal/models"
	"github.com/monstrong/gracker2/torrent-service/internal/service"
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
	id, err := uuid.Parse(in.GetId())
	if err != nil {
		return nil, fmt.Errorf("%w: invalid torrent id", models.ErrInvalidData)
	}

	m, err := s.s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &pb.GetTorrentResponse{Torrent: toProto(m)}, nil
}

func (s *Server) ListTorrents(ctx context.Context, in *pb.ListTorrentsRequest) (*pb.ListTorrentsResponse, error) {
	ms, err := s.s.List(ctx, int(in.GetLimit()), int(in.GetOffset()))
	if err != nil {
		return nil, err
	}
	pbTorrents := make([]*pb.Torrent, 0, len(ms))

	for _, m := range ms {
		pbTorrent := toProto(m)
		pbTorrents = append(pbTorrents, pbTorrent)
	}

	return &pb.ListTorrentsResponse{Torrents: pbTorrents}, nil
}

func (s *Server) CreateTorrent(ctx context.Context, in *pb.CreateTorrentRequest) (*pb.CreateTorrentResponse, error) {
	author_id, err := uuid.Parse(in.GetAuthorId())
	if err != nil {
		return nil, fmt.Errorf("%w: invalid format for author_id", models.ErrInvalidData)
	}
	category_id, err := uuid.Parse(in.GetCategoryId())
	if err != nil {
		return nil, fmt.Errorf("%w: invalid format for category_id", models.ErrInvalidData)
	}
	t := models.Torrent{Name: in.GetName(), 
		Description: in.GetDescription(), 
		InfoHash: in.GetInfoHash(), 
		AuthorID: author_id,
		CategoryID: category_id,
	}
	m, err := s.s.Create(ctx, &t)
	if err != nil{
		return nil, err
	}

	return &pb.CreateTorrentResponse{Torrent: toProto(m)}, nil
}


func (s *Server) DeleteTorrent(ctx context.Context, in *pb.DeleteTorrentRequest) (*pb.DeleteTorrentResponse, error) {
	id, err := uuid.Parse(in.GetId())
	if err != nil {
		return nil, fmt.Errorf("%w: invalid format for torrent id", models.ErrInvalidData)
	}

	ok, err := s.s.Delete(ctx, id)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteTorrentResponse{Deleted: ok}, nil
}


func toProto(m *models.Torrent) *pb.Torrent {
	return &pb.Torrent{
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
	}
}