package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/monstrong/gracker2/api-gateway/internal/gatewayconfig"
	"github.com/monstrong/gracker2/api-gateway/internal/middleware"
	"github.com/monstrong/gracker2/api-gateway/internal/transport/httpapi/dto"
	"github.com/monstrong/gracker2/api-gateway/internal/transport/httpapi/mapper"
	"github.com/monstrong/gracker2/pkg/contracts/grpcmetakey"
	"github.com/monstrong/gracker2/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	authv1 "github.com/monstrong/gracker2/proto/gen/go/auth/v1"
)

type AuthHandler struct {
    authClient AuthClient
	cfg gatewayconfig.Auth
	l logger.Logger
}


func NewAuthHandler(authClient AuthClient, cfg gatewayconfig.Auth, l logger.Logger) *AuthHandler {
	return &AuthHandler{authClient: authClient, cfg: cfg, l: l}
}

// 	пример того как надо передавать роль или id в metadata
// ctx = metadata.AppendToOutgoingContext(
//     ctx,
//     grpcmetakey.GRPCMetaKeyUserID, claims.UserID.String(),
//     grpcmetakey.GRPCMetaKeyRole, strconv.Itoa(int(claims.Role)),
// )

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.Register"

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.l.Info(ctx, "user sent invalid request", logger.String("op", op))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	grpcReq := mapper.ToRegisterGRPC(req)

	resp, err := h.authClient.Register(ctx, grpcReq)
	if err != nil {
		h.l.Error(ctx, op, logger.Error(err))
		writeGRPCError(w, err)
		return
	}

	httpResp, err := mapper.ToRegisterHTTP(resp)
	if err != nil {
		h.l.Error(ctx, op, logger.Error(err))
		writeGRPCError(w, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, httpResp); err != nil {
		h.l.Error(ctx, op, logger.Error(err))
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.Login"

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.l.Info(ctx, "user sent invalid request", logger.String("op", op))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	grpcReq := mapper.ToLoginGRPC(req)

    resp, err := h.authClient.Login(ctx, grpcReq)
    if err != nil {
		h.l.Error(ctx, op, logger.Error(err))
        writeGRPCError(w, err)
        return
    }

    httpResp := mapper.ToLoginHTTP(resp)

    if err := writeJSON(w, http.StatusOK, httpResp); err != nil {
		h.l.Error(ctx, op, logger.Error(err))
	}
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.Refresh"

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	var req dto.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.l.Info(ctx, "user sent invalid request", logger.String("op", op))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	grpcReq := mapper.ToRefreshGRPC(req)

	resp, err := h.authClient.Refresh(ctx, grpcReq)
	if err != nil {
		h.l.Error(ctx, op, logger.Error(err))
		writeGRPCError(w, err)
		return
	}

	httpResp := mapper.ToRefreshHTTP(resp)

	if err := writeJSON(w, http.StatusOK, httpResp); err != nil {
		h.l.Error(ctx, op, logger.Error(err))
	}
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.Logout"

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	var req dto.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.l.Info(ctx, "user sent invalid request", logger.String("op", op))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	grpcReq := mapper.ToLogoutGRPC(req)

	_, err := h.authClient.Logout(ctx, grpcReq)
	if err != nil {
		h.l.Error(ctx, op, logger.Error(err))
		writeGRPCError(w, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, struct{}{}); err != nil {
		h.l.Error(ctx, op, logger.Error(err))
	}
}

func (h *AuthHandler) LogoutEverywhere(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.LogoutEverywhere"

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	claims, ok := middleware.GetClaims(ctx)
	if !ok {
		h.l.Info(ctx, "user sent invalid claims", logger.String("op", op))
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	ctx = metadata.AppendToOutgoingContext(
		ctx,
		grpcmetakey.GRPCMetaKeyUserID, claims.UserID.String(),
	)

	var req *authv1.LogoutEverywhereRequest
	_, err := h.authClient.LogoutEverywhere(ctx, req)
	if err != nil {
		h.l.Error(ctx, op, logger.Error(err))
		writeGRPCError(w, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, struct{}{}); err != nil {
		h.l.Error(ctx, op, logger.Error(err))
	}
}

func (h *AuthHandler) RoleUpdate(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.RoleUpdate"

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	claims, ok := middleware.GetClaims(ctx)
	if !ok {
		h.l.Info(ctx, "user sent invalid claims", logger.String("op", op))
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	ctx = metadata.AppendToOutgoingContext(
		ctx,
		grpcmetakey.GRPCMetaKeyRole, strconv.Itoa(int(claims.Role)),
	)

	var req dto.RoleUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.l.Info(ctx, "user sent invalid request", logger.String("op", op))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	grpcReq := mapper.ToRoleUpdateGRPC(req)

	resp, err := h.authClient.RoleUpdate(ctx, grpcReq)
	if err != nil {
		h.l.Error(ctx, op, logger.Error(err))
		writeGRPCError(w, err)
		return
	}

	httpResp := mapper.ToRoleUpdateHTTP(resp)

	if err := writeJSON(w, http.StatusOK, httpResp); err != nil {
		h.l.Error(ctx, op, logger.Error(err))
	}
}
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.GetMe"

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		h.l.Info(ctx, "user sent invalid claims", logger.String("op", op))
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	ctx = metadata.AppendToOutgoingContext(
        ctx,
        grpcmetakey.GRPCMetaKeyUserID, claims.UserID.String(),
    )

	resp, err := h.authClient.GetMe(ctx, &authv1.GetMeRequest{})
	if err != nil {
		h.l.Error(ctx, op, logger.Error(err))
		writeGRPCError(w, err)
		return
	}

	httpResp, err := mapper.ToGetMeHTTP(resp)
	if err != nil {
		h.l.Error(ctx, op, logger.Error(err))
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := writeJSON(w, http.StatusOK, httpResp); err != nil {
		h.l.Error(ctx, op, logger.Error(err))
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// возвращает пользователю весь message ошибки.
// Нельзя возвращать сюда из сервисов сырые ошибки
func writeGRPCError(w http.ResponseWriter, err error) {
	st, ok := status.FromError(err)
	if !ok {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return 
		// ok == false возвращается только когда ошибка вообще не gRPC-статус 
		// сюда попадут всякие ошибки маппера и тд
	}
	
	switch st.Code() {
    case codes.InvalidArgument:
        http.Error(w, st.Message(), http.StatusBadRequest)
    case codes.Unauthenticated:
        http.Error(w, st.Message(), http.StatusUnauthorized)
    case codes.PermissionDenied:
        http.Error(w, st.Message(), http.StatusForbidden)
    case codes.NotFound:
        http.Error(w, st.Message(), http.StatusNotFound)
    case codes.AlreadyExists:
        http.Error(w, st.Message(), http.StatusConflict)
    case codes.DeadlineExceeded:
        http.Error(w, st.Message(), http.StatusGatewayTimeout)
    case codes.Unavailable:
        http.Error(w, st.Message(), http.StatusServiceUnavailable)
    default:
        http.Error(w, "internal server error", http.StatusInternalServerError)
    }
}