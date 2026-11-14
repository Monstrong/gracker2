package clients

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authv1 "github.com/monstrong/gracker2/proto/gen/go/auth/v1" // при смене версии сменить версию
)

type AuthGrpcClient struct {
	conn   *grpc.ClientConn
	client authv1.AuthServiceClient
}

func (c *AuthGrpcClient) Close() error {
	return c.conn.Close()
}
func NewAuthClient(addr string) (*AuthGrpcClient, error) {

	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		return nil, err
	}

	client := authv1.NewAuthServiceClient(conn)

	return &AuthGrpcClient{conn: conn, client: client}, nil
}

func (c *AuthGrpcClient) Register(
	ctx context.Context,
	req *authv1.RegisterRequest,
) (*authv1.RegisterResponse, error) {
	return c.client.Register(ctx, req)
}

func (c *AuthGrpcClient) Login(
	ctx context.Context,
	req *authv1.LoginRequest,
) (*authv1.LoginResponse, error) {
	return c.client.Login(ctx, req)
}

func (c *AuthGrpcClient) Refresh(
	ctx context.Context,
	req *authv1.RefreshRequest,
) (*authv1.RefreshResponse, error) {
	return c.client.Refresh(ctx, req)
}

func (c *AuthGrpcClient) Logout(
	ctx context.Context,
	req *authv1.LogoutRequest,
) (*authv1.LogoutResponse, error) {
	return c.client.Logout(ctx, req)
}

func (c *AuthGrpcClient) LogoutEverywhere(
	ctx context.Context,
	req *authv1.LogoutEverywhereRequest,
) (*authv1.LogoutEverywhereResponse, error) {
	return c.client.LogoutEverywhere(ctx, req)
}

func (c *AuthGrpcClient) RoleUpdate(
	ctx context.Context,
	req *authv1.RoleUpdateRequest,
) (*authv1.RoleUpdateResponse, error) {
	return c.client.RoleUpdate(ctx, req)
}

func (c *AuthGrpcClient) GetMe(
	ctx context.Context,
	req *authv1.GetMeRequest,
) (*authv1.GetMeResponse, error) {
	return c.client.GetMe(ctx, req)
}

func (c *AuthGrpcClient) GetUser(
	ctx context.Context,
	req *authv1.GetUserRequest,
) (*authv1.GetUserResponse, error) {
	return c.client.GetUser(ctx, req)
}

func (c *AuthGrpcClient) AdminGetUser(
	ctx context.Context,
	req *authv1.AdminGetUserRequest,
) (*authv1.AdminGetUserResponse, error) {
	return c.client.AdminGetUser(ctx, req)
}

func (c *AuthGrpcClient) UpdateUser(
	ctx context.Context,
	req *authv1.UpdateUserRequest,
) (*authv1.UpdateUserResponse, error) {
	return c.client.UpdateUser(ctx, req)
}

func (c *AuthGrpcClient) DeleteUser(
	ctx context.Context,
	req *authv1.DeleteUserRequest,
) (*authv1.DeleteUserResponse, error) {
	return c.client.DeleteUser(ctx, req)
}

func (c *AuthGrpcClient) BanUser(
	ctx context.Context,
	req *authv1.BanUserRequest,
) (*authv1.BanUserResponse, error) {
	return c.client.BanUser(ctx, req)
}

func (c *AuthGrpcClient) UnBanUser(
	ctx context.Context,
	req *authv1.UnBanUserRequest,
) (*authv1.UnBanUserResponse, error) {
	return c.client.UnBanUser(ctx, req)
}

func (c *AuthGrpcClient) ChangePassword(
	ctx context.Context,
	req *authv1.ChangePasswordRequest,
) (*authv1.ChangePasswordResponse, error) {
	return c.client.ChangePassword(ctx, req)
}
