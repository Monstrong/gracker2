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

func (c *AuthGrpcClient) Login(
	ctx context.Context,
	req *authv1.LoginRequest,
) (*authv1.LoginResponse, error) {
	return c.client.Login(ctx, req)
}

func (c *AuthGrpcClient) GetMe(
	ctx context.Context,
	req *authv1.GetMeRequest,
) (*authv1.GetMeResponse, error) {
	return c.client.GetMe(ctx, req)
}