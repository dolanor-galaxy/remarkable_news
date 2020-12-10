package main

import (
	"context"
	"log"
	"time"

	"github.com/owulveryck/remarkable_news/provider"
	"github.com/owulveryck/remarkable_news/provider/generic"
	"github.com/disintegration/imaging"
	"github.com/sethvargo/go-envconfig"
)

type configuration struct {
	Output          string        `env:"OUTPUT,required"`
	UpdateFrequency time.Duration `env:"UPDATE_FREQUENCY,default=1h"`
}

func main() {
	var c configuration
	ctx := context.Background()
	l := envconfig.PrefixLookuper("RKNEWS_", envconfig.OsLookuper())
	err := envconfig.ProcessWith(ctx, &c, l)
	if err != nil {
		log.Fatal(err)
	}
	err = run(ctx, c)
	if err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, c configuration) error {
	prvd := &generic.Provider{}
	prvd.Configure(ctx)
	ticker := time.NewTicker(c.UpdateFrequency)
	err := runProvider(ctx, c, ticker, prvd)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			err := runProvider(ctx, c, ticker, prvd)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func runProvider(ctx context.Context, c configuration, ticker *time.Ticker, prvd provider.Provider) error {
	err := prvd.WaitOnline(ctx)
	if err != nil {
		return err
	}
	ticker.Reset(c.UpdateFrequency)
	img, err := prvd.Fetch(ctx)
	if err != nil {
		return err
	}
	img, err = prvd.Process(ctx, img)
	if err != nil {
		return err
	}
	err = imaging.Save(img, c.Output)
	if err != nil {
		return err
	}
	return nil
}
