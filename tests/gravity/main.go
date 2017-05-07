package main

import (
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"math/rand"
	"time"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/lifecycle"
)

// Asteroid represents an asteroid in space.
type Asteroid struct {
	// The bounds of the Window.
	bounds     WidthHeight
	size       WidthHeight
	x, y       int
	directionX int
	directionY int
	delay      time.Duration
	color      color.Color
	image      *image.RGBA
}

// NewAsteroid creates an Asteroid type.
func NewAsteroid(bounds WidthHeight, size int, directionX int, directionY int, delay time.Duration, color color.Color) *Asteroid {
	a := &Asteroid{
		size:       Bounds{size, size},
		bounds:     bounds,
		directionX: directionX,
		directionY: directionY,
		delay:      delay,
		color:      color,
	}

	a.x = rand.Intn(a.bounds.Width())
	a.y = rand.Intn(a.bounds.Height())

	return a
}

// Start starts rendering the Asteroid to the canvas.
func (a *Asteroid) Tick() (at image.Point, img *image.RGBA) {
	if a.x < 0 {
		a.x = 0
	}
	if a.x > a.bounds.Width() {
		a.x = a.bounds.Width()
	}
	if a.y < 0 {
		a.y = 0
	}
	if a.y > a.bounds.Height() {
		a.y = a.bounds.Height()
	}

	if (a.x+a.size.Width()) >= a.bounds.Width() || (a.x-a.size.Width()) <= 0 {
		a.directionX = a.directionX * -1
	}
	if (a.y+a.size.Height()) >= a.bounds.Height() || (a.y-a.size.Height()) <= 0 {
		a.directionY = a.directionY * -1
	}

	// Get the canvas and draw on it.
	img = image.NewRGBA(image.Rect(0, 0, a.size.Width(), a.size.Height()))

	for ix := 0; ix < a.size.Width(); ix++ {
		for iy := 0; iy < a.size.Height(); iy++ {
			img.Set(ix, iy, a.color)
		}
	}

	// Return where to draw the sprite.
	at = image.Point{a.x, a.y}

	a.x += a.directionX
	a.y += a.directionY

	return at, img
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

		for i := 0; i < len(asteroids); i++ {
			size := rand.Intn(80) + 20
			directionX := rand.Intn(5) - rand.Intn(5)
			directionY := rand.Intn(5) - rand.Intn(5)
			delay := rand.Intn(200000) + 1
			randomColor := rand.Intn(215) + 1
			color := palette.WebSafe[randomColor]

			asteroids[i] = NewAsteroid(bounds, size, directionX, directionY, time.Duration(delay), color)
		}

		defer repaint(w)

		background, _ := s.NewBuffer(image.Point{bounds.Width(), bounds.Height()})
		tickIndex := 0

		for {
			drawBackground(background.RGBA())
			for _, a := range asteroids {
				at, img := a.Tick()
				draw.Draw(background.RGBA(),
					image.Rect(at.X, at.Y, at.X+img.Rect.Dx(), at.Y+img.Rect.Dy()),
					img,
					image.Point{0, 0},
					draw.Over)
			}
			w.Upload(image.Point{0, 0}, background, image.Rect(0, 0, 1800, 1000))
			w.Publish()
			tickIndex++
		}
	})
}

func drawBackground(img *image.RGBA) {
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			img.Set(x, y, color.White)
		}
	}
}

func repaint(w screen.Window) {
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
