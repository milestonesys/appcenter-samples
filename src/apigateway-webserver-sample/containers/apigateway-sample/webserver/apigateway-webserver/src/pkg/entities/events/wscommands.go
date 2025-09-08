package events

import (
	"encoding/json"
	"fmt"
)

// WebSocket command request.
type WsCommandRequest struct {
	Command     string               `json:"command"`
	CommandID   int                  `json:"commandId"`
	SessionID   string               `json:"sessionId"`
	LastEventID string               `json:"eventId"`
	Filters     []SubscriptionFilter `json:"filters"`
}

// WebSocket command response.
type WsCommandResponse struct {
	SessionID string `json:"sessionId"`
	CommandID int    `json:"commandId"`
	Status    int    `json:"status"`
	Error     struct {
		ErrorText string `json:"errorText"`
	} `json:"error"`
}

func (cr *WsCommandResponse) ToJSON() (string, error) {
	jsonData, err := json.Marshal(cr)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}
	return string(jsonData), nil
}

// Subscriptions filter
type SubscriptionFilter struct {
	Modifier      string   `json:"modifier"`
	ResourceTypes []string `json:"resourceTypes"`
	SourceIDs     []string `json:"sourceIds"`
	EventTypes    []string `json:"eventTypes"`
}

type SubscriptionFilters struct {
	Filters []SubscriptionFilter `json:"filters"`
}
