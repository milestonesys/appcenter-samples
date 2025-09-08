package services

import (
	"context"

	"apigateway-webserver/src/pkg/entities/vms"
	"apigateway-webserver/src/pkg/repositories"
)

// Defines the interface for interacting with the API gateway.
type GatewayService interface {
	// Queries the API gateway for well-known URIs.
	// Check constants.ApiWellKnownUris for more information.
	RequestGatewayWellKnownUris(ctx context.Context, s *vms.Server) (*vms.ApiWellKnownUrisSchema, error)
	// Queries all cameras related to a given hardware.
	RequestEnabledCameras(ctx context.Context, s *vms.Server, t vms.Token) (*vms.CamerasList, error)
	// Queries all analytic event types.
	RequestAnalyticEventTypes(ctx context.Context, s *vms.Server, t vms.Token) (*vms.AnalyticEventTypes, error)
}

type gatewayService struct {
	gr repositories.GatewayRepository
}

// Creates a new instance of GatewayService.
func NewGatewayService() GatewayService {
	return &gatewayService{
		gr: repositories.NewGatewayRepository(),
	}
}

func (gs *gatewayService) RequestGatewayWellKnownUris(ctx context.Context, s *vms.Server) (*vms.ApiWellKnownUrisSchema, error) {
	return gs.gr.RequestGatewayWellKnownUris(ctx, *s)
}

func (gs *gatewayService) RequestEnabledCameras(ctx context.Context, s *vms.Server, t vms.Token) (*vms.CamerasList, error) {
	return gs.gr.RequestEnabledCameras(ctx, *s, t)
}

func (gs *gatewayService) RequestAnalyticEventTypes(ctx context.Context, s *vms.Server, t vms.Token) (*vms.AnalyticEventTypes, error) {
	return gs.gr.RequestAnalyticEventTypes(ctx, *s, t)
}
