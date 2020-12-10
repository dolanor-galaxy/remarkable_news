package generic

import (
	"context"
	"net/url"
	"reflect"
	"testing"

	"github.com/sethvargo/go-envconfig"
)

func TestProvider_Configure(t *testing.T) {
	u, _ := url.Parse("https://example.com")
	tests := map[string]func(t *testing.T){
		"no env": func(t *testing.T) {
			p := &Provider{
				lookuper: envconfig.MapLookuper(map[string]string{}),
			}
			ctx := context.Background()
			err := p.Configure(ctx)
			if err == nil {
				t.Fail()
			}
		},
		"url": func(t *testing.T) {
			p := &Provider{
				lookuper: envconfig.MapLookuper(map[string]string{
					"URL": u.String(),
				}),
			}
			ctx := context.Background()
			err := p.Configure(ctx)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(p.URL, u) {
				t.Errorf("Expected url %v, got %v", u, p.URL)
			}
		},
	}

	for name, tt := range tests {
		t.Run(name, tt)
	}
}
