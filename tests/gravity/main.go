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
	x, y       float64
	directionX float64
	directionY float64
	delay      time.Duration
	color      color.Color
	image      *image.RGBA
}

// NewAsteroid creates an Asteroid type.
func NewAsteroid(bounds WidthHeight, size int, directionX float64, directionY float64, delay time.Duration, color color.Color) *Asteroid {
	a := &Asteroid{
		size:       Bounds{size, size},
		bounds:     bounds,
		directionX: directionX,
		directionY: directionY,
		delay:      delay,
		color:      color,
	}

	a.x = rand.Float64() * float64(a.bounds.Width())
	a.y = rand.Float64() * float64(a.bounds.Height())

	return a
}

// Start starts rendering the Asteroid to the canvas.
func (a *Asteroid) Tick() (at image.Point, img *image.RGBA, changed bool) {
	if a.x < 0 {
		a.x = 0
	}
	if a.x > float64(a.bounds.Width()) {
		a.x = float64(a.bounds.Width())
	}
	if a.y < 0 {
		a.y = 0
	}
	if a.y > float64(a.bounds.Height()) {
		a.y = float64(a.bounds.Height())
	}

	if (a.x+float64(a.size.Width())) >= float64(a.bounds.Width()) || (a.x-float64(a.size.Width())) <= 0 {
		a.directionX = a.directionX * -1
	}
	if (float64(a.y)+float64(a.size.Height())) >= float64(a.bounds.Height()) || (a.y-float64(a.size.Height())) <= 0 {
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
	oldX := int(a.x)
	oldY := int(a.y)

	a.x += a.directionX
	a.y += a.directionY

	hasMoved := (int(a.x) != oldX) || (int(a.y) != oldY)

	return image.Point{oldX, oldY}, img, hasMoved
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

		asteroids := make([]*Asteroid, 30)

		for i := 0; i < len(asteroids); i++ {
			size := rand.Intn(80) + 20
			directionX := (rand.Float64() * 10.0) - (rand.Float64() * 20.0)
			directionY := (rand.Float64() * 10.0) - (rand.Float64() * 20.0)
			delay := rand.Intn(200000) + 1
			randomColor := rand.Intn(215) + 1
			color := palette.WebSafe[randomColor]

			asteroids[i] = NewAsteroid(bounds, size, directionX, directionY, time.Duration(delay), color)
		}

		defer repaint(w)

		background, _ := s.NewBuffer(image.Point{bounds.Width(), bounds.Height()})
		tickIndex := 0
		bkg := background.RGBA()
		drawBackground(bkg)

		spriteLocations := make([]image.Rectangle, len(asteroids))

		for {
			clearSpriteLocations(bkg, spriteLocations)
			for i, a := range asteroids {
				at, img, hasMoved := a.Tick()
				if !hasMoved {
					continue
				}
				location := image.Rect(at.X, at.Y, at.X+img.Rect.Dx(), at.Y+img.Rect.Dy())
				draw.Draw(bkg,
					location,
					img,
					image.Point{0, 0},
					draw.Over)
				spriteLocations[i] = location
			}
			w.Upload(image.Point{0, 0}, background, image.Rect(0, 0, 1800, 1000))
			w.Publish()
			tickIndex++
			fmt.Println(tickIndex)
		}
	})
}

func clearSpriteLocations(img *image.RGBA, locations []image.Rectangle) {
	for _, l := range locations {
		for x := l.Min.X; x < l.Max.X; x++ {
			for y := l.Min.Y; y < l.Max.Y; y++ {
				img.Set(x, y, color.White)
			}
		}
	}
}

func drawBackground(img *image.RGBA) {
	for i := range img.Pix {
		img.Pix[i] = 0xff
	}
	/*
		for i := 0; i < len(img.Pix); i += 4 {
			img.Pix[i+0] = 0xff // R
			img.Pix[i+1] = 0xff // G
			img.Pix[i+2] = 0xff // B
			img.Pix[i+3] = 0xff // A
		}*/
	/*
		for x := 0; x < img.Bounds().Dx(); x++ {
			for y := 0; y < img.Bounds().Dy(); y++ {
				img.Set(x, y, color.White)
			}
		}
	*/
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
