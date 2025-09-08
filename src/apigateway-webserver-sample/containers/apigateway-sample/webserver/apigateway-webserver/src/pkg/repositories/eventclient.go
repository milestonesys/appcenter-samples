package repositories

import (
	"context"
	"errors"
	"net/url"
	"strconv"
	"strings"

	"apigateway-webserver/src/pkg/constants"
	"apigateway-webserver/src/pkg/entities/events"
	"apigateway-webserver/src/pkg/entities/vms"
	"apigateway-webserver/src/pkg/repositories/base"
)

var commandRequestsCounter int = 0

func newStartSessionRequest() *events.WsCommandRequest {
	return &events.WsCommandRequest{
		Command: "startSession",
	}
}

func newAddSubscriptionRequest(filters *events.SubscriptionFilters) *events.WsCommandRequest {
	return &events.WsCommandRequest{
		Command: "addSubscription",
		Filters: filters.Filters,
	}
}

func newSubscriptionFilter() *events.SubscriptionFilter {
	return &events.SubscriptionFilter{
		Modifier:      "include",
		ResourceTypes: []string{"cameras"},
		SourceIDs:     []string{},
		EventTypes:    []string{},
	}
}

func newSubscriptionFilters() *events.SubscriptionFilters {
	return &events.SubscriptionFilters{
		Filters: []events.SubscriptionFilter{*newSubscriptionFilter()},
	}
}

type WsEventsRepository interface {
	// 1- Start a new session
	// To resume a previous session (within 30 seconds), session_id and event_id of the last received event must be provided.
	// This will resume any subscriptions previously made and recover lost events.
	// If the session cannot be resumed (e.g., due to timeout), a new session will be created, and any subscriptions must be created again.
	// response["status"] will contain the status:
	// - 200 indicates an existing session was successfully resumed.
	// - 201 indicates a new session was created.
	RequestStartSession(ctx context.Context, s vms.Server, t vms.Token) (*events.WsCommandResponse, error)

	// 2- Subscribe to a topic
	RequestSubscribe(ctx context.Context, cameraID string, eventTypeID string) (*events.WsCommandResponse, error)

	// 3- Read events from an open session
	RequestEvents(ctx context.Context) (*events.AnalyticsEvents, error)

	// 4- Close communication
	RequestClose() error
}

type wsEventsRepository struct {
	base.WsBaseRepository
	sessionID   string
	lastEventID string
}

func NewWsEventsRepository() WsEventsRepository {
	return &wsEventsRepository{
		WsBaseRepository: base.NewWsBaseRepository(),
		sessionID:        "",
		lastEventID:      "",
	}
}

func (wer *wsEventsRepository) sendCommand(ctx context.Context, wreq *events.WsCommandRequest) (*events.WsCommandResponse, error) {
	commandRequestsCounter++
	wreq.CommandID = commandRequestsCounter

	// Send request to the websocket server
	if err := wer.SendRequest(ctx, wreq); err != nil {
		return nil, err
	}

	// Read from the websocket until we get a command response (discarding any other messages - e.g., events)
	wres := new(events.WsCommandResponse)
	if err := wer.ReadResponse(ctx, wres); err != nil {
		return nil, err
	}

	// Raise an exception with errorText if the status does not indicate success
	if wres.Status < 200 || wres.Status > 299 {
		return nil, errors.New("command failed - status: " + strconv.Itoa(wres.Status) + " - error: " + wres.Error.ErrorText)
	}
	return wres, nil
}

func (wer *wsEventsRepository) RequestStartSession(ctx context.Context, s vms.Server, t vms.Token) (*events.WsCommandResponse, error) {
	requestUrl, err := url.ParseRequestURI(s.ApiWellKnownUris.ApiGateways[0])
	if err != nil {
		return nil, err
	}

	requestUrl.Scheme = "ws"
	if s.IsSecure() {
		requestUrl.Scheme = "wss"
	}
	requestUrl.Path = constants.EventsWebsocket

	// Dial
	if err := wer.MakeConnect(ctx, requestUrl, t); err != nil {
		return nil, err
	}

	request := newStartSessionRequest()
	// If the session id or last event id are empty or null then we start a new session
	if strings.TrimSpace(wer.sessionID) == "" || strings.TrimSpace(wer.lastEventID) == "" {
		wer.sessionID = ""
		wer.lastEventID = ""
	}
	request.SessionID = wer.sessionID
	request.LastEventID = wer.lastEventID

	// Send request, read response, and parse to object
	wsCommandResponse, err := wer.sendCommand(ctx, request)
	if err != nil {
		return nil, err
	}
	wer.sessionID = wsCommandResponse.SessionID
	return wsCommandResponse, nil
}

func (wer *wsEventsRepository) RequestSubscribe(ctx context.Context, cameraID string, eventTypeID string) (*events.WsCommandResponse, error) {
	filters := newSubscriptionFilters()
	filters.Filters[0].SourceIDs = []string{cameraID}
	filters.Filters[0].EventTypes = []string{eventTypeID}
	request := newAddSubscriptionRequest(filters)

	// Send request, read response, and parse to object
	return wer.sendCommand(ctx, request)
}

func (wer *wsEventsRepository) RequestEvents(ctx context.Context) (*events.AnalyticsEvents, error) {
	aes := events.NewAnalyticsEvents()
	if err := wer.ReadResponse(ctx, aes); err != nil {
		return nil, err
	}
	// Get id of the last event
	wer.lastEventID = aes.Events[len(aes.Events)-1].ID
	return aes, nil
}

func (wer *wsEventsRepository) RequestClose() error {
	// Close and ignore any error
	defer wer.CloseConnect()
	wer.sessionID = ""
	wer.lastEventID = ""
	return nil
}
