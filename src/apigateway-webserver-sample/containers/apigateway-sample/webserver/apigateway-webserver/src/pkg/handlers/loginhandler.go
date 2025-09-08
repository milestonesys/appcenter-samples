package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"

	"apigateway-webserver/src/pkg/constants/enums"
	"apigateway-webserver/src/pkg/entities/appcenter"
	"apigateway-webserver/src/pkg/entities/vms"
	handlers_context "apigateway-webserver/src/pkg/handlers/context"
	"apigateway-webserver/src/pkg/services"
)

type LoginHandler struct {
	mu sync.Mutex
}

func NewLoginHandler() *LoginHandler {
	return &LoginHandler{}
}

func (lh *LoginHandler) Handle(w http.ResponseWriter, r *http.Request) {
	lh.mu.Lock()
	defer lh.mu.Unlock()
	log.Println("LoginHandler.Handle() called")

	var data struct {
		Username            string `json:"username"`
		Password            string `json:"password"`
		Hostname            string `json:"hostname"`
		Secure              bool   `json:"secure"`
		CredentialsFlowType string `json:"credentialsFlowType"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON format: %v", err), http.StatusBadRequest)
		return
	}

	if data.Hostname == "" {
		http.Error(w, "Missing required field: hostname", http.StatusBadRequest)
		return
	}

	// If the user selected the login form and didn't left the password field empty
	if data.CredentialsFlowType == enums.LoginForm.String() {
		if data.Password == "" {
			http.Error(w, "Missing required field: password", http.StatusBadRequest)
			return
		}
		if data.Hostname == "" {
			http.Error(w, "Missing required field: hostname", http.StatusBadRequest)
			return
		}
	}

	scheme := "http"
	if data.Secure {
		scheme = "https"
	}

	credentialsFlowType, err := enums.ParseCredentialsFlowType(data.CredentialsFlowType)
	if err != nil {
		http.Error(w, "Couldn't parse the provided credential flow type", http.StatusBadRequest)
		return
	}

	username, password, err := appcenter.ReadCredentialsFlowFiles(data.Username, data.Password, credentialsFlowType)
	if err != nil {
		http.Error(w, fmt.Sprintf("Couldn't read the credentials files: %v", err), http.StatusInternalServerError)
		return
	}

	appCtx, err := setupAppContext(data.Hostname, username, password, scheme, credentialsFlowType)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to perform login: %v", err), http.StatusInternalServerError)
		return
	}

	// Always override the previous user login session
	handlers_context.GetAppContextsInstance().AddAppContext(data.Username, appCtx)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{ "message": "Login success" }`))
}

func setupAppContext(hostname, username, password, scheme string, credentialsFlowType enums.CredentialsFlowType) (handlers_context.AppContext, error) {
	// Create services
	gatewayService := services.NewGatewayService()
	idpService := services.NewIdpService()
	wsEventsService := services.NewWsEventsService()

	// Parse the url string into a real url
	parsedServerUrl := &url.URL{
		Host:   hostname,
		Scheme: scheme,
	}

	// Create endpoint needed data structures
	user := vms.NewUser(username, password, credentialsFlowType)
	server := vms.NewServer(parsedServerUrl)

	var err error

	// Request gateway uris
	server.ApiWellKnownUris, err = gatewayService.RequestGatewayWellKnownUris(context.Background(), server)
	if err != nil {
		return nil, err
	}

	// Request idp openid config
	server.IdpOpenIdConfig, err = idpService.RequestIdpWellKnownConfig(context.Background(), server)
	if err != nil {
		return nil, err
	}

	// Create access token for the given management server and user
	token, err := idpService.RequestAccessToken(context.Background(), user, server)
	if err != nil {
		return nil, err
	}

	return handlers_context.NewAppContext(idpService, gatewayService, wsEventsService, server, user, token), nil
}
