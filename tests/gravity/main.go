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
	"golang.org/x/image/math/f64"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"

	"github.com/fogleman/gg"
)

// Asteroid represents an asteroid in space.
type Asteroid struct {
	canvas     *gg.Context
	x, y       int
	size       int
	directionX int
	directionY int
	delay      time.Duration
	color      color.Color
}

// NewAsteroid creates an Asteroid type.
func NewAsteroid(canvas *gg.Context, size int, directionX int, directionY int, delay time.Duration, color color.Color) *Asteroid {
	return &Asteroid{
		canvas:     canvas,
		size:       size,
		directionX: directionX,
		directionY: directionY,
		delay:      delay,
		color:      color,
	}
}

// Start starts rendering the Asteroid to the canvas.
func (a *Asteroid) Start(repaint chan image.Rectangle) {
	a.x = rand.Intn(a.canvas.Width())
	a.y = rand.Intn(a.canvas.Height())

	for {
		if (a.x+a.size) > a.canvas.Width() || (a.x-a.size) < 0 {
			a.x = rand.Intn(a.canvas.Width())
		}
		if (a.y+a.size) > a.canvas.Height() || (a.y-a.size) < 0 {
			a.y = rand.Intn(a.canvas.Height())
		}

		repaintBounds := image.Rect(int(a.x-a.size), int(a.y-a.size), int(a.x+a.size), int(a.y+a.size))

		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("crashed drawing circle at %v, %v with size %v\n", a.x, a.y, a.size)
			}
		}()
		a.canvas.DrawCircle(float64(a.x), float64(a.y), float64(a.size))
		a.canvas.SetColor(a.color)
		// a.canvas.SetRGB(0xFF, 0x66, 0)
		a.canvas.Stroke()
		//a.canvas.Fill()
		repaint <- repaintBounds

		a.x += a.directionX
		a.y += a.directionY

		time.Sleep(a.delay)
	}
}

func main() {
	windowWidth := 1800
	windowHeight := 1000

	driver.Main(func(s screen.Screen) {
		w, err := s.NewWindow(&screen.NewWindowOptions{
			Height: windowHeight,
			Width:  windowWidth,
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		defer w.Release()

		// Create a buffer to write to.
		var b screen.Buffer
		b, err = s.NewBuffer(image.Point{windowWidth, windowHeight})
		defer func() {
			if b != nil {
				b.Release()
			}
		}()

		dc := gg.NewContextForRGBA(b.RGBA())
		// dc.SavePNG("out.png")

		repaint := make(chan image.Rectangle)

		asteroids := make([]*Asteroid, 2)
		for i := 0; i < len(asteroids); i++ {
			size := rand.Intn(40) + 5
			directionX := rand.Intn(5) - rand.Intn(10)
			directionY := rand.Intn(5) - rand.Intn(10)
			delay := rand.Intn(200000) + 1
			randomColor := rand.Intn(215) + 1
			color := palette.WebSafe[randomColor]
			asteroids[i] = NewAsteroid(dc, size, directionX, directionY, time.Duration(delay), color)
			go asteroids[i].Start(repaint)
		}

		go painter(w, b, repaint)

		for {
			switch e := w.NextEvent().(type) {
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					return
				}
			case paint.Event:
				repaint <- b.Bounds()
			case error:
			}
		}
	})
}

func painter(w screen.Window, b screen.Buffer, repaint chan image.Rectangle) {
	paints := 0
	since := time.Now()

	timer := time.NewTimer(time.Second / 30)

	for {
		select {
		case bounds := <-repaint:
			w.Upload(bounds.Min, b, bounds)
			//w.Fill(bounds, color.RGBA{0x00, 0x7f, 0x00, 0x7f}, screen.Src)

		case <-timer.C:
			w.Publish()
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

func transformExample(w screen.Window, b screen.Buffer) {
	op := screen.Src

	//cos30 := math.Cos(math.Pi / 6)
	// sin30 := math.Sin(math.Pi / 6)
	src2dst := f64.Aff3{
		+0.5, -1.0, 200,
		+0.5, +1.0, 200,
	}
	w.DrawUniform(src2dst, color.RGBA{R: 0xFF, G: 0x00, B: 0x00, A: 0xFF},
		b.Bounds(), op, &screen.DrawOptions{})
}
