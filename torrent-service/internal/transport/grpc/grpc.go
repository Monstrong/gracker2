package grpc

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	"github.com/monstrong/gracker2/pkg/contracts/grpcmetakey"
	"github.com/monstrong/gracker2/pkg/domain"
	"github.com/monstrong/gracker2/pkg/errs"
	pb "github.com/monstrong/gracker2/proto/gen/go/torrent/v1" // <-- при смене версии изменить
	"github.com/monstrong/gracker2/torrent-service/internal/models"
	"github.com/monstrong/gracker2/torrent-service/internal/service"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	pb.UnimplementedTorrentServiceServer
	s service.Service
}

func NewServer(s service.Service) *Server {
	return &Server{s: s}
}

// Тут нет смысла добавлять "op" в ошибки. вызываемый метод уже известен клиенту

// StatusUpdate wants SenderRole in metadata
func (s *Server) StatusUpdate(ctx context.Context, in *pb.StatusUpdateRequest) (*pb.StatusUpdateResponse, error) {
	
}

// UpdateTorrent wants SenderRole AND SenderID in metadata
func (s *Server) UpdateTorrent(ctx context.Context, in *pb.UpdateTorrentRequest) (*pb.UpdateTorrentResponse, error) {
	categoryID, err := uuid.Parse(in.GetCategoryId())
	if err != nil {
		return nil, &errs.ValidationError{Field: "categoryID", Reason: "invalid categoryID id format"}
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
		Name:        in.GetName,
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
func (s *Server) GetTorrent(ctx context.Context, in *pb.GetTorrentRequest) (*pb.GetTorrentResponse, error) {
	id, err := uuid.Parse(in.GetId())
	if err != nil {
		return nil, &errs.ValidationError{Field: "id", Reason: "invalid torrent id format"}
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

func checkIDfromMD(MD metadata.MD) (userID uuid.UUID, err error) {
	vals := MD.Get(grpcmetakey.GRPCMetaKeyUserID)

	if len(vals) > 0 {
		if userID, err = uuid.Parse(vals[0]); err != nil {
			return uuid.Nil, &errs.GatewayWrongFieldError{Field: grpcmetakey.GRPCMetaKeyUserID, Reason: "invalid user id format"}
		}
	} else {
		return uuid.Nil, &errs.GatewayWrongFieldError{Field: grpcmetakey.GRPCMetaKeyUserID, Reason: "empty user id"}
	}

	return
}

func checkSenderRolefromMD(MD metadata.MD) (SenderRole domain.Role, err error) {
	vals := MD.Get(grpcmetakey.GRPCMetaKeyRole)

	if len(vals) > 0 {
		roleValue, err := strconv.ParseInt(vals[0], 10, 8)
		if err != nil {
			return domain.RoleUnspecified, &errs.GatewayWrongFieldError{Field: grpcmetakey.GRPCMetaKeyRole, Reason: "invalid role format"}
		}
		SenderRole = domain.Role(roleValue)
	} else {
		return domain.RoleUnspecified, &errs.GatewayWrongFieldError{Field: grpcmetakey.GRPCMetaKeyRole, Reason: "empty role"}
	}

	return
}
