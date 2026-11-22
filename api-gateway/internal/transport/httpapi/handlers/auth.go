package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
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

// Swagger комментарии были сгенерены нейронкой, да

// Register godoc
// @Summary Register a new user
// @Description Creates a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration data"
// @Success 201 {object} dto.RegisterResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.Register"

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.l.Info(ctx, "user sent invalid request", logger.String("op", op))
		writeError(w, http.StatusBadRequest, "bad request")
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

	if err := writeJSON(w, http.StatusCreated, httpResp); err != nil {
		h.l.Error(ctx, op, logger.Error(err))
	}
}

// Login godoc
// @Summary Login user
// @Description Authenticates a user and returns access and refresh tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.Login"

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.l.Info(ctx, "user sent invalid request", logger.String("op", op))
		writeError(w, http.StatusBadRequest, "bad request")
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

// Refresh godoc
// @Summary Refresh authentication tokens
// @Description Generates a new access and refresh token pair
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshRequest true "Refresh token"
// @Success 200 {object} dto.RefreshResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.Refresh"

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	var req dto.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.l.Info(ctx, "user sent invalid request", logger.String("op", op))
		writeError(w, http.StatusBadRequest, "bad request")
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

// Logout godoc
// @Summary Logout user
// @Description Revokes the provided refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LogoutRequest true "Refresh token"
// @Success 200
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.Logout"

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	var req dto.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.l.Info(ctx, "user sent invalid request", logger.String("op", op))
		writeError(w, http.StatusBadRequest, "bad request")
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

// LogoutEverywhere godoc
// @Summary Logout from all devices
// @Description Revokes all refresh tokens belonging to the authenticated user
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/logout-all [post]
func (h *AuthHandler) LogoutEverywhere(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.LogoutEverywhere"

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	claims, ok := middleware.GetClaims(ctx)
	if !ok {
		h.l.Info(ctx, "user sent invalid claims", logger.String("op", op))
		writeError(w, http.StatusUnauthorized, "unauthorized")
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

// RoleUpdate godoc
// @Summary Update user role
// @Description Changes the role of a target user. Requires administrator permissions
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.RoleUpdateRequest true "Role update data"
// @Success 200 {object} dto.RoleUpdateResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/users/role [patch]
func (h *AuthHandler) RoleUpdate(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.RoleUpdate"

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	claims, ok := middleware.GetClaims(ctx)
	if !ok {
		h.l.Info(ctx, "user sent invalid claims", logger.String("op", op))
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	ctx = metadata.AppendToOutgoingContext(
		ctx,
		grpcmetakey.GRPCMetaKeyRole, strconv.Itoa(int(claims.Role)),
	)

	var req dto.RoleUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.l.Info(ctx, "user sent invalid request", logger.String("op", op))
		writeError(w, http.StatusBadRequest, "bad request")
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

// GetMe godoc
// @Summary Get current user
// @Description Returns information about the authenticated user
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.GetMeResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/me [get]
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.GetMe"

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		h.l.Info(ctx, "user sent invalid claims", logger.String("op", op))
		writeError(w, http.StatusUnauthorized, "unauthorized")
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
		writeGRPCError(w, err)
		return
	}

	if err := writeJSON(w, http.StatusOK, httpResp); err != nil {
		h.l.Error(ctx, op, logger.Error(err))
	}
}

// GetUser godoc
// @Summary Get public user information
// @Description Returns public information about a user by UUID
// @Tags Users
// @Produce json
// @Param id path string true "User UUID"
// @Success 200 {object} dto.GetUserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/{id} [get]
func (h *AuthHandler) GetUser(w http.ResponseWriter, r *http.Request, id string) {
	const op = "handlers.auth.GetUser"

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	uuidId, err := uuid.Parse(id)
	if err != nil {
		h.l.Info(ctx, "user sent invalid id", logger.String("op", op))
		writeError(w, http.StatusBadRequest, "id must be uuid")
		return
	}

	req := dto.GetUserRequest{
		UserId: uuidId,
	}


	grpcReq, err := mapper.ToGetUserGRPC(req)
	if err != nil {
		h.l.Info(ctx, "user sent invalid request", logger.String("op", op))
		writeGRPCError(w, err)
		return
	}

	resp, err := h.authClient.GetUser(ctx, grpcReq)
	if err != nil {
		h.l.Error(ctx, op, logger.Error(err))
		writeGRPCError(w, err)
		return
	}

	httpResp, err := mapper.ToGetUserHTTP(resp)
	if err != nil {
		h.l.Error(ctx, op, logger.Error(err))
		writeGRPCError(w, err)
		return 
	}

	if err := writeJSON(w, http.StatusOK, httpResp); err != nil {
		h.l.Error(ctx, op, logger.Error(err))
	}
}

// AdminGetUser godoc
// @Summary Get user information as administrator
// @Description Returns full user information by UUID. Requires administrator permissions
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param id path string true "User UUID"
// @Success 200 {object} dto.AdminGetUserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/users/{id} [get]
func (h *AuthHandler) AdminGetUser(w http.ResponseWriter, r *http.Request, id string) {
	const op = "handlers.auth.AdminGetUser"

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	uuidId, err := uuid.Parse(id)
	if err != nil {
		h.l.Info(ctx, "user sent invalid id", logger.String("op", op))
		writeError(w, http.StatusBadRequest, "id must be uuid")
		return
	}

	req := dto.AdminGetUserRequest{
		Id: uuidId,
	}

	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		h.l.Info(ctx, "user sent invalid claims", logger.String("op", op))
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	ctx = metadata.AppendToOutgoingContext(
		ctx,
		grpcmetakey.GRPCMetaKeyRole, strconv.Itoa(int(claims.Role)),
	)

	grpcReq, err := mapper.ToAdminGetUserGRPC(req)
	if err != nil {
		h.l.Info(ctx, "user sent invalid request", logger.String("op", op))
		writeGRPCError(w, err)
		return
	}

	resp, err := h.authClient.AdminGetUser(ctx, grpcReq)
	if err != nil {
		h.l.Error(ctx, op, logger.Error(err))
		writeGRPCError(w, err)
		return
	}

	httpResp, err := mapper.ToAdminGetUserHTTP(resp)
	if err != nil {
		h.l.Error(ctx, op, logger.Error(err))
		writeGRPCError(w, err)
		return 
	}

	if err := writeJSON(w, http.StatusOK, httpResp); err != nil {
		h.l.Error(ctx, op, logger.Error(err))
	}
}

// UpdateUser godoc
// @Summary Update current user profile
// @Description Updates profile fields of the authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UpdateUserRequest true "Profile update data"
// @Success 200 {object} dto.UpdateUserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/me [patch]
func (h *AuthHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.UpdateUser"

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		h.l.Info(ctx, "user sent invalid claims", logger.String("op", op))
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	ctx = metadata.AppendToOutgoingContext(
		ctx,
		grpcmetakey.GRPCMetaKeyUserID, claims.UserID.String(),
	)

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.l.Info(ctx, "user sent invalid request", logger.String("op", op))
		writeError(w, http.StatusBadRequest, "bad request")
		return
	}

	grpcReq := mapper.ToUpdateUserGRPC(req)

	resp, err := h.authClient.UpdateUser(ctx, grpcReq)
	if err != nil {
		h.l.Error(ctx, op, logger.Error(err))
		writeGRPCError(w, err)
		return
	}

	httpResp, err := mapper.ToUpdateUserHTTP(resp)
	if err != nil {
		h.l.Error(ctx, op, logger.Error(err))
		writeGRPCError(w, err)
		return 
	}

	if err := writeJSON(w, http.StatusOK, httpResp); err != nil {
		h.l.Error(ctx, op, logger.Error(err))
	}
}

// DeleteUser godoc
// @Summary Delete user
// @Description Deletes the authenticated user account
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.DeleteUserRequest true "User deletion data"
// @Success 200 {object} dto.DeleteUserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users [delete]
func (h *AuthHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.DeleteUser"

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		h.l.Info(ctx, "user sent invalid claims", logger.String("op", op))
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	ctx = metadata.AppendToOutgoingContext(
		ctx,
		grpcmetakey.GRPCMetaKeyRole, strconv.Itoa(int(claims.Role)),
		grpcmetakey.GRPCMetaKeyUserID, claims.UserID.String(),
	)

	var req dto.DeleteUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.l.Info(ctx, "user sent invalid request", logger.String("op", op))
		writeError(w, http.StatusBadRequest, "bad request")
		return
	}

	grpcReq, err := mapper.ToDeleteUserGRPC(req)
	if err != nil {
		h.l.Info(ctx, "user sent invalid request", logger.String("op", op))
		writeGRPCError(w, err)
		return
	}

	resp, err := h.authClient.DeleteUser(ctx, grpcReq)
	if err != nil {
		h.l.Error(ctx, op, logger.Error(err))
		writeGRPCError(w, err)
		return
	}

	httpResp := mapper.ToDeleteUserHTTP(resp)
	

	if err := writeJSON(w, http.StatusOK, httpResp); err != nil {
		h.l.Error(ctx, op, logger.Error(err))
	}
}

// BanUser godoc
// @Summary Ban user
// @Description Bans a target user. Requires administrator permissions
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.BanUserRequest true "User ban data"
// @Success 200 {object} dto.BanUserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/users/ban [post]
func (h *AuthHandler) BanUser(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.BanUser"

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		h.l.Info(ctx, "user sent invalid claims", logger.String("op", op))
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	ctx = metadata.AppendToOutgoingContext(
		ctx,
		grpcmetakey.GRPCMetaKeyRole, strconv.Itoa(int(claims.Role)),
	)

	var req dto.BanUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.l.Info(ctx, "user sent invalid request", logger.String("op", op))
		writeError(w, http.StatusBadRequest, "bad request")
		return
	}

	grpcReq, err := mapper.ToBanUserGRPC(req)
	if err != nil {
		h.l.Info(ctx, "user sent invalid request", logger.String("op", op))
		writeGRPCError(w, err)
		return
	}

	resp, err := h.authClient.BanUser(ctx, grpcReq)
	if err != nil {
		h.l.Error(ctx, op, logger.Error(err))
		writeGRPCError(w, err)
		return
	}

	httpResp := mapper.ToBanUserHTTP(resp)

	if err := writeJSON(w, http.StatusOK, httpResp); err != nil {
		h.l.Error(ctx, op, logger.Error(err))
	}
}

// UnbanUser godoc
// @Summary Unban user
// @Description Removes a ban from a target user. Requires administrator permissions
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.UnBanUserRequest true "User unban data"
// @Success 200 {object} dto.UnBanUserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/users/unban [post]
func (h *AuthHandler) UnbanUser(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.UnbanUser"

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		h.l.Info(ctx, "user sent invalid claims", logger.String("op", op))
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	ctx = metadata.AppendToOutgoingContext(
		ctx,
		grpcmetakey.GRPCMetaKeyRole, strconv.Itoa(int(claims.Role)),

	)

	var req dto.UnBanUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.l.Info(ctx, "user sent invalid request", logger.String("op", op))
		writeError(w, http.StatusBadRequest, "bad request")
		return
	}

	grpcReq, err := mapper.ToUnBanUserGRPC(req)
	if err != nil {
		h.l.Info(ctx, "invalid id", logger.String("op", op))
		writeGRPCError(w, err)
		return
	}

	resp, err := h.authClient.UnBanUser(ctx, grpcReq)
	if err != nil {
		h.l.Error(ctx, op, logger.Error(err))
		writeGRPCError(w, err)
		return
	}

	httpResp := mapper.ToUnBanUserHTTP(resp)

	if err := writeJSON(w, http.StatusOK, httpResp); err != nil {
		h.l.Error(ctx, op, logger.Error(err))
	}
}

// ChangePassword godoc
// @Summary Change password
// @Description Changes the password of the authenticated user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ChangePasswordRequest true "Old and new password"
// @Success 200 {object} dto.ChangePasswordResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/me/password [patch]
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.auth.ChangePassword"

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.Timeout)
	defer cancel()

	claims, ok := middleware.GetClaims(r.Context())
	if !ok {
		h.l.Info(ctx, "user sent invalid claims", logger.String("op", op))
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	ctx = metadata.AppendToOutgoingContext(
		ctx,
		grpcmetakey.GRPCMetaKeyUserID, claims.UserID.String(),
	)

	var req dto.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.l.Info(ctx, "user sent invalid request", logger.String("op", op))
		writeError(w, http.StatusBadRequest, "bad request")
		return
	}

	grpcReq := mapper.ToChangePasswordGRPC(req)

	resp, err := h.authClient.ChangePassword(ctx, grpcReq)
	if err != nil {
		h.l.Error(ctx, op, logger.Error(err))
		writeGRPCError(w, err)
		return
	}

	httpResp := mapper.ToChangePasswordHTTP(resp)

	if err := writeJSON(w, http.StatusOK, httpResp); err != nil {
		h.l.Error(ctx, op, logger.Error(err))
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
    _ = writeJSON(w, status, dto.ErrorResponse{
        Error: message,
    })
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
		writeError(w, http.StatusInternalServerError, "internal server error")
		return 
		// ok == false возвращается только когда ошибка вообще не gRPC-статус 
		// сюда попадут всякие ошибки маппера и тд
	}
	
	switch st.Code() {
    case codes.InvalidArgument:
        writeError(w, http.StatusBadRequest, st.Message())
    case codes.Unauthenticated:
        writeError(w, http.StatusUnauthorized, st.Message())
    case codes.PermissionDenied:
        writeError(w, http.StatusForbidden, st.Message())
    case codes.NotFound:
        writeError(w, http.StatusNotFound, st.Message())
    case codes.AlreadyExists:
        writeError(w, http.StatusConflict, st.Message())
    case codes.DeadlineExceeded:
        writeError(w, http.StatusGatewayTimeout, st.Message())
    case codes.Unavailable:
        writeError(w, http.StatusServiceUnavailable, st.Message())
    default:
        writeError(w, http.StatusInternalServerError, "internal server error")
    }
}