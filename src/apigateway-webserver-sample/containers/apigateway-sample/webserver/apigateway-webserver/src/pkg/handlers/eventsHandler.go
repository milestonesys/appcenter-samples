package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	handlers_context "apigateway-webserver/src/pkg/handlers/context"
)

type EventHandler struct {
	mu sync.Mutex
}

func NewEventHandler() *EventHandler {
	return &EventHandler{}
}

func (eh *EventHandler) StartSubscriptionHandle(w http.ResponseWriter, r *http.Request) {
	eh.mu.Lock()
	defer eh.mu.Unlock()
	log.Println("EventHandler.StartSubscriptionHandle() called")

	var data struct {
		Username    string `json:"username"`
		CameraId    string `json:"cameraId"`
		EventTypeId string `json:"eventTypeId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON format: %v", err), http.StatusBadRequest)
		return
	}

	if data.Username == "" || data.CameraId == "" || data.EventTypeId == "" {
		http.Error(w, "Missing required fields: Username, CameraId or EventTypeId.", http.StatusBadRequest)
		return
	}

	appCtx, exists := handlers_context.GetAppContextsInstance().GetAppContext(data.Username)
	if !exists {
		http.Error(w, "App handlers_context not found.", http.StatusBadRequest)
		return
	}

	// Close existing WebSocket connection
	if err := appCtx.WsEventsService().RequestClose(); err != nil {
		http.Error(w, fmt.Sprintf("While closing the previous websocket connection: %v", err), http.StatusInternalServerError)
		return
	}

	// Start new WebSocket connection
	wsResponse, err := appCtx.WsEventsService().RequestStartSession(r.Context(), appCtx.Server(), appCtx.Token())
	if err != nil {
		http.Error(w, fmt.Sprintf("While starting a new websocket connection: %v", err), http.StatusInternalServerError)
		return
	}

	appCtx.SetWsCommandResponse(wsResponse)

	// Subscribe for events filtered by type and source
	if _, err := appCtx.WsEventsService().RequestSubscribe(r.Context(), data.CameraId, data.EventTypeId); err != nil {
		http.Error(w, fmt.Sprintf("While creating a new subscription: %v", err), http.StatusInternalServerError)
		return
	}

	sessionJson, err := wsResponse.ToJSON()
	if err != nil {
		http.Error(w, fmt.Sprintf("Converting session info to JSON: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("{ \"message\": \"Processing started\", \"session\": %s }", sessionJson)))
}

func (eh *EventHandler) RequestEventsHandle(w http.ResponseWriter, r *http.Request) {
	eh.mu.Lock()
	defer eh.mu.Unlock()
	log.Println("EventHandler.RequestEventsHandle() called")

	var data struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON format: %v", err), http.StatusBadRequest)
		return
	}

	if data.Username == "" {
		http.Error(w, "Missing required fields: username", http.StatusBadRequest)
		return
	}

	appCtx, exists := handlers_context.GetAppContextsInstance().GetAppContext(data.Username)
	if !exists {
		http.Error(w, "App handlers_context not found.", http.StatusBadRequest)
		return
	}

	// Start reading events, if the communication was closed will return an error but we can ignore it since the we closed the communication gracefully
	aes, err := appCtx.WsEventsService().RequestEvents(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("Requesting events: %v", err), http.StatusInternalServerError)
		return
	}

	aesJson, err := aes.ToJSON()
	if err != nil {
		http.Error(w, fmt.Sprintf("Converting cameras list to JSON: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(aesJson))
}
