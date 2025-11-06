package auth

import (
	"context"
	"fmt"

	authpb "auth2/api/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	authpb.UnimplementedAuthServiceServer
	svc *Service
}

func NewAuthHandler(svc *Service) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	id, err := h.svc.Register(req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &authpb.RegisterResponse{Token: fmt.Sprintf("user_%d", id)}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	token, err := h.svc.Login(req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, err.Error())
	}
	return &authpb.LoginResponse{Token: token}, nil
}
