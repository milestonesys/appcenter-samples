package services

import (
	"context"

	"apigateway-webserver/src/pkg/entities/events"
	"apigateway-webserver/src/pkg/entities/vms"
	"apigateway-webserver/src/pkg/repositories"
)

// Interface for implementing the gateway api repository
type WsEventsService interface {
	// 1- start new session
	// To resume a previous session (within 30 seconds), session_id and event_id of the last received event must be provided.
	// This will resume any subscriptions previously made and recover lost events.
	// If the session cannot be resumed (e.g. due to timeout) a new session will be created, and any subscriptions must be created again.
	// response["status"] will contain the status
	// - 200 indicates an existing session was successfully resumed.
	// - 201 indicates a new session was created.
	RequestStartSession(ctx context.Context, s *vms.Server, t vms.Token) (*events.WsCommandResponse, error)

	// 2- subscribe to topic
	RequestSubscribe(ctx context.Context, cameraId string, eventTypeId string) (*events.WsCommandResponse, error)

	// 3- Subscribe to topic and loop
	RequestEvents(ctx context.Context) (*events.AnalyticsEvents, error)

	// 4- Close communication
	RequestClose() error
}

type wsEventsService struct {
	wer repositories.WsEventsRepository
}

func NewWsEventsService() WsEventsService {
	return &wsEventsService{
		wer: repositories.NewWsEventsRepository(),
	}
}

func (wes *wsEventsService) RequestStartSession(ctx context.Context, s *vms.Server, t vms.Token) (*events.WsCommandResponse, error) {
	return wes.wer.RequestStartSession(ctx, *s, t)
}

func (wes *wsEventsService) RequestSubscribe(ctx context.Context, cameraId string, eventTypeId string) (*events.WsCommandResponse, error) {
	return wes.wer.RequestSubscribe(ctx, cameraId, eventTypeId)
}

func (wes *wsEventsService) RequestEvents(ctx context.Context) (*events.AnalyticsEvents, error) {
	return wes.wer.RequestEvents(ctx)
}

func (wes *wsEventsService) RequestClose() error {
	return wes.wer.RequestClose()
}
