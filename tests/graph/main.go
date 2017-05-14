package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/a-h/raster"
	"github.com/a-h/scaler"

	"math"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
)

func main() {
	bounds := Bounds{w: 1800, h: 1000}

	driver.Main(func(s screen.Screen) {
		w, err := s.NewWindow(&screen.NewWindowOptions{
			Height: bounds.Height(),
			Width:  bounds.Width(),
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		defer w.Release()

		background, _ := s.NewBuffer(image.Point{bounds.Width(), bounds.Height()})
		img := background.RGBA()

		backgroundColor := color.RGBA{R: 0x00, G: 0x11, B: 0x11}
		axesColor := color.RGBA{R: 0x66, G: 0x66, B: 0x66}
		colorA := color.RGBA{R: 0xFF, G: 0x55, B: 0x55}
		colorB := color.RGBA{R: 0x55, G: 0xFF, B: 0x55}

		drawBackground(img, backgroundColor)

		border := 50

		xScaler := scaler.New(
			0+border,              // Input Minimum
			bounds.Width()-border, // Input Maximum
			-10, // Output Minimum
			+10) // Output Maximum

		yScaler := scaler.New(-1, +1, bounds.Height()-border, 0+border)

		// Draw the x-axis.
		for ix := 0; ix < bounds.Width(); ix++ {
			x := ix
			y, _ := yScaler.Scale(0)
			img.Set(x, int(y), axesColor)
		}

		// Draw the y-axis.
		for iy := 0; iy < bounds.Height(); iy++ {
			x, _ := xScaler.Invert(0)
			y := iy
			img.Set(int(x), y, axesColor)
		}

		for ix := 0; ix < bounds.Width(); ix++ {
			// Draw htan.
			x, _ := xScaler.Scale(float64(ix))
			y := htan(x)
			iy, _ := yScaler.Scale(y)
			img.Set(ix, int(iy), colorA)

			// Draw sigmoid.
			x, _ = xScaler.Scale(float64(ix))
			y = sigmoid(x)
			iy, _ = yScaler.Scale(y)
			img.Set(ix, int(iy), colorB)
		}

		// Middle to top right.
		raster.DrawLine(img, 900, 500, 1800, 0, color.RGBA{R: 0xFF, G: 0x00, B: 0x00})
		// Middle to bottom right.
		raster.DrawLine(img, 900, 500, 1800, 1000, color.RGBA{R: 0x00, G: 0xFF, B: 0x00})
		// Middle to top left.
		raster.DrawLine(img, 900, 500, 0, 0, color.RGBA{R: 0x00, G: 0x00, B: 0xFF})
		//Middle to bottom left.
		raster.DrawLine(img, 900, 500, 0, 1000, color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF})

		// Top to bottom.
		// raster.Line(img, 30, 0, 30, 1800, color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF})
		// Left to right.
		// raster.Line(img, 0, 30, 1800, 30, color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF})

		raster.DrawCircle(img, 900, 500, 60, color.RGBA{R: 0xFF, G: 0x00, B: 0x00})

		raster.DrawDisc(img, 1200, 700, 60, color.RGBA{R: 0x00, G: 0x33, B: 0x00})

		w.Upload(image.Point{0, 0}, background, image.Rect(0, 0, bounds.Width(), bounds.Height()))
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

func htan(x float64) float64 {
	e := math.Pow(math.E, (2 * x))
	return (e - 1) / (e + 1)
}

func sigmoid(x float64) float64 {
	return 1 / (1 + math.Pow(math.E, -x))
}

func drawBackground(img *image.RGBA, c color.RGBA) {
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			img.Set(x, y, c)
		}
	}
}

// Bounds provides Width and Height.
type Bounds struct {
	// Width
	w int
	// Height
	h int
}

// Width is the maximum pixel width.
func (b Bounds) Width() int {
	return b.w
}

// Height is the maximum pixel height.
func (b Bounds) Height() int {
	return b.h
}

// WidthHeight is any type which has Width() and Height() functions.
type WidthHeight interface {
	Width() int
	Height() int
}
