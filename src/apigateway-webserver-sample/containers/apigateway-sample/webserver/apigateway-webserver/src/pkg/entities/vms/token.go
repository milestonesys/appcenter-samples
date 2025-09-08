package vms

import (
	"context"
	"encoding/json"
	"time"
)

// Ensures every concrete implementation can parse a token schema, check for the token expiration, and provide access to the token value.
type Token interface {
	// Returns a boolean indicating whether the token has expired or not.
	HasExpired() bool

	// Checks if the token has expired and renews it if necessary.
	// It returns the token value.
	DispatchToken(ctx context.Context) (string, error)

	// Copies all token schema values from the given token.
	Copy(copy Token) error

	// Returns the token schema.
	GetSchema() TokenSchema

	// Returns the token parsing timestamp in UTC.
	GetTimestampUTC() time.Time
}

// External function that checks for the token expiration and returns the dispatched (current or renewed) token based on that condition.
type TokenDispatchFunc func(ctx context.Context, current Token) error

type TokenSchema struct {
	AccessToken string `json:"access_token"` // active token
	ExpiresIn   int64  `json:"expires_in"`
	Type        string `json:"token_type"`
	Scope       string `json:"scope"`
}

type token struct {
	schema       TokenSchema
	timestampUTC time.Time
	dispatchFunc TokenDispatchFunc
}

func NewToken(tokenData []byte, tokenDispatchFunc TokenDispatchFunc) (Token, error) {
	t := &token{
		dispatchFunc: tokenDispatchFunc,
	}

	if err := json.Unmarshal(tokenData, &t.schema); err != nil {
		return nil, err
	}

	t.timestampUTC = time.Now().UTC()
	return t, nil
}

func (t *token) HasExpired() bool {
	return time.Since(t.timestampUTC) > time.Duration(t.schema.ExpiresIn)*time.Second
}

func (t *token) DispatchToken(ctx context.Context) (string, error) {
	if err := t.dispatchFunc(ctx, t); err != nil {
		return "", err
	}
	return t.schema.AccessToken, nil
}

func (t *token) Copy(copy Token) error {
	t.schema = copy.GetSchema()
	t.timestampUTC = copy.GetTimestampUTC()
	return nil
}

func (t *token) GetSchema() TokenSchema {
	return t.schema
}

func (t *token) GetTimestampUTC() time.Time {
	return t.timestampUTC
}
