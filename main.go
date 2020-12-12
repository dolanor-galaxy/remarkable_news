package main

import (
	"context"
	"image"
	"image/color"
	"log"
	"time"

	"github.com/disintegration/imaging"
	"github.com/golang/freetype/truetype"
	"github.com/owulveryck/remarkable_news/provider"
	"github.com/owulveryck/remarkable_news/provider/generic"
	"github.com/sethvargo/go-envconfig"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	_ "golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/math/fixed"
)

type configuration struct {
	Output          string        `env:"OUTPUT,required"`
	UpdateFrequency time.Duration `env:"UPDATE_FREQUENCY,default=1h"`
	HealthCheck     time.Duration `env:"HEALTH_CHECK,default=30s"`
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

func hasSlept(ctx context.Context, healthCheck time.Duration, sleptTime time.Duration) <-chan struct{} {
	signalC := make(chan struct{})
	go func(chan<- struct{}) {
		tick := time.NewTicker(healthCheck)
		last := time.Now()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				if time.Now().Sub(last) > sleptTime {
					signalC <- struct{}{}
				}
				last = time.Now()
			}
		}
	}(signalC)
	return signalC
}

func run(ctx context.Context, c configuration) error {
	prvd := &generic.Provider{}
	prvd.Configure(ctx)
	ticker := time.NewTicker(c.UpdateFrequency)
	err := runProvider(ctx, c, ticker, prvd)
	if err != nil {
		log.Fatal(err)
	}
	runC := make(chan struct{})
	sleptC := hasSlept(ctx, c.HealthCheck, c.HealthCheck*2)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			runC <- struct{}{}
		case <-sleptC:
			runC <- struct{}{}
			ticker.Reset(c.UpdateFrequency)
		case <-runC:
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
	addLabel(img, 10, 40, "Remarkable is sleeping")
	addLabel(img, 10, 1850, "Updated at "+time.Now().Format("2/1/2006 15:04:05"))
	err = imaging.Save(img, c.Output)
	if err != nil {
		return err
	}
	return nil
}

func addLabel(i image.Image, x, y int, label string) error {
	fnt, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return err
	}
	face := truetype.NewFace(fnt, &truetype.Options{
		Size: 36,
	})
	if img, ok := i.(*image.NRGBA); ok {
		// img is now an *image.RGBA
		col := color.RGBA{0, 0, 0, 255}
		point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}

		d := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(col),
			Face: face, //basicfont.Face7x13,
			Dot:  point,
		}
		d.DrawString(label)
	}
	return nil
}
