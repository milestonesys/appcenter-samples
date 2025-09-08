package base

import (
	"context"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"

	"apigateway-webserver/src/pkg/entities/vms"
)

type WsBaseRepository struct {
	BaseRepository
	conn   *websocket.Conn
	wg     sync.WaitGroup
	cancel context.CancelFunc
	mu     sync.Mutex
}

func NewWsBaseRepository() WsBaseRepository {
	return WsBaseRepository{
		BaseRepository: NewBaseRepository(),
		conn:           nil,
		wg:             sync.WaitGroup{},
		cancel:         nil,
	}
}

// Establishes a WebSocket connection to the given URL with the provided token
func (wbr *WsBaseRepository) MakeConnect(ctx context.Context, requestUrl *url.URL, token vms.Token) error {
	wbr.mu.Lock()
	defer wbr.mu.Unlock()

	var err error
	var ctxWithCancel context.Context

	// Close the connection if it is already open
	if err := wbr.CloseConnect(); err != nil {
		return err
	}

	// Set the request transport to support both encrypted and unencrypted communication
	wbr.setRequestTransport(requestUrl)

	header := make(http.Header)
	// Check if the token was provided and add it to the request header
	if token != nil {
		bearerToken, err := token.DispatchToken(ctx)
		if err != nil {
			return err
		}
		header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Perform a WebSocket handshake
	wbr.conn, _, err = websocket.Dial(ctx, requestUrl.String(), &websocket.DialOptions{
		HTTPClient: wbr.client,
		HTTPHeader: header,
		Host:       requestUrl.Host,
	})
	if err != nil {
		return err
	}

	// Start a ping pong chat with the server
	ctxWithCancel, wbr.cancel = context.WithCancel(ctx)
	wbr.keepAlive(ctxWithCancel)

	return nil
}

// Closes the WebSocket connection if it is open
func (wbr *WsBaseRepository) CloseConnect() error {
	if wbr.conn == nil {
		return nil
	}

	if wbr.cancel != nil {
		wbr.cancel()
		wbr.wg.Wait()
		wbr.cancel = nil
	}

	defer wbr.conn.CloseNow()

	err := wbr.conn.Close(websocket.StatusNormalClosure, "")
	if err != nil {
		return err
	}

	wbr.conn = nil
	return nil
}

// Sends a request over the WebSocket connection
func (wbr *WsBaseRepository) SendRequest(ctx context.Context, v any) error {
	wbr.mu.Lock()
	defer wbr.mu.Unlock()
	return wsjson.Write(ctx, wbr.conn, v)
}

// Reads a response from the WebSocket connection
func (wbr *WsBaseRepository) ReadResponse(ctx context.Context, v any) error {
	wbr.mu.Lock()
	defer wbr.mu.Unlock()
	return wsjson.Read(ctx, wbr.conn, v)
}

func (wbr *WsBaseRepository) keepAlive(ctx context.Context) {
	// Keep the server alive - ping every minute
	wbr.wg.Add(1)
	go func() {
		defer wbr.wg.Done()
		t := time.NewTimer(time.Minute)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				// Stopping the goroutines when close communication
				return
			case <-t.C:
			}

			// send ping
			err := wbr.conn.Ping(ctx)
			if err != nil {
				return
			}

			t.Reset(time.Minute)
		}
	}()
}
