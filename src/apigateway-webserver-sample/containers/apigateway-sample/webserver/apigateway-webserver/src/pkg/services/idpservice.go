package services

import (
	"context"

	"apigateway-webserver/src/pkg/entities/vms"
	"apigateway-webserver/src/pkg/repositories"
)

type IdpService interface {
	// Queries the IDP API well-known configuration (OpenID Configuration).
	RequestIdpWellKnownConfig(ctx context.Context, s *vms.Server) (*vms.IdpOpenIdConfigSchema, error)

	// Sends a POST request to get an access token for a basic user for management server scope (not supported for Windows users by design)..
	RequestAccessToken(ctx context.Context, u *vms.User, s *vms.Server) (vms.Token, error)
}

type idpService struct {
	ir repositories.IdpRepository
}

func NewIdpService() IdpService {
	return &idpService{
		ir: repositories.NewIdpRepository(),
	}
}

func (is *idpService) RequestIdpWellKnownConfig(ctx context.Context, s *vms.Server) (*vms.IdpOpenIdConfigSchema, error) {
	return is.ir.RequestIdpWellKnownConfig(ctx, *s)
}

func (is *idpService) RequestAccessToken(ctx context.Context, u *vms.User, s *vms.Server) (vms.Token, error) {
	// Create a token dispatcher (one per user and server combination).
	// The reference to this instance will be stored in the token object.
	// Every time the token value is requested, the token dispatcher will be called to check whether the token needs to be renewed or not.
	tokenDispatcher := repositories.NewTokenDispatcher(is.ir, u, s)

	// Sends a POST request to get the access token for the management server scope.
	return is.ir.RequestAccessToken(ctx, *u, *s, tokenDispatcher)
}
