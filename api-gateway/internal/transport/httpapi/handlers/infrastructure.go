package handlers

import (
	"context"
	
	authv1 "github.com/monstrong/gracker2/proto/gen/go/auth/v1"
)


type AuthClient interface {
    Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error)
    GetMe(ctx context.Context, req *authv1.GetMeRequest) (*authv1.GetMeResponse, error)
}