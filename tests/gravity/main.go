package main

import (
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"math"
	"math/rand"
	"time"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/image/math/f64"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
)

// Asteroid represents an asteroid in space.
type Asteroid struct {
	w screen.Window
	s screen.Screen
	// The bounds of the Window.
	bounds     WidthHeight
	x, y       int
	size       int
	directionX int
	directionY int
	delay      time.Duration
	color      color.Color
}

// NewAsteroid creates an Asteroid type.
func NewAsteroid(w screen.Window, s screen.Screen, bounds WidthHeight, size int, directionX int, directionY int, delay time.Duration, color color.Color) *Asteroid {
	return &Asteroid{
		w:          w,
		s:          s,
		bounds:     bounds,
		size:       size,
		directionX: directionX,
		directionY: directionY,
		delay:      delay,
		color:      color,
	}
}

// Start starts rendering the Asteroid to the canvas.
func (a *Asteroid) Start(repaint chan PaintRequest) {
	a.x = rand.Intn(a.bounds.Width())
	a.y = rand.Intn(a.bounds.Height())

	var b screen.Buffer
	b, _ = a.s.NewBuffer(image.Point{a.size, a.size})
	img := b.RGBA()
	for ix := 0; ix < img.Bounds().Dx(); ix++ {
		for iy := 0; iy < img.Bounds().Dy(); iy++ {
			img.Set(ix, iy, a.color)
		}
	}

	for {
		if (a.x+a.size) > a.bounds.Width() || (a.x-a.size) < 0 {
			a.x = rand.Intn(a.bounds.Width())
		}
		if (a.y+a.size) > a.bounds.Height() || (a.y-a.size) < 0 {
			a.y = rand.Intn(a.bounds.Height())
		}

		// Draw with sd.
		/* sd := NewSimpleDraw(a.w, a.s, a.bounds)
		request := sd.Rectangle(a.x, a.y, a.size, a.size, a.color)
		*/

		repaint <- PaintRequest{
			dp:  image.Point{a.x, a.y},
			src: b,
			sr:  b.Bounds(),
		}

		a.x += a.directionX
		a.y += a.directionY

		time.Sleep(a.delay)
	}
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

		repaint := make(chan PaintRequest)

		asteroids := make([]*Asteroid, 15)
		for i := 0; i < len(asteroids); i++ {
			size := rand.Intn(80) + 20
			directionX := rand.Intn(10) - rand.Intn(20)
			directionY := rand.Intn(10) - rand.Intn(20)
			delay := rand.Intn(200000) + 1
			randomColor := rand.Intn(215) + 1
			color := palette.WebSafe[randomColor]
			asteroids[i] = NewAsteroid(w, s, bounds, size, directionX, directionY, time.Duration(delay), color)
			go asteroids[i].Start(repaint)
		}

		go painter(w, repaint)

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
	})
}

func painter(w screen.Window, repaint chan PaintRequest) {
	paints := 0
	since := time.Now()

	d, _ := time.ParseDuration("33ms")
	timer := time.NewTimer(d)

	for {
		select {
		case pr := <-repaint:
			w.Upload(pr.dp, pr.src, pr.src.Bounds())
			w.Publish()
			go pr.src.Release()

		case <-timer.C:
			paints++
			if time.Now().Sub(since).Seconds() >= 1 {
				fmt.Printf("%v fps\n", float64(paints)/time.Now().Sub(since).Seconds())
				paints = 0
				since = time.Now()
			}
			timer.Reset(time.Second / 30)
		}
	}
}

type PaintRequest struct {
	dp  image.Point
	src screen.Buffer
	sr  image.Rectangle
}

// SimpleDraw is a simple drawing tool.
type SimpleDraw struct {
	w      screen.Window
	s      screen.Screen
	bounds WidthHeight
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

func NewSimpleDraw(w screen.Window, s screen.Screen, bounds WidthHeight) *SimpleDraw {
	return &SimpleDraw{
		w:      w,
		s:      s,
		bounds: bounds,
	}
}

func (sd *SimpleDraw) Rectangle(x, y int, width, height int, c color.Color) PaintRequest {
	// Use the DrawUniform function.
	// sd.w.DrawUniform(affine.NewIdentity(), c,
	//	image.Rect(x, y, x+width, y+height), screen.Src, &screen.DrawOptions{})

	// Use an image.
	var b screen.Buffer
	b, _ = sd.s.NewBuffer(image.Point{width, height})
	img := b.RGBA()
	for ix := 0; ix < img.Bounds().Dx(); ix++ {
		for iy := 0; iy < img.Bounds().Dy(); iy++ {
			img.Set(ix, iy, c)
		}
	}
	return PaintRequest{
		dp:  image.Point{x, y},
		src: b,
		sr:  b.Bounds(),
	}

	// Use the Window fill directly.
	//sd.w.Fill(image.Rect(x, y, x+width, y+height), c, screen.Src)
}

func (sd *SimpleDraw) Line(src image.Point, dst image.Point, color color.Color) {
	op := screen.Src

	//cos30 := math.Cos(math.Pi / 6)
	// sin30 := math.Sin(math.Pi / 6)
	src2dst := f64.Aff3{
		1, .75, 500,
		-.75, 1, 200,
	}

	r := image.Rect(src.X, src.Y, dst.X, dst.Y)

	sd.w.DrawUniform(src2dst, color,
		r, op, &screen.DrawOptions{})
}

func (sd *SimpleDraw) Circle(dst image.Point, radius int, color color.Color) {
	op := screen.Over

	nib := image.Rect(dst.X, dst.Y, 1, dst.Y+1)

	for i := 1.0; i < 360; i++ {
		rads := i * math.Pi / 180.0
		src2dst := f64.Aff3{
			math.Cos(rads), -math.Sin(rads), float64(dst.X),
			math.Sin(rads), math.Cos(rads), float64(dst.Y),
		}
		sd.w.DrawUniform(src2dst, color,
			nib, op, &screen.DrawOptions{})
	}
}

/*
func (sd *SimpleDraw) Line(src image.Point, dst image.Point, color color.Color) {
	width := dst.X - src.X
	height := dst.Y - src.Y
	var slope float64
	if width != 0 {
		slope = float64(height) / float64(width)
	}

	// y = mx + b
	fmt.Printf("width: %v: height: %v\n", width, height)
	xFrom := min(src.X, dst.X)
	xTo := max(src.X, dst.X)
	for x := xFrom; x < xTo; x++ {
		y := (slope * float64(x))
		//fmt.Printf("draw point at %v, %v: y = %v * %v\n", x, y, slope, x)
		sd.Point(int(x), int(y), color)
	}
}
*/

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (sd *SimpleDraw) Point(x, y int, color color.Color) {
	r := image.Rect(x, y, x+1, y+1)
	sd.w.Fill(r, color, screen.Over)
}
