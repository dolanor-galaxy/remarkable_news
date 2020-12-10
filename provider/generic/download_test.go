package generic

import (
	"context"
	"encoding/base64"
	"html/template"
	"image"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestProvider_Fetch(t *testing.T) {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(picture))
	m, _, err := image.Decode(reader)
	if err != nil {
		t.Fatal(err)
	}
	mux := http.NewServeMux()
	ts := httptest.NewServer(mux)
	defer ts.Close()
	mux.HandleFunc("/image", func(w http.ResponseWriter, r *http.Request) {
		io.Copy(w, base64.NewDecoder(base64.StdEncoding, strings.NewReader(picture)))
	})
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
		httpClient       *http.Client
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    image.Image
		wantErr bool
	}{
		{
			"grab image",
			fields{
				URL:        urlize("/image"),
				httpClient: http.DefaultClient,
			},
			args{context.TODO()},
			m,
			false,
		},
		{
			"grab html, nil expression",
			fields{
				URL:        urlize("/html"),
				httpClient: http.DefaultClient,
			},
			args{context.TODO()},
			nil,
			true,
		},
		{
			"grab html, notfound path",
			fields{
				URL:        urlize("/html"),
				httpClient: http.DefaultClient,
				Path:       "/",
			},
			args{context.TODO()},
			nil,
			true,
		},
		{
			"grab html, path no match",
			fields{
				URL:        urlize("/html"),
				httpClient: http.DefaultClient,
				Path:       `"@id="XXXXX"`,
			},
			args{context.TODO()},
			nil,
			true,
		},
		{
			"grab html, path relative",
			fields{
				URL:        urlize("/html"),
				httpClient: http.DefaultClient,
				Path:       `//img[@id="relative"]/@src`,
			},
			args{context.TODO()},
			m,
			false,
		},
		{
			"grab html, path absolute",
			fields{
				URL:        urlize("/html"),
				httpClient: http.DefaultClient,
				Path:       `//img[@id="absolute"]/@src`,
			},
			args{context.TODO()},
			m,
			false,
		},
		{
			"grab html, path full",
			fields{
				URL:        urlize("/html"),
				httpClient: http.DefaultClient,
				Path:       `//img[@id="full"]/@src`,
			},
			args{context.TODO()},
			m,
			false,
		},
		// TODO: Add test cases.
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
				httpClient:       tt.fields.httpClient,
			}
			got, err := c.Fetch(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Provider.Fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Provider.Fetch() = %v, want %v", got, tt.want)
			}
		})
	}
}
