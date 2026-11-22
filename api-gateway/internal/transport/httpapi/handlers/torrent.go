
package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/go-chi/chi/v5"

	"github.com/monstrong/gracker2/api-gateway/internal/gatewayconfig"
	"github.com/monstrong/gracker2/api-gateway/internal/middleware"
	"github.com/monstrong/gracker2/api-gateway/internal/transport/httpapi/dto"
	"github.com/monstrong/gracker2/api-gateway/internal/transport/httpapi/mapper"
	"github.com/monstrong/gracker2/pkg/contracts/grpcmetakey"
	"github.com/monstrong/gracker2/pkg/logger"

	torrentv1 "github.com/monstrong/gracker2/proto/gen/go/torrent/v1"
	"google.golang.org/grpc/metadata"
)

type TorrentHandler struct {
	torrentClient TorrentClient
	cfg          gatewayconfig.Torrent
	l            logger.Logger
}

func NewTorrentHandler(
	torrentClient TorrentClient,
	cfg gatewayconfig.Torrent,
	l logger.Logger,
) *TorrentHandler {
	return &TorrentHandler{
		torrentClient: torrentClient,
		cfg:           cfg,
		l:             l,
	}
}

// GetTorrent godoc
// @Summary Get torrent by ID
// @Tags Torrents
// @Produce json
// @Param id path string true "Torrent UUID"
// @Success 200 {object} dto.GetTorrentResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /torrents/{id} [get]
func (h *TorrentHandler) GetTorrent(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.torrent.GetTorrent"

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid torrent id")
		return
	}

	claims, ok := middleware.GetClaims(ctx)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	ctx = metadata.AppendToOutgoingContext(
		ctx,
		grpcmetakey.GRPCMetaKeyRole,
		strconv.Itoa(int(claims.Role)),
	)

	resp, err := h.torrentClient.GetTorrent(ctx, &torrentv1.GetTorrentRequest{
		Id: id.String(),
	})
	if err != nil {
		h.l.Error(ctx, op, logger.Error(err))
		writeGRPCError(w, err)
		return
	}

	httpResp, err := mapper.ToGetTorrentHTTP(resp)
	if err != nil {
		h.l.Error(ctx, op, logger.Error(err))
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if err := writeJSON(w, http.StatusOK, httpResp); err != nil {
		h.l.Error(ctx, op, logger.Error(err))
	}
}

// ListTorrents godoc
// @Summary List torrents
// @Tags Torrents
// @Produce json
// @Param limit query int false "Page size"
// @Param offset query int false "Offset"
// @Success 200 {object} dto.ListTorrentsResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /torrents [get]
func (h *TorrentHandler) ListTorrents(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.torrent.ListTorrents"

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	limit := int32(20)
	offset := int32(0)

	if raw := r.URL.Query().Get("limit"); raw != "" {
		value, err := strconv.ParseInt(raw, 10, 32)
		if err != nil || value <= 0 {
			writeError(w, http.StatusBadRequest, "invalid limit")
			return
		}
		limit = int32(value)
	}

	if raw := r.URL.Query().Get("offset"); raw != "" {
		value, err := strconv.ParseInt(raw, 10, 32)
		if err != nil || value < 0 {
			writeError(w, http.StatusBadRequest, "invalid offset")
			return
		}
		offset = int32(value)
	}

	resp, err := h.torrentClient.ListTorrents(ctx, &torrentv1.ListTorrentsRequest{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		h.l.Error(ctx, op, logger.Error(err))
		writeGRPCError(w, err)
		return
	}

	httpResp, err := mapper.ToListTorrentsHTTP(resp)
	if err != nil {
		h.l.Error(ctx, op, logger.Error(err))
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if err := writeJSON(w, http.StatusOK, httpResp); err != nil {
		h.l.Error(ctx, op, logger.Error(err))
	}
}

// CreateTorrent godoc
// @Summary Create torrent
// @Tags Torrents
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateTorrentRequest true "Torrent data"
// @Success 201 {object} dto.CreateTorrentResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /torrents [post]
func (h *TorrentHandler) CreateTorrent(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.torrent.CreateTorrent"

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	claims, ok := middleware.GetClaims(ctx)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req dto.CreateTorrentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad request")
		return
	}

	ctx = metadata.AppendToOutgoingContext(
		ctx,
		grpcmetakey.GRPCMetaKeyUserID, claims.UserID.String(),
		grpcmetakey.GRPCMetaKeyRole, strconv.Itoa(int(claims.Role)),
	)

	grpcReq := mapper.ToCreateTorrentGRPC(req)
	grpcReq.AuthorId = claims.UserID.String()

	resp, err := h.torrentClient.CreateTorrent(ctx, grpcReq)
	if err != nil {
		h.l.Error(ctx, op, logger.Error(err))
		writeGRPCError(w, err)
		return
	}

	httpResp, err := mapper.ToCreateTorrentHTTP(resp)
	if err != nil {
		h.l.Error(ctx, op, logger.Error(err))
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if err := writeJSON(w, http.StatusCreated, httpResp); err != nil {
		h.l.Error(ctx, op, logger.Error(err))
	}
}

// UpdateTorrent godoc
// @Summary Update torrent
// @Tags Torrents
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Torrent UUID"
// @Param request body dto.UpdateTorrentRequest true "Torrent update data"
// @Success 200 {object} dto.UpdateTorrentResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /torrents/{id} [patch]
func (h *TorrentHandler) UpdateTorrent(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.torrent.UpdateTorrent"

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid torrent id")
		return
	}

	claims, ok := middleware.GetClaims(ctx)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req dto.UpdateTorrentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad request")
		return
	}

	ctx = metadata.AppendToOutgoingContext(
		ctx,
		grpcmetakey.GRPCMetaKeyUserID, claims.UserID.String(),
		grpcmetakey.GRPCMetaKeyRole, strconv.Itoa(int(claims.Role)),
	)

	grpcReq := mapper.ToUpdateTorrentGRPC(id, req)

	resp, err := h.torrentClient.UpdateTorrent(ctx, grpcReq)
	if err != nil {
		h.l.Error(ctx, op, logger.Error(err))
		writeGRPCError(w, err)
		return
	}

	httpResp, err := mapper.ToUpdateTorrentHTTP(resp)
	if err != nil {
		h.l.Error(ctx, op, logger.Error(err))
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if err := writeJSON(w, http.StatusOK, httpResp); err != nil {
		h.l.Error(ctx, op, logger.Error(err))
	}
}

// DeleteTorrent godoc
// @Summary Delete torrent
// @Tags Torrents
// @Produce json
// @Security BearerAuth
// @Param id path string true "Torrent UUID"
// @Success 200 {object} dto.DeleteTorrentResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /torrents/{id} [delete]
func (h *TorrentHandler) DeleteTorrent(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.torrent.DeleteTorrent"

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid torrent id")
		return
	}

	claims, ok := middleware.GetClaims(ctx)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	ctx = metadata.AppendToOutgoingContext(
		ctx,
		grpcmetakey.GRPCMetaKeyUserID, claims.UserID.String(),
		grpcmetakey.GRPCMetaKeyRole, strconv.Itoa(int(claims.Role)),
	)

	resp, err := h.torrentClient.DeleteTorrent(ctx, &torrentv1.DeleteTorrentRequest{
		Id: id.String(),
	})
	if err != nil {
		h.l.Error(ctx, op, logger.Error(err))
		writeGRPCError(w, err)
		return
	}

	httpResp := mapper.ToDeleteTorrentHTTP(resp)

	if err := writeJSON(w, http.StatusOK, httpResp); err != nil {
		h.l.Error(ctx, op, logger.Error(err))
	}
}

// UpdateStatus godoc
// @Summary Update torrent status
// @Tags Torrents
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Torrent UUID"
// @Param request body dto.UpdateStatusRequest true "New status"
// @Success 200 {object} dto.UpdateStatusResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /torrents/{id}/status [patch]
func (h *TorrentHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.torrent.UpdateStatus"

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid torrent id")
		return
	}

	claims, ok := middleware.GetClaims(ctx)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req dto.UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "bad request")
		return
	}

	ctx = metadata.AppendToOutgoingContext(
		ctx,
		grpcmetakey.GRPCMetaKeyRole,
		strconv.Itoa(int(claims.Role)),
	)

	resp, err := h.torrentClient.UpdateStatus(ctx, &torrentv1.UpdateStatusRequest{
		TorrentId: id.String(),
		NewStatus: req.NewStatus,
	})
	if err != nil {
		h.l.Error(ctx, op, logger.Error(err))
		writeGRPCError(w, err)
		return
	}

	httpResp := mapper.ToUpdateStatusHTTP(resp)

	if err := writeJSON(w, http.StatusOK, httpResp); err != nil {
		h.l.Error(ctx, op, logger.Error(err))
	}
}

// ListCategory godoc
// @Summary List torrent categories
// @Tags Categories
// @Produce json
// @Success 200 {object} dto.ListCategoryResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /categories [get]
func (h *TorrentHandler) ListCategory(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.torrent.ListCategory"

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	resp, err := h.torrentClient.ListCategory(ctx, &torrentv1.ListCategoryRequest{})
	if err != nil {
		h.l.Error(ctx, op, logger.Error(err))
		writeGRPCError(w, err)
		return
	}

	httpResp, err := mapper.ToListCategoryHTTP(resp)
	if err != nil {
		h.l.Error(ctx, op, logger.Error(err))
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if err := writeJSON(w, http.StatusOK, httpResp); err != nil {
		h.l.Error(ctx, op, logger.Error(err))
	}
}
