package generic

import (
	"context"
	"encoding/base64"
	"image"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/sethvargo/go-envconfig"
)

func TestProvider_Process(t *testing.T) {
	must := func(img image.Image, _ string, _ error) image.Image {
		return img
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
		img image.Image
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"simple",
			fields{
				Mode: "fill",
			},
			args{
				context.TODO(),
				must(image.Decode(base64.NewDecoder(base64.StdEncoding, strings.NewReader(picture)))),
			},
			false,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Provider{
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
			img, err := p.Process(tt.args.ctx, tt.args.img)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provider.Process() error = %v, wantErr %v", err, tt.wantErr)
			}
			if img.Bounds().Max.X != reWidth {
				t.Errorf("Provider.Process() size = %v, size %v", img.Bounds(), 1492)
			}
		})
	}
}
