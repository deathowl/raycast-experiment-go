package assets

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
)

func LoadTextures() *image.RGBA {
	infile, err := os.Open("assets/tiles.png")
	if err != nil {
		// replace this with real error handling
		panic(err)
	}
	defer infile.Close()
	p, err := png.Decode(infile)
	if err != nil {
		panic(err)
	}

	m := image.NewRGBA(p.Bounds())
	draw.Draw(m, m.Bounds(), p, image.ZP, draw.Src)
	fmt.Println("Assets loaded successfully.")
	return m
}
