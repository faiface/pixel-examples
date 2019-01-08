package main

import (
	"time"
	"log"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

import "flag"
import "os"

var (
	version       string
	race          bool
	debug         = os.Getenv("BUILDDEBUG") != ""
	filename      string
	width         int
	height        int
	timeout       = "120s"
	uDrift        float32
)

func run() {
	// Set up window configs
	cfg := pixelgl.WindowConfig{ // Default: 1024 x 768
		Title:  "Golang GLSL",
		Bounds: pixel.R(0, 0, float64(width), float64(height)),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	camVector := win.Bounds().Center()

	bounds := win.Bounds()
	bounds.Max = bounds.Max.ScaledXY(pixel.V(1.0, 1.0))

	// I am putting all shader example initializing stuff here for
	// easier reference to those learning to use this functionality

	fragSource, err := LoadFileToString(filename)

	if err != nil {
		panic(err)
	}

	var uMouse mgl32.Vec4
	var uTime float32

	canvas := win.Canvas()
	uResolution := mgl32.Vec2{float32(win.Bounds().W()), float32(win.Bounds().H())}

	EasyBindUniforms(canvas,
		"uResolution", &uResolution,
		"uTime", &uTime,
		"uMouse", &uMouse,
		"uDrift", &uDrift,
	)

	canvas.SetFragmentShader(fragSource)

	start := time.Now()

	// Game Loop
	for !win.Closed() {
		uTime = float32(time.Since(start).Seconds())
		mpos := win.MousePosition()
		uMouse[0] = float32(mpos.X)
		uMouse[1] = float32(mpos.Y)

		win.Clear(colornames.Black)

		// Drawing to the screen
		canvas.Draw(win, pixel.IM.Moved(camVector))

		win.Update()
	}

}

func parseFlags() {
	flag.StringVar  (&version,       "version",       "v0.1",                     "Set compiled in version string")
	flag.StringVar  (&filename,      "filename",      "shaders/seascape.glsl",    "path to GLSL file")
	flag.IntVar     (&width,         "width",         1024,                       "Width of the OpenGL Window")
	flag.IntVar     (&height,        "height",        768,                        "Height of the OpenGL Window")
	var tmp float64
	flag.Float64Var (&tmp,           "drift",         0.01,                       "Speed of the gradual camera drift")
	flag.BoolVar    (&race,          "race",          race,                       "Use race detector")

	// this parses the arguements
	flag.Parse()

	uDrift = float32(tmp)
	log.Println("width=",width)
	log.Println("height=",height)
	log.Println("uDrift=",uDrift)

}

func main() {
	parseFlags()

	pixelgl.Run(run)
}
