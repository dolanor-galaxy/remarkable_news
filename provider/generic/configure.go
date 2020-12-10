package generic

import (
	"context"
	"net"
	"net/http"

	"github.com/sethvargo/go-envconfig"
)

// Configure the provider with environment variables
func (c *Provider) Configure(ctx context.Context) error {
	var l envconfig.Lookuper
	l = c.lookuper
	if c.lookuper == nil {
		l = envconfig.PrefixLookuper("NEWFECTHER_GENERIC_", envconfig.OsLookuper())
	}
	err := envconfig.ProcessWith(ctx, c, l)
	if err != nil {
		return err
	}
	// Check is URL is valid

	// Create the default client
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: c.TransportTimeout,
		}).Dial,
		TLSHandshakeTimeout: c.TransportTimeout,
	}

	c.httpClient = &http.Client{
		Timeout:   c.HTTPTimeout,
		Transport: netTransport,
	}
	return nil
}
