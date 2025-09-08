package vms

import (
	"net/url"
	"strings"
)

type ApiWellKnownUrisSchema struct {
	ProductVersion           string   `json:"ProductVersion"`
	UnsecureManagementServer string   `json:"UnsecureManagementServer"`
	SecureManagementServer   string   `json:"SecureManagementServer"`
	IdentityProvider         string   `json:"IdentityProvider"`
	ApiGateways              []string `json:"ApiGateways"`
}

type IdpOpenIdConfigSchema struct {
	Issuer        string `json:"issuer"`
	TokenEndPoint string `json:"token_endpoint"`
	ServerVersion string `json:"server_version"`
}

type serverInputInfo struct {
	ServerURL *url.URL
}

type Server struct {
	serverInputInfo  serverInputInfo
	IdpOpenIdConfig  *IdpOpenIdConfigSchema
	ApiWellKnownUris *ApiWellKnownUrisSchema
}

func NewServer(serverURL *url.URL) *Server {
	return &Server{
		serverInputInfo:  serverInputInfo{ServerURL: serverURL},
		IdpOpenIdConfig:  &IdpOpenIdConfigSchema{},
		ApiWellKnownUris: &ApiWellKnownUrisSchema{},
	}
}

func (s *Server) ServerInputInfo() serverInputInfo {
	return s.serverInputInfo
}

func (s *Server) Hostname() string {
	return s.serverInputInfo.ServerURL.Host
}

func (s *Server) IsSecure() bool {
	return strings.HasPrefix(s.serverInputInfo.ServerURL.Scheme, "https")
}
