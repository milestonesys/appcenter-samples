package repositories

import (
	"context"
	"sync"

	"apigateway-webserver/src/pkg/entities/vms"
)

// Ensure a token dispatcher is responsible for the token dispatch function
type TokenDispatcher interface {
	// Return the active bearer token
	// If the token is nil or expired, will request new one
	DispatchFunc() vms.TokenDispatchFunc
}

// Given an user and a server will implement a function that when a token is provided will call the IDP of that server to renew the user token
type tokenDispatcher struct {
	idpRepo IdpRepository
	user    *vms.User
	server  *vms.Server
	mu      sync.Mutex
}

func NewTokenDispatcher(idpRepo IdpRepository, u *vms.User, s *vms.Server) TokenDispatcher {
	return &tokenDispatcher{
		idpRepo: idpRepo,
		user:    u,
		server:  s,
	}
}

// Return an implementation of the TokenDispatchFunc function
func (td *tokenDispatcher) DispatchFunc() vms.TokenDispatchFunc {
	return func(ctx context.Context, current vms.Token) error {
		// A function that returns true or false based on some checks over the token to whether renew it or not
		tokenRenewCondition := func() bool {
			if current != nil {
				return current.HasExpired()
			}
			// Token is nil
			return true
		}

		// Execute function defined above
		if tokenRenewCondition() {
			td.mu.Lock()
			defer td.mu.Unlock()
			// Execute function defined above again after being inside the mutex lock
			if tokenRenewCondition() {
				dispatched, err := td.idpRepo.RequestAccessToken(ctx, *td.user, *td.server, td)
				if err != nil {
					return err
				}

				// Copy new token into current token
				return current.Copy(dispatched)
			}
		}
		return nil
	}
}
