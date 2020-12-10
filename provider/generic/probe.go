package generic

import (
	"context"
	"net/http"
	"time"
)

// WaitOnline check if the service is online, every tick; it returns when the service is up
// StartProbe only returns when the provider is reachable or in case of context cancelation
func (c *Provider) WaitOnline(ctx context.Context) error {
	if c.isOnline(ctx) {
		return nil
	}
	liveness := time.NewTicker(c.LivenessCheck)
	timeout := time.NewTicker(c.ProbeTimeout)
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()
	defer liveness.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-liveness.C:
			if c.isOnline(ctx) {
				return nil
			}
		case <-timeout.C:
			cancel()
		}
	}
}

func (c *Provider) isOnline(ctx context.Context) bool {
	headRequest, err := http.NewRequestWithContext(ctx, http.MethodHead, c.URL.String(), nil)
	if err != nil {
		return false
	}
	_, err = c.httpClient.Do(headRequest)
	if err != nil {
		return false
	}
	return true
}
