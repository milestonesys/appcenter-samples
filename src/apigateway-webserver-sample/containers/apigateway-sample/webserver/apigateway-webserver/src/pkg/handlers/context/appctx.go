package context

import (
	"apigateway-webserver/src/pkg/entities/events"
	"apigateway-webserver/src/pkg/entities/vms"
	"apigateway-webserver/src/pkg/services"
)

type AppContext interface {
	IdpService() services.IdpService
	GatewayService() services.GatewayService
	WsEventsService() services.WsEventsService

	Server() *vms.Server
	User() *vms.User
	Token() vms.Token

	SetWsCommandResponse(wsCommandResponse *events.WsCommandResponse)
	GetWsCommandResponse() *events.WsCommandResponse
}

type appContext struct {
	idpService      services.IdpService
	gatewayService  services.GatewayService
	wsEventsService services.WsEventsService

	server *vms.Server
	user   *vms.User
	token  vms.Token

	wsCommandResponse *events.WsCommandResponse
}

func NewAppContext(
	idpService services.IdpService,
	gatewayService services.GatewayService,
	wsEventsService services.WsEventsService,
	server *vms.Server,
	user *vms.User,
	token vms.Token) AppContext {
	return &appContext{
		idpService:      idpService,
		gatewayService:  gatewayService,
		wsEventsService: wsEventsService,
		server:          server,
		user:            user,
		token:           token,
	}
}

func (a *appContext) IdpService() services.IdpService {
	return a.idpService
}

func (a *appContext) GatewayService() services.GatewayService {
	return a.gatewayService
}

func (a *appContext) WsEventsService() services.WsEventsService {
	return a.wsEventsService
}

func (a *appContext) Server() *vms.Server {
	return a.server
}

func (a *appContext) User() *vms.User {
	return a.user
}

func (a *appContext) Token() vms.Token {
	return a.token
}

func (a *appContext) SetWsCommandResponse(wsCommandResponse *events.WsCommandResponse) {
	a.wsCommandResponse = wsCommandResponse
}

func (a *appContext) GetWsCommandResponse() *events.WsCommandResponse {
	if a.wsCommandResponse == nil {
		return &events.WsCommandResponse{}
	}
	return a.wsCommandResponse
}
