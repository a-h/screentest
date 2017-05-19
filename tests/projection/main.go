package main

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/a-h/raster"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
)

const (
	width, height = 1800, 900
	cells         = 100
	xyrange       = 30.0
	xyscale       = width / 2 / xyrange
	zscale        = height * 0.4
	angle         = math.Pi / 6
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle)

func main() {
	driver.Main(func(s screen.Screen) {
		w, err := s.NewWindow(&screen.NewWindowOptions{
			Height: height,
			Width:  width,
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		defer w.Release()

		background, _ := s.NewBuffer(image.Point{width, height})
		img := background.RGBA()

		backgroundColor := color.RGBA{R: 0x00, G: 0x11, B: 0x11}
		colorA := color.RGBA{R: 0xFF, G: 0x55, B: 0x55}
		//colorB := color.RGBA{R: 0x55, G: 0xFF, B: 0x55}

		drawBackground(img, backgroundColor)

		for i := 0; i < cells; i++ {
			for j := 0; j < cells; j++ {
				ax, ay, ok := corner(i+1, j)
				if !ok {
					continue
				}
				bx, by, ok := corner(i, j)
				if !ok {
					continue
				}
				cx, cy, ok := corner(i, j+1)
				if !ok {
					continue
				}
				dx, dy, ok := corner(i+1, j+1)
				if !ok {
					continue
				}

				a := image.Point{int(ax), int(ay)}
				b := image.Point{int(bx), int(by)}
				c := image.Point{int(cx), int(cy)}
				d := image.Point{int(dx), int(dy)}

				raster.DrawPolygon(img, colorA, a, b, c, d)
				//w.Upload(image.Point{0, 0}, background, image.Rect(0, 0, width, height))
				//w.Publish()
			}
		}

		//border := 50

		/*
				xScaler := scaler.New(
			/*
				xScaler := scaler.New(
					0+border,              // Input Minimum
					bounds.Width()-border, // Input Maximum
					-30, // Output Minimum
					+30) // Output Maximum

				//yScaler := scaler.New(-1, +1, bounds.Height()-border, 0+border)

				for ix := 0; ix < bounds.Width(); ix++ {
					x, _ := xScaler.Scale(float64(ix))
					y := sin30 * x
					img.Set(int(ix), int(y), colorA)
				}
		*/
		// Left to right.
		// raster.Line(img, 0, 30, 1800, 30, color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF})

		w.Upload(image.Point{0, 0}, background, image.Rect(0, 0, width, height))
		w.Publish()

		// Keep looking for events.
		for {
			switch e := w.NextEvent().(type) {
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					fmt.Println("StageDead...")
					return
				}
			case key.Event:
				if e.Code == key.CodeUpArrow {
					// Do something.
				}
			case error:
			}
		}
	})
}

func corner(i, j int) (float64, float64, bool) {
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)
	z := f(x, y)

	sx := width/2 + (x-y)*cos30*xyscale //- z*zscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale

	if math.IsNaN(sx) || math.IsNaN(sy) {
		return 0, 0, false
	}

	return sx, sy, true
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y)
	return math.Sin(r) / r
}

func drawBackground(img *image.RGBA, c color.RGBA) {
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			img.Set(x, y, c)
		}
	}
}
