package vms

import (
	"encoding/json"
	"fmt"
)

type AnalyticEventType struct {
	ID           string   `json:"id"`
	Name         string   `json:"displayName"`
	Description  string   `json:"description"`
	LastModified string   `json:"lastModified"`
	Sources      []string `json:"sourceArray"`
}

type AnalyticEventTypes struct {
	Types []*AnalyticEventType
}

func (aets *AnalyticEventTypes) ToJSON() (string, error) {
	jsonData, err := json.Marshal(aets.Types)
	if err != nil {
		return "", fmt.Errorf("failed to marshal analytic event types: %w", err)
	}
	return string(jsonData), nil
}

func NewAnalyticEventTypes() *AnalyticEventTypes {
	return &AnalyticEventTypes{
		Types: []*AnalyticEventType{},
	}
}

func (aets *AnalyticEventTypes) Add(newEventType AnalyticEventType) {
	aets.Types = append(aets.Types, &newEventType)
}
