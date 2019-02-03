package main

import (
	"bytes"
	"encoding/json"
	"image/png"
	"io"
	"io/ioutil"
	"os"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/salviati/go-tmx/tmx"
)

var clearColor = colornames.Skyblue

var tilemap *tmx.Map
var sprites []*pixel.Sprite

// Tile maps a tilemap coordinate to be drawn at `GamePos`
type Tile struct {
	MapPos  pixel.Vec `json:"mapPos"`
	GamePos pixel.Vec `json:"gamePos"`
}

// Level represents a single game scene composed of tiles
//   - Tiles []*tile
type Level struct {
	Name  string  `json:"name"`
	Tiles []*Tile `json:"tiles"`
}

func gameloop(win *pixelgl.Window, level *Level) {
	tm := tilemap.Tilesets[0]
	w := float64(tm.TileWidth)
	h := float64(tm.TileHeight)
	sprite := loadSprite(tm.Image.Source)

	var iX, iY float64
	var fX = float64(tm.TileWidth)
	var fY = float64(tm.TileHeight)

	for !win.Closed() {
		win.Clear(clearColor)

		for _, coord := range level.Tiles {
			iX = coord.MapPos.X * w
			fX = iX + w
			iY = coord.MapPos.Y * h
			fY = iY + h
			sprite.Set(sprite.Picture(), pixel.R(iX, iY, fX, fY))
			pos := coord.GamePos.ScaledXY(pixel.V(w, h))
			sprite.Draw(win, pixel.IM.Moved(pos.Add(pixel.V(0, h))))
		}
		win.Update()
	}
}

func run() {
	// Create the window with OpenGL
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Tilemaps",
		Bounds: pixel.R(0, 0, 800, 600),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	panicIfErr(err)

	// Initialize art assets (i.e. the tilemap)
	tilemap, err = tmx.ReadFile("gameart2d-desert.tmx")
	panicIfErr(err)

	// Load the level from file
	level, err := ParseLevelFile("level.json")
	panicIfErr(err)

	gameloop(win, level)
}

func loadSprite(path string) *pixel.Sprite {
	f, err := os.Open(path)
	panicIfErr(err)

	img, err := png.Decode(f)
	panicIfErr(err)

	pd := pixel.PictureDataFromImage(img)
	return pixel.NewSprite(pd, pd.Bounds())
}

// ParseLevelFile reads a file from the disk at `path`
// and unmarshals it to a `*Level`
func ParseLevelFile(path string) (*Level, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var newLevel Level
	json.Unmarshal(bytes, &newLevel)
	return &newLevel, nil
}

// Save serializes a level to a JSON data file
func (level *Level) Save(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := json.MarshalIndent(level, "", "  ")
	if err != nil {
		return err
	}

	r := bytes.NewReader(b)
	_, err = io.Copy(f, r)
	return err
}

func main() {
	pixelgl.Run(run)
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
