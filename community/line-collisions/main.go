package main

import (
	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

const (
	modeR = iota
	modeL
)

const (
	clickLineA = iota
	clickLineB
)

var (
	winBounds = pixel.R(0, 0, 1024, 768)

	r = pixel.R(10, 10, 70, 50)
	l = pixel.L(pixel.V(20, 20), pixel.V(100, 30))

	mode      = modeR
	clickLine = clickLineA
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Line collision",
		Bounds: winBounds,
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)

	for !win.Closed() {
		if mode == modeR {
			win.Clear(color.RGBA{R: 23, G: 39, B: 58, A: 125})
		} else {
			win.Clear(color.RGBA{R: 21, G: 55, B: 18, A: 125})
		}
		imd.Clear()

		if win.JustPressed(pixelgl.KeyR) {
			mode = modeR
		}
		if win.JustPressed(pixelgl.KeyL) {
			mode = modeL
			clickLine = clickLineA
		}

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			if mode == modeR {
				rectToMouse := r.Center().To(win.MousePosition())
				r = r.Moved(rectToMouse)
			}

			if mode == modeL {
				if clickLine == clickLineA {
					l = pixel.L(win.MousePosition(), win.MousePosition().Add(pixel.V(1, 1)))
					clickLine = clickLineB
				} else {
					l = pixel.L(l.A, win.MousePosition())
					clickLine = clickLineA
				}
			}
		}

		imd.Color = color.Black
		imd.Push(r.Min, r.Max)
		imd.Rectangle(3)

		imd.Color = color.RGBA{R: 10, G: 10, B: 250, A: 255}
		imd.Push(l.A, l.B)
		imd.Line(3)

		imd.Color = color.RGBA{R: 250, G: 10, B: 10, A: 255}
		for _, i := range r.IntersectionPoints(l) {
			imd.Push(i)
			imd.Circle(4, 0)
		}

		imd.Draw(win)
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
