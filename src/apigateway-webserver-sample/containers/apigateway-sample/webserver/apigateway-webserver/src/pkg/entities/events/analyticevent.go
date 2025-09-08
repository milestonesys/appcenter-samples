package events

import (
	"encoding/json"
	"fmt"
)

type AnalyticsEvent struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Source    string `json:"source"`
	Timestamp string `json:"time"`
	Datatype  string `json:"datatype"`
}

type AnalyticsEvents struct {
	Events []AnalyticsEvent `json:"events"`
}

func (ae *AnalyticsEvents) ToJSON() (string, error) {
	jsonData, err := json.Marshal(ae.Events)
	if err != nil {
		return "", fmt.Errorf("failed to marshal analytics events: %w", err)
	}
	return string(jsonData), nil
}

func NewAnalyticsEvents() *AnalyticsEvents {
	return &AnalyticsEvents{
		Events: []AnalyticsEvent{},
	}
}

func (ae *AnalyticsEvents) AddAnalyticsEvents(newEvents *AnalyticsEvents) {
	ae.Events = append(ae.Events, newEvents.Events...)
}

func (ae *AnalyticsEvents) Clear() {
	ae.Events = []AnalyticsEvent{}
}
