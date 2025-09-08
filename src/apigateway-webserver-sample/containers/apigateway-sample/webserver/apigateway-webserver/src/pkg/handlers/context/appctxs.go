package context

import (
	"sync"
)

type AppContexts interface {
	AddAppContext(username string, ac AppContext)
	GetAppContext(username string) (AppContext, bool)
}

type appContexts struct {
	ctxs map[string]AppContext
	mu   sync.Mutex
}

var (
	instance AppContexts
	once     sync.Once
)

func GetAppContextsInstance() AppContexts {
	once.Do(func() {
		instance = &appContexts{
			ctxs: make(map[string]AppContext),
		}
	})
	return instance
}

func (acs *appContexts) AddAppContext(username string, ac AppContext) {
	acs.mu.Lock()
	defer acs.mu.Unlock()
	acs.ctxs[username] = ac
}

func (acs *appContexts) GetAppContext(username string) (AppContext, bool) {
	acs.mu.Lock()
	defer acs.mu.Unlock()
	ac, ok := acs.ctxs[username]
	return ac, ok
}
