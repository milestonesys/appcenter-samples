package repositories

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"apigateway-webserver/src/pkg/constants"
	"apigateway-webserver/src/pkg/constants/enums"
	"apigateway-webserver/src/pkg/entities/vms"
	"apigateway-webserver/src/pkg/repositories/base"
)

// Interface for implementing the gr api repository
type GatewayRepository interface {
	// Given a copy of a management server Query the api gr well known Uri. Check constants.ApiWellKnownUris for more information
	RequestGatewayWellKnownUris(ctx context.Context, s vms.Server) (*vms.ApiWellKnownUrisSchema, error)
	// Query all cameras related to a given hardware
	RequestEnabledCameras(ctx context.Context, s vms.Server, t vms.Token) (*vms.CamerasList, error)
	// Query all analytic events types
	RequestAnalyticEventTypes(ctx context.Context, s vms.Server, t vms.Token) (*vms.AnalyticEventTypes, error)
}

type gatewayRepository struct {
	base.HttpBaseRepository
}

func NewGatewayRepository() GatewayRepository {
	return &gatewayRepository{
		HttpBaseRepository: base.NewHttpBaseRepository(),
	}
}

func (gr gatewayRepository) RequestGatewayWellKnownUris(ctx context.Context, s vms.Server) (*vms.ApiWellKnownUrisSchema, error) {
	// Build the request url from the management server url
	requestUrl := s.ServerInputInfo().ServerURL

	// Build the request path
	requestUrl.Path = constants.ApiWellKnownUris

	// Execute Get request
	response, _, err := gr.DoFromArgs(ctx, http.MethodGet, requestUrl, nil, nil, enums.None)
	if err != nil {
		return nil, err
	}

	// Build the response data structure
	d := new(vms.ApiWellKnownUrisSchema)
	if err := json.Unmarshal(response, d); err != nil {
		return nil, err
	}

	return d, nil
}

func (gr gatewayRepository) RequestEnabledCameras(ctx context.Context, s vms.Server, t vms.Token) (*vms.CamerasList, error) {
	// Build the request url
	requestUrl, err := url.ParseRequestURI(s.ApiWellKnownUris.ApiGateways[0])
	if err != nil {
		return nil, err
	}

	// Build the request path
	requestUrl.Path = constants.EnabledCameras

	// Execute Get request
	response, _, err := gr.DoFromArgs(ctx, http.MethodGet, requestUrl, t, nil, enums.None)
	if err != nil {
		return nil, err
	}

	// Build the response data structure
	type data struct {
		Array []vms.Camera `json:"array"`
	}
	d := new(data)
	if err := json.Unmarshal(response, d); err != nil {
		return nil, err
	}

	cameras := vms.NewCamerasList()
	for _, camera := range d.Array {
		cameras.Add(&camera)
	}

	return cameras, nil
}

func (gr gatewayRepository) RequestAnalyticEventTypes(ctx context.Context, s vms.Server, t vms.Token) (*vms.AnalyticEventTypes, error) {
	// Build the request url
	requestUrl, err := url.ParseRequestURI(s.ApiWellKnownUris.ApiGateways[0])
	if err != nil {
		return nil, err
	}

	// Build the request path
	requestUrl.Path = constants.AnalyticEventTypes

	// Execute Get request
	response, _, err := gr.DoFromArgs(ctx, http.MethodGet, requestUrl, t, nil, enums.None)
	if err != nil {
		return nil, err
	}

	type data struct {
		Array []vms.AnalyticEventType `json:"array"`
	}
	d := new(data)
	if err := json.Unmarshal(response, d); err != nil {
		return nil, err
	}

	analyticEvents := vms.NewAnalyticEventTypes()
	for _, analyticEventType := range d.Array {
		analyticEvents.Add(analyticEventType)
	}

	return analyticEvents, nil
}
