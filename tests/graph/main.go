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

		border := 50

		xScale := Scaler{
			InputMinimum:  0 + border,
			InputMaximum:  bounds.Width() - border,
			OutputMinimum: -10,
			OutputMaximum: +10,
		}

		yScale := Scaler{
			InputMinimum:  -1,
			InputMaximum:  +1,
			OutputMinimum: bounds.Height() - border,
			OutputMaximum: 0 + border,
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

		// Middle to top right.
		Line(img, 900, 500, 1800, 0, color.RGBA{R: 0xFF, G: 0x00, B: 0x00})
		// Middle to bottom right.
		Line(img, 900, 500, 1800, 1000, color.RGBA{R: 0x00, G: 0xFF, B: 0x00})
		// Middle to top left.
		Line(img, 900, 500, 0, 0, color.RGBA{R: 0x00, G: 0x00, B: 0xFF})
		//Middle to bottom left.
		Line(img, 900, 500, 0, 1000, color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF})

		// Top to bottom.
		// Line(img, 30, 0, 30, 1800, color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF})
		// Left to right.
		// Line(img, 0, 30, 1800, 30, color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF})

		Circle(img, 900, 500, 60, color.RGBA{R: 0xFF, G: 0x00, B: 0x00})

		Disc(img, 1200, 700, 60, color.RGBA{R: 0x00, G: 0x33, B: 0x00})

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

func Circle(img *image.RGBA, x, y int, radius int, c color.RGBA) {
	bounds := image.Rect(x-radius-2, y-radius-2, x+radius+2, y+radius+2)
	for ix := bounds.Min.X; ix < bounds.Max.X; ix++ {
		for iy := bounds.Min.Y; iy < bounds.Max.Y; iy++ {
			width := x - ix
			height := y - iy

			distanceFromCenter := math.Sqrt(float64(((width * width) + (height * height))))
			if int(distanceFromCenter) == radius {
				img.Set(ix, iy, c)
			}
		}
	}
}

func Disc(img *image.RGBA, x, y int, radius int, c color.RGBA) {
	bounds := image.Rect(x-radius-2, y-radius-2, x+radius+2, y+radius+2)
	for ix := bounds.Min.X; ix < bounds.Max.X; ix++ {
		for iy := bounds.Min.Y; iy < bounds.Max.Y; iy++ {
			width := x - ix
			height := y - iy

			distanceFromCenter := math.Sqrt(float64(((width * width) + (height * height))))
			if int(distanceFromCenter) <= radius {
				img.Set(ix, iy, c)
			}
		}
	}
}

func Line(img *image.RGBA, fromX, fromY int, toX, toY int, c color.RGBA) {
	// Vertical line.
	if fromX == toX {
		for y := fromY; y < toY; y++ {
			img.Set(fromX, y, c)
		}
		return
	}

	// Horizontal line.
	if fromY == toY {
		for x := fromX; x < toX; x++ {
			img.Set(x, fromY, c)
		}
		return
	}

	// It's a slope.
	// We're moving from fromX to toX, so make sure they're in the right order.
	if toX < fromX {
		toX, toY, fromX, fromY = fromX, fromY, toX, toY
	}

	var b int
	if toY < fromY {
		b = img.Bounds().Dy()
	}

	rise := toY - fromY
	run := toX - fromX
	m := float64(rise) / float64(run)

	for x := fromX; x <= toX; x++ {
		y := (m * float64(x)) + float64(b)
		img.Set(x, int(y), c)
	}
}
