package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
)



func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {
	const op = "service.User.Login"
	return "", nil
}

func (s *UserService) Logout(ctx context.Context, token string) (bool, error) {
	const op = "service.User.Logout"
	return false, nil
}

func (s *UserService) Refresh(ctx context.Context, token string) (string, error) {
	const op = "service.User.Refresh"
	return "", nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}