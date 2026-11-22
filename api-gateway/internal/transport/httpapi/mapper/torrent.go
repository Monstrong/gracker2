package mapper

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/monstrong/gracker2/api-gateway/internal/transport/httpapi/dto"
	torrentv1 "github.com/monstrong/gracker2/proto/gen/go/torrent/v1"
)

func ToCreateTorrentGRPC(req dto.CreateTorrentRequest) *torrentv1.CreateTorrentRequest {
	return &torrentv1.CreateTorrentRequest{
		Name:        req.Name,
		Description: req.Description,
		InfoHash:    req.InfoHash,
		CategoryId:  req.CategoryID.String(),
	}
}

func ToCreateTorrentHTTP(resp *torrentv1.CreateTorrentResponse) (dto.CreateTorrentResponse, error) {
	if resp == nil || resp.Torrent == nil {
		return dto.CreateTorrentResponse{}, fmt.Errorf("empty torrent response")
	}

	t, err := ToTorrentHTTP(resp.Torrent)
	if err != nil {
		return dto.CreateTorrentResponse{}, err
	}

	return dto.CreateTorrentResponse{Torrent: t}, nil
}

func ToGetTorrentHTTP(resp *torrentv1.GetTorrentResponse) (dto.GetTorrentResponse, error) {
	if resp == nil || resp.Torrent == nil {
		return dto.GetTorrentResponse{}, fmt.Errorf("empty torrent response")
	}

	t, err := ToTorrentHTTP(resp.Torrent)
	if err != nil {
		return dto.GetTorrentResponse{}, err
	}

	return dto.GetTorrentResponse{Torrent: t}, nil
}

func ToListTorrentsHTTP(resp *torrentv1.ListTorrentsResponse) (dto.ListTorrentsResponse, error) {
	if resp == nil {
		return dto.ListTorrentsResponse{}, fmt.Errorf("empty torrents response")
	}

	torrents := make([]dto.Torrent, 0, len(resp.Torrents))

	for _, torrent := range resp.Torrents {
		t, err := ToTorrentHTTP(torrent)
		if err != nil {
			return dto.ListTorrentsResponse{}, err
		}

		torrents = append(torrents, t)
	}

	return dto.ListTorrentsResponse{Torrents: torrents}, nil
}

func ToUpdateTorrentGRPC(id uuid.UUID, req dto.UpdateTorrentRequest) *torrentv1.UpdateTorrentRequest {
	grpcReq := &torrentv1.UpdateTorrentRequest{
		Id:          id.String(),
		Name:        req.Name,
		Description: req.Description,
		InfoHash:    req.InfoHash,
	}

	if req.CategoryID != nil {
		categoryID := req.CategoryID.String()
		grpcReq.CategoryId = &categoryID
	}

	return grpcReq
}

func ToUpdateTorrentHTTP(resp *torrentv1.UpdateTorrentResponse) (dto.UpdateTorrentResponse, error) {
	if resp == nil || resp.Torrent == nil {
		return dto.UpdateTorrentResponse{}, fmt.Errorf("empty torrent response")
	}

	t, err := ToTorrentHTTP(resp.Torrent)
	if err != nil {
		return dto.UpdateTorrentResponse{}, err
	}

	return dto.UpdateTorrentResponse{Torrent: t}, nil
}

func ToTorrentHTTP(pb *torrentv1.Torrent) (dto.Torrent, error) {
	if pb == nil {
		return dto.Torrent{}, fmt.Errorf("empty torrent")
	}

	id, err := uuid.Parse(pb.Id)
	if err != nil {
		return dto.Torrent{}, fmt.Errorf("invalid torrent UUID: %w", err)
	}

	authorID, err := uuid.Parse(pb.AuthorId)
	if err != nil {
		return dto.Torrent{}, fmt.Errorf("invalid author UUID: %w", err)
	}

	categoryID, err := uuid.Parse(pb.CategoryId)
	if err != nil {
		return dto.Torrent{}, fmt.Errorf("invalid category UUID: %w", err)
	}

	if pb.CreatedAt == nil || pb.UpdatedAt == nil {
		return dto.Torrent{}, fmt.Errorf("missing torrent timestamps")
	}

	return dto.Torrent{
		ID:          id,
		Name:        pb.Name,
		Description: pb.Description,
		InfoHash:    pb.InfoHash,
		AuthorID:    authorID,
		CategoryID:  categoryID,
		Status:      pb.Status,
		Downloads:   pb.Downloads,
		CreatedAt:   pb.CreatedAt.AsTime(),
		UpdatedAt:   pb.UpdatedAt.AsTime(),
	}, nil
}

func ToDeleteTorrentHTTP(resp *torrentv1.DeleteTorrentResponse) dto.DeleteTorrentResponse {
	return dto.DeleteTorrentResponse{
		Deleted: resp.GetDeleted(),
	}
}

func ToUpdateStatusHTTP(resp *torrentv1.UpdateStatusResponse) dto.UpdateStatusResponse {
	return dto.UpdateStatusResponse{
		Updated: resp.GetUpdated(),
	}
}

func ToListCategoryHTTP(resp *torrentv1.ListCategoryResponse) (dto.ListCategoryResponse, error) {
	if resp == nil {
		return dto.ListCategoryResponse{}, fmt.Errorf("empty categories response")
	}

	categories := make([]dto.Category, 0, len(resp.Categories))

	for _, category := range resp.Categories {
		if category == nil {
			continue
		}

		id, err := uuid.Parse(category.Id)
		if err != nil {
			return dto.ListCategoryResponse{}, fmt.Errorf("invalid category UUID: %w", err)
		}

		item := dto.Category{
			ID:   id,
			Name: category.Name,
		}

		if category.ParentId != nil {
			parentID, err := uuid.Parse(*category.ParentId)
			if err != nil {
				return dto.ListCategoryResponse{}, fmt.Errorf("invalid parent category UUID: %w", err)
			}

			item.ParentID = &parentID
		}

		categories = append(categories, item)
	}

	return dto.ListCategoryResponse{Categories: categories}, nil
}
