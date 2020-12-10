// Package generic implements a generic provider that can fetch an url and grap the picture form an XPath
package generic

import (
	"net/http"
	"net/url"
	"time"

	"github.com/sethvargo/go-envconfig"
)

// Provider provider
type Provider struct {
	LivenessCheck    time.Duration `env:"LIVENESS_CHECK,default=5m"`
	ProbeTimeout     time.Duration `env:"PROBE_TIMEOUT,default=60m"`
	HTTPTimeout      time.Duration `env:"HTTP_TIMEOUT,default=10s"`
	TransportTimeout time.Duration `env:"TRANSPORT_TIMEOUT,default=5s"`
	URL              *url.URL      `env:"URL,required"`
	Path             string        `env:"XPATH"`
	Mode             string        `env:"MODE,default=center"` // fill is possible
	Scale            float64       `env:"SCALE, default=1"`
	lookuper         envconfig.Lookuper
	httpClient       *http.Client
}
