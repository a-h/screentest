package main

import (
	"fmt"
	"image"

	"time"

	"github.com/fogleman/gg"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
)

func main() {
	windowWidth := 1200
	windowHeight := 800

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

		repaint := make(chan bool)
		go draw(dc, repaint)
		go painter(w, b, repaint)

		for {
			switch e := w.NextEvent().(type) {
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					return
				}
			case paint.Event:
				repaint <- true
			case error:
			}
		}
	})
}

func draw(dc *gg.Context, repaint chan bool) {
	for i := 0; i < 1000; i++ {
		x := float64(i)
		y := float64(i)
		dc.DrawCircle(x, y, 400)
		dc.SetRGB(0xFF, 0x66, 0)
		dc.Fill()
		repaint <- true
		time.Sleep(10)
	}
}

func painter(w screen.Window, b screen.Buffer, repaint chan bool) {
	for {
		<-repaint
		// This just repaints the whole screen.
		// There will be a way to make this better.
		// Maybe the channel should state which area should be redrawn?
		w.Upload(image.Point{}, b, b.Bounds())
		w.Publish()
	}
}
