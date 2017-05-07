package main

import (
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"math/rand"
	"time"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
)

// Asteroid represents an asteroid in space.
type Asteroid struct {
	// The bounds of the Window.
	bounds     WidthHeight
	x, y       int
	directionX int
	directionY int
	delay      time.Duration
	color      color.Color
	image      *image.RGBA
}

// NewAsteroid creates an Asteroid type.
func NewAsteroid(b screen.Buffer, bounds WidthHeight, size int, directionX int, directionY int, delay time.Duration, color color.Color) *Asteroid {
	a := &Asteroid{
		bounds:     bounds,
		directionX: directionX,
		directionY: directionY,
		delay:      delay,
		color:      color,
		image:      b.RGBA(),
	}

	a.x = rand.Intn(a.bounds.Width())
	a.y = rand.Intn(a.bounds.Height())

	for ix := 0; ix < a.image.Bounds().Dx(); ix++ {
		for iy := 0; iy < a.image.Bounds().Dy(); iy++ {
			a.image.Set(ix, iy, a.color)
		}
	}

	return a
}

type PaintResult struct {
	Image *image.RGBA
	At    image.Point
}

func NewPaintResult(x, y int, img *image.RGBA) PaintResult {
	return PaintResult{
		Image: img,
		At:    image.Point{x, y},
	}
}

// Start starts rendering the Asteroid to the canvas.
func (a *Asteroid) Tick() PaintResult {
	sizeX, sizeY := a.image.Bounds().Dx(), a.image.Bounds().Dy()

	if (a.x+sizeX) > a.bounds.Width() || (a.x-sizeX) < 0 {
		a.x = rand.Intn(a.bounds.Width())
	}
	if (a.y+sizeY) > a.bounds.Height() || (a.y-sizeY) < 0 {
		a.y = rand.Intn(a.bounds.Height())
	}

	rv := NewPaintResult(a.x, a.y, a.image)

	a.x += a.directionX
	a.y += a.directionY

	return rv
}

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

		asteroids := make([]*Asteroid, 15)
		buffers := make([]screen.Buffer, 15)

		for i := 0; i < len(asteroids); i++ {
			size := rand.Intn(80) + 20
			directionX := rand.Intn(10) - rand.Intn(20)
			directionY := rand.Intn(10) - rand.Intn(20)
			delay := rand.Intn(200000) + 1
			randomColor := rand.Intn(215) + 1
			color := palette.WebSafe[randomColor]

			b, _ := s.NewBuffer(image.Point{size, size})
			buffers[i] = b
			asteroids[i] = NewAsteroid(b, bounds, size, directionX, directionY, time.Duration(delay), color)
		}

		defer repaint(w)

		for {
			for i, a := range asteroids {
				b := buffers[i]
				pr := a.Tick()
				w.Upload(pr.At, b, b.Bounds())
			}
			w.Publish()

			time.Sleep(time.Second / 1)
		}
	})
}

func repaint(w screen.Window) {
	for {
		switch e := w.NextEvent().(type) {
		case lifecycle.Event:
			if e.To == lifecycle.StageDead {
				return
			}
		case paint.Event:
			//repaint <- image.Rect(0, 0, bounds.Width(), bounds.Height())
		case error:
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
