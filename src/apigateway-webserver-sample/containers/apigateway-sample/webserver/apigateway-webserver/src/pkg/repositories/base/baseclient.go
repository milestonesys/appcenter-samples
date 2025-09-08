package base

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type BaseRepository struct {
	client *http.Client
}

func NewBaseRepository() BaseRepository {
	return BaseRepository{
		client: &http.Client{
			Timeout: 2 * time.Minute,
		},
	}
}

// Configures the request URL and sets a transport. This must be done once per request.
func (br *BaseRepository) setRequestTransport(requestUrl *url.URL) {
	br.client.Transport = &http.Transport{
		IdleConnTimeout: 30 * time.Second,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			dialer := net.Dialer{}
			return dialer.DialContext(ctx, network, addr)
		},
		DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			dialer := tls.Dialer{
				Config: &tls.Config{
					ServerName:         requestUrl.Hostname(),
					InsecureSkipVerify: strings.ToLower(requestUrl.Scheme) == "http",
				},
			}
			return dialer.DialContext(ctx, network, addr)
		},
	}
}
