package provider

import (
	"context"
	"image"
)

// Provider is an image provider
type Provider interface {
	ImageFetcher
	Configurater
	Waiter
	ImageProcesser
}

// ImageFetcher can fetch an image remotly.
type ImageFetcher interface {
	//Fetch an image. It returns any error that prevents the image to be downloaded
	Fetch(context.Context) (image.Image, error)
}

// Configurater is any objet that can grab its configuration
type Configurater interface {
	Configure(context.Context) error
}

// Waiter waits for a service to be available
type Waiter interface {
	WaitOnline(context.Context) error
}

// ImageProcesser ...
type ImageProcesser interface {
	Process(context.Context, image.Image) (image.Image, error)
}

// OnlineData carries an error and a cancelfunc
type OnlineData struct {
	Err    error
	Cancel context.CancelFunc
}
