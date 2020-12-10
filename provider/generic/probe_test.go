package generic

import (
	"context"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/sethvargo/go-envconfig"
)

func TestProvider_WaitOnline(t *testing.T) {
	mux := http.NewServeMux()
	ts := httptest.NewServer(mux)
	defer ts.Close()
	mux.HandleFunc("/html", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.New("html").Parse(htmlData))
		tmpl.Execute(w, ts)
	})

	urlize := func(s string) *url.URL {
		u, _ := url.Parse(ts.URL + s)
		return u
	}
	type fields struct {
		LivenessCheck    time.Duration
		ProbeTimeout     time.Duration
		HTTPTimeout      time.Duration
		TransportTimeout time.Duration
		URL              *url.URL
		Path             string
		Mode             string
		Scale            float64
		lookuper         envconfig.Lookuper
		httpClient       *http.Client
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"online",
			fields{
				httpClient:    http.DefaultClient,
				URL:           urlize("/html"),
				LivenessCheck: 5 * time.Millisecond,
				ProbeTimeout:  10 * time.Millisecond,
			},
			args{context.TODO()},
			false,
		},
		{
			"offline",
			fields{
				httpClient:    http.DefaultClient,
				URL:           &url.URL{},
				LivenessCheck: 5 * time.Millisecond,
				ProbeTimeout:  10 * time.Millisecond,
			},
			args{context.TODO()},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Provider{
				LivenessCheck:    tt.fields.LivenessCheck,
				ProbeTimeout:     tt.fields.ProbeTimeout,
				HTTPTimeout:      tt.fields.HTTPTimeout,
				TransportTimeout: tt.fields.TransportTimeout,
				URL:              tt.fields.URL,
				Path:             tt.fields.Path,
				Mode:             tt.fields.Mode,
				Scale:            tt.fields.Scale,
				lookuper:         tt.fields.lookuper,
				httpClient:       tt.fields.httpClient,
			}
			if err := c.WaitOnline(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Provider.WaitOnline() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
