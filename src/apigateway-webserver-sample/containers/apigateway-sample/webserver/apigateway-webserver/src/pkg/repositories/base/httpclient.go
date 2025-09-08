package base

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"

	"apigateway-webserver/src/pkg/constants/enums"
	"apigateway-webserver/src/pkg/entities/vms"
)

type HttpBaseRepository struct {
	BaseRepository
}

func NewHttpBaseRepository() HttpBaseRepository {
	return HttpBaseRepository{
		BaseRepository: NewBaseRepository(),
	}
}

// Sends a request using the HTTP client.
// It creates the request from the given arguments.
// - method: Specifies whether it's a GET, POST, DELETE, etc.
// - requestUrl: The URL to which the request is sent.
// - token: Optional parameter that gets added to the request header if provided.
// - body: Optional parameter that gets added as the request body if provided.
// - contentType: The content type of the request body.
func (hbr HttpBaseRepository) DoFromArgs(ctx context.Context, method string, requestUrl *url.URL, token vms.Token, body io.Reader, contentType enums.RequestContentType) ([]byte, int, error) {
	if body == nil {
		body = http.NoBody
	}

	// Create request from given arguments
	request, err := http.NewRequestWithContext(ctx, method, requestUrl.String(), body)
	if err != nil {
		return nil, -1, err
	}

	// Set the content type
	if err := hbr.setContentType(request, contentType); err != nil {
		return nil, -1, err
	}

	// Check if the token was provided and add it to the request header
	if token != nil {
		bearerToken, err := token.DispatchToken(ctx)
		if err != nil {
			return nil, -1, err
		}
		request.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Execute request
	return hbr.doFromRequest(request)
}

// Sets the content type in the request header based on the given content type enum.
func (hbr HttpBaseRepository) setContentType(request *http.Request, contentType enums.RequestContentType) error {
	switch contentType {
	case enums.None:
		request.Header.Del("Content-Type")
		request.Header.Set("Content-Length", "0")
	case enums.Urlencoded, enums.Raw, enums.Json:
		request.Header.Set("Content-Type", contentType.String())
	default:
		return errors.New("Body type not supported: " + contentType.String())
	}
	return nil
}

// Executes any HTTP request and returns the response as bytes.
func (hbr HttpBaseRepository) doFromRequest(request *http.Request) ([]byte, int, error) {
	// Set the request transport to support both encrypted and unencrypted communication
	hbr.setRequestTransport(request.URL)

	// Execute request
	resp, err := hbr.client.Do(request)
	if err != nil {
		return nil, -1, err
	}

	// Ensure the response body is closed before exiting the function
	defer resp.Body.Close()

	// Check the status code and return it as an error if it is not OK
	if resp.StatusCode >= 400 && resp.StatusCode <= 511 {
		return nil, resp.StatusCode, errors.New(resp.Status)
	}

	// Read the body and convert it to bytes
	// If the body is empty, return an empty array of bytes
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	// Success: return body content and status code
	return bytes, resp.StatusCode, nil
}
