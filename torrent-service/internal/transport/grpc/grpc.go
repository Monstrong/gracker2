package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/monstrong/gracker2/pkg/errs"
	pb "github.com/monstrong/gracker2/proto/gen/go/torrent/v1" // <-- при смене версии изменить
	"github.com/monstrong/gracker2/torrent-service/internal/models"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Тут нет смысла добавлять "op" в ошибки. вызываемый метод уже известен клиенту

// StatusUpdate wants SenderRole in metadata
func (s *Server) UpdateStatus(ctx context.Context, in *pb.UpdateStatusRequest) (*pb.UpdateStatusResponse, error) {
	torrentID, err := uuid.Parse(in.GetTorrentId())
	if err != nil {
		return nil, &errs.ValidationError{Field: "id", Reason: "invalid torrent id format"}
	}
	MD, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errs.ErrMissingMetadata
	}
	SenderRole, err := checkSenderRolefromMD(MD)
	if err != nil {
		return nil, err
	}
	ok, err = s.s.StatusUpdate(ctx, torrentID, int(in.GetNewStatus()), SenderRole)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateStatusResponse{Updated: ok}, nil
}

// UpdateTorrent wants SenderRole AND SenderID in metadata
func (s *Server) UpdateTorrent(ctx context.Context, in *pb.UpdateTorrentRequest) (*pb.UpdateTorrentResponse, error) {
	var categoryID *uuid.UUID

	if in.CategoryId != nil {
		parsed, err := uuid.Parse(in.GetCategoryId())
		if err != nil {
			return nil, &errs.ValidationError{
				Field:  "categoryID",
				Reason: "invalid categoryID format",
			}
		}

		categoryID = &parsed
	}
	torrentID, err := uuid.Parse(in.GetId())
	if err != nil {
		return nil, &errs.ValidationError{Field: "id", Reason: "invalid torrent id format"}
	}
	MD, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errs.ErrMissingMetadata
	}

	SenderRole, err := checkSenderRolefromMD(MD)
	if err != nil {
		return nil, err
	}
	SenderID, err := checkIDfromMD(MD)
	if err != nil {
		return nil, err
	}

	patch := &models.TorrentPatch{
		Name:        in.Name,
		Description: in.Description,
		InfoHash:    in.InfoHash,
		CategoryID:  categoryID,
	}

	m, err := s.s.Update(ctx, patch, torrentID, SenderID, SenderRole)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateTorrentResponse{Torrent: toProto(m)}, nil
}

// GetTorrent wants SenderRole in metadata
func (s *Server) GetTorrent(ctx context.Context, in *pb.GetTorrentRequest) (*pb.GetTorrentResponse, error) {
	id, err := uuid.Parse(in.GetId())
	if err != nil {
		return nil, &errs.ValidationError{Field: "id", Reason: "invalid torrent id format"}
	}

	MD, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errs.ErrMissingMetadata
	}
	SenderRole, err := checkSenderRolefromMD(MD)
	if err != nil {
		return nil, err
	}

	m, err := s.s.GetByID(ctx, id, SenderRole)
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

// CreateTorrent wants SenderRole AND SenderID in metadata
func (s *Server) CreateTorrent(ctx context.Context, in *pb.CreateTorrentRequest) (*pb.CreateTorrentResponse, error) {
	authorID, err := uuid.Parse(in.GetAuthorId())
	if err != nil {
		return nil, &errs.ValidationError{Field: "authorID", Reason: "invalid authorID id format"}
	}
	categoryID, err := uuid.Parse(in.GetCategoryId())
	if err != nil {
		return nil, &errs.ValidationError{Field: "categoryID", Reason: "invalid categoryID id format"}
	}
	MD, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errs.ErrMissingMetadata
	}

	SenderID, err := checkIDfromMD(MD)
	if err != nil {
		return nil, err
	}
	SenderRole, err := checkSenderRolefromMD(MD)
	if err != nil {
		return nil, err
	}

	t := models.Torrent{Name: in.GetName(),
		Description: in.GetDescription(),
		InfoHash:    in.GetInfoHash(),
		AuthorID:    authorID,
		CategoryID:  categoryID,
	}
	m, err := s.s.Create(ctx, &t, SenderID, SenderRole)
	if err != nil {
		return nil, err
	}

	return &pb.CreateTorrentResponse{Torrent: toProto(m)}, nil
}

// DeleteTorrent wants SenderRole AND SenderID in metadata
func (s *Server) DeleteTorrent(ctx context.Context, in *pb.DeleteTorrentRequest) (*pb.DeleteTorrentResponse, error) {
	id, err := uuid.Parse(in.GetId())
	if err != nil {
		return nil, &errs.ValidationError{Field: "id", Reason: "invalid torrent id format"}
	}

	MD, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errs.ErrMissingMetadata
	}

	SenderID, err := checkIDfromMD(MD)
	if err != nil {
		return nil, err
	}
	SenderRole, err := checkSenderRolefromMD(MD)
	if err != nil {
		return nil, err
	}

	ok, err = s.s.Delete(ctx, id, SenderID, SenderRole)
	if err != nil {
		return nil, err
	}
	return &pb.DeleteTorrentResponse{Deleted: ok}, nil
}

func toProto(m *models.Torrent) *pb.Torrent {
	return &pb.Torrent{
		Id:          m.ID.String(),
		Name:        m.Name,
		Description: m.Description,
		InfoHash:    m.InfoHash,
		AuthorId:    m.AuthorID.String(),
		CategoryId:  m.CategoryID.String(),
		Status:      int32(m.Status),
		Downloads:   int64(m.Downloads),
		CreatedAt:   timestamppb.New(m.CreatedAt),
		UpdatedAt:   timestamppb.New(m.UpdatedAt),
	}
}

