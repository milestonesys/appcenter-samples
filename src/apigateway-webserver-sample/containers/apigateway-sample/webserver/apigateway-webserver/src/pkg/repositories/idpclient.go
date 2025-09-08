package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"apigateway-webserver/src/pkg/constants"
	"apigateway-webserver/src/pkg/constants/enums"
	"apigateway-webserver/src/pkg/entities/vms"
	"apigateway-webserver/src/pkg/repositories/base"
)

type IdpRepository interface {
	// Queries the IDP API well-known configuration (OpenId Configuration)
	RequestIdpWellKnownConfig(ctx context.Context, s vms.Server) (*vms.IdpOpenIdConfigSchema, error)

	// Sends a post request to get an access token for a "basic user" for the management server scope (not supported for windows users by design)
	RequestAccessToken(ctx context.Context, u vms.User, s vms.Server, td TokenDispatcher) (vms.Token, error)
}

type idpRepository struct {
	base.HttpBaseRepository
}

func NewIdpRepository() IdpRepository {
	return &idpRepository{
		HttpBaseRepository: base.NewHttpBaseRepository(),
	}
}

func (ir idpRepository) RequestIdpWellKnownConfig(ctx context.Context, s vms.Server) (*vms.IdpOpenIdConfigSchema, error) {
	serverUrl := s.ServerInputInfo().ServerURL
	// Configure request url path (so far the request url was the server url. Now, we add the api server endpoint)
	serverUrl.Path = constants.IdpWellKnownOpenIdConfig

	// Execute GET request
	response, _, err := ir.DoFromArgs(ctx, http.MethodGet, serverUrl, nil, nil, enums.None)
	if err != nil {
		return nil, fmt.Errorf("failed to execute GET request: %w", err)
	}

	var config vms.IdpOpenIdConfigSchema
	if err := json.Unmarshal(response, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &config, nil
}

func (ir idpRepository) RequestAccessToken(ctx context.Context, u vms.User, s vms.Server, td TokenDispatcher) (vms.Token, error) {
	// Build the request url
	requestUrl, err := url.ParseRequestURI(s.IdpOpenIdConfig.TokenEndPoint)
	if err != nil {
		return nil, fmt.Errorf("invalid token endpoint URL: %w", err)
	}

	// Create basic user token request
	payload := url.Values{}
	if u.CredentialsFlowType() == enums.ClientCredentialsFlow {
		payload.Set("grant_type", "client_credentials")
		payload.Set("scope", "managementserver")
		payload.Set("client_id", u.Username())     // At the ReadCredentialsFlowFiles the username is the client_id
		payload.Set("client_secret", u.Password()) // And the password contains the client_secret
	} else {
		payload.Set("grant_type", "password")
		payload.Set("username", u.Username())
		payload.Set("password", u.Password())
		payload.Set("client_id", "GrantValidatorClient")
	}

	// Execute request
	response, _, err := ir.DoFromArgs(ctx, http.MethodPost, requestUrl, nil, strings.NewReader(payload.Encode()), enums.Urlencoded)
	if err != nil {
		return nil, fmt.Errorf("failed to execute POST request: %w", err)
	}

	// Load response into the token and return copy of the modified token
	return vms.NewToken(response, td.DispatchFunc())
}
