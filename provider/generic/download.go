package generic

import (
	"bytes"
	"context"
	"errors"
	"image"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/disintegration/imaging"
)

//Fetch an image. It returns any error that prevents the image to be downloaded
func (c *Provider) Fetch(ctx context.Context) (image.Image, error) {
	return c.fetch(ctx, c.URL.String())
}

func (c *Provider) fetch(ctx context.Context, u string) (image.Image, error) {
	ur, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	if ur.Host == "" {
		ur.Host = c.URL.Host
	}
	if ur.Scheme == "" {
		ur.Scheme = c.URL.Scheme
	}
	u = ur.String()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	contentType := http.DetectContentType(body)
	switch {
	case strings.Contains(contentType, "image/"):
		return imaging.Decode(bytes.NewBuffer(body))
	case strings.Contains(contentType, "text/html"):
		doc, err := htmlquery.LoadURL(u)
		if err != nil {
			return nil, err
		}

		list, err := htmlquery.QueryAll(doc, c.Path)
		if err != nil {
			return nil, err
		}

		if len(list) == 0 {
			return nil, errors.New("no match")
		}
		bla := htmlquery.InnerText(list[0])
		//return c.fetch(ctx, htmlquery.InnerText(list[0]))
		return c.fetch(ctx, bla)
	default:
		return nil, errors.New("unhandled content type")
	}
}
