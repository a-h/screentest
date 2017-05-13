package main

import (
	"fmt"
	"image"
	"image/color"

	"math"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
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

		xScale := Scaler{
			InputMinimum:  0,
			InputMaximum:  bounds.Width(),
			OutputMinimum: -10,
			OutputMaximum: +10,
		}

		yScale := Scaler{
			InputMinimum:  -1,
			InputMaximum:  +1,
			OutputMinimum: bounds.Height(),
			OutputMaximum: 0,
		}

		// Draw the x-axis.
		for ix := 0; ix < bounds.Width(); ix++ {
			x := ix
			y := yScale.Calculate(0)
			img.Set(x, int(y), axesColor)
		}

		// Draw the y-axis.
		for iy := 0; iy < bounds.Height(); iy++ {
			x := xScale.CalculateInverse(0)
			y := iy
			img.Set(int(x), y, axesColor)
		}

		for ix := 0; ix < bounds.Width(); ix++ {
			// Draw htan.
			x := xScale.Calculate(float64(ix))
			y := htan(x)
			iy := int(yScale.Calculate(y))
			img.Set(ix, iy, colorA)

			// Draw sigmoid.
			x = xScale.Calculate(float64(ix))
			y = sigmoid(x)
			iy = int(yScale.Calculate(y))
			img.Set(ix, iy, colorB)
		}
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
			case error:
			}
		}
	})
}

type Scaler struct {
	InputMinimum  int
	InputMaximum  int
	OutputMinimum int
	OutputMaximum int
}

func (s Scaler) Calculate(in float64) float64 {
	outputRange := float64(s.OutputMaximum - s.OutputMinimum)
	inputRange := float64(s.InputMaximum - s.InputMinimum)

	return ((outputRange / inputRange) * (in - float64(s.InputMinimum))) + float64(s.OutputMinimum)
}

func (s Scaler) CalculateInverse(in float64) float64 {
	outputRange := float64(s.InputMaximum - s.InputMinimum)
	inputRange := float64(s.OutputMaximum - s.OutputMinimum)

	return ((outputRange / inputRange) * (in - float64(s.OutputMinimum))) + float64(s.InputMinimum)
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
