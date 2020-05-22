package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"time"

	"github.com/deathowl/raycast-experiment-go/player"
	"github.com/deathowl/raycast-experiment-go/renderer"
	"github.com/deathowl/raycast-experiment-go/world"
	w "github.com/deathowl/raycast-experiment-go/world"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	texSize = 64
)

var (
	fullscreen   = false
	showMap      = true
	width        = 960
	height       = 540
	scale        = 2.0
	wallDistance = 3.0

	as       actionSquare
	inv      Inventory
	selected int

	pos, dir, plane pixel.Vec

	frames = 0
	second = time.Tick(time.Second)
)

//Typedfs

type Inventory struct {
	items []string
}

type actionSquare struct {
	X      int
	Y      int
	block  int
	active bool
}

type Entity struct {
	Health  int
	texture int
}

func setup() {
	pos = pixel.V(12.0, 14.5)
	dir = pixel.V(-1.0, 0.0)
	plane = pixel.V(0.0, 0.66)
	world.InitWorld(1000)
}

func getColor(x, y int) color.RGBA {
	switch w.World[x][y] {
	case 0:
		return color.RGBA{43, 30, 24, 255}
	case 1:
		return color.RGBA{100, 89, 73, 255}
	case 2:
		return color.RGBA{110, 23, 0, 255}
	case 3:
		return color.RGBA{45, 103, 171, 255}
	case 4:
		return color.RGBA{123, 84, 33, 255}
	case 5:
		return color.RGBA{158, 148, 130, 255}
	case 6:
		return color.RGBA{203, 161, 47, 255}
	case 7:
		return color.RGBA{255, 107, 0, 255}
	case 9:
		return color.RGBA{0, 0, 0, 0}
	default:
		return color.RGBA{255, 194, 32, 255}
	}
}

func minimap() *image.RGBA {
	m := image.NewRGBA(image.Rect(0, 0, 60, 65))
	startx := math.Max(0., pos.X-float64(12))
	starty := math.Max(0., pos.Y-float64(13))
	endx := 24.
	endy := 26.
	if startx != 0 {
		endx = pos.X + float64(12)
	}
	if starty != 0 {
		endy = pos.Y + float64(13)
	}

	for x := startx; x < endx; x++ {
		for y := starty; y < endy; y++ {
			c := getColor(int(x), int(y))
			if c.A == 255 {
				c.A = 96
			}
			for rx := 0; rx <= 2; rx++ {
				for ry := 0; ry <= 3; ry++ {
					m.Set(2*int(x-startx)+int(x-startx)+rx, 2*int(y-starty)+int(y-starty)+ry, c)
					m.Set(2*int(x-startx)+int(x-startx)-rx, 2*int(y-starty)+int(y-starty)-ry, c)
					m.Set(2*int(x-startx)+int(x-startx)+rx, 2*int(y-starty)+int(y-starty)-ry, c)
					m.Set(2*int(x-startx)+int(x-startx)-rx, 2*int(y-starty)+int(y-starty)+ry, c)
				}
			}
		}
	}

	for rx := 0; rx <= 1; rx++ {
		for ry := 0; ry <= 1; ry++ {
			m.Set(2*int(pos.X-startx)+int(pos.X-startx)+rx, 2*int(pos.Y-starty)+int(pos.Y-starty)+ry, color.RGBA{0, 255, 0, 180})
			m.Set(2*int(pos.X-startx)+int(pos.X-startx)-rx, 2*int(pos.Y-starty)+int(pos.Y-starty)-ry, color.RGBA{0, 255, 0, 180})
			m.Set(2*int(pos.X-startx)+int(pos.X-startx)+rx, 2*int(pos.Y-starty)+int(pos.Y-starty)-ry, color.RGBA{0, 255, 0, 180})
			m.Set(2*int(pos.X-startx)+int(pos.X-startx)-rx, 2*int(pos.Y-starty)+int(pos.Y-starty)+ry, color.RGBA{0, 255, 0, 180})
		}
	}
	return m
}

func getInventory() Inventory {

	return Inventory{}

}

func (i Inventory) getBitMap() *image.RGBA {
	m := image.NewRGBA(image.Rect(10, 10, 8, 12))
	mygreen := color.RGBA{0, 0, 0, 255} //  R, G, B, Alpha
	draw.Draw(m, m.Bounds(), &image.Uniform{mygreen}, image.ZP, draw.Src)

	return m

}

func getActionSquare() actionSquare {
	pt := image.Pt(int(pos.X)+1, int(pos.Y))

	a := dir.Angle()

	switch {
	case a > 2.8 || a < -2.8:
		pt = image.Pt(int(pos.X)-1, int(pos.Y))
	case a > -2.8 && a < -2.2:
		pt = image.Pt(int(pos.X)-1, int(pos.Y)-1)
	case a > -2.2 && a < -1.4:
		pt = image.Pt(int(pos.X), int(pos.Y)-1)
	case a > -1.4 && a < -0.7:
		pt = image.Pt(int(pos.X)+1, int(pos.Y)-1)
	case a > 0.4 && a < 1.0:
		pt = image.Pt(int(pos.X)+1, int(pos.Y)+1)
	case a > 1.0 && a < 1.7:
		pt = image.Pt(int(pos.X), int(pos.Y)+1)
	case a > 1.7:
		pt = image.Pt(int(pos.X)-1, int(pos.Y)+1)
	}
	block := -1
	active := pt.X > 0 && pt.X < len(w.World) && pt.Y > 0 && pt.Y < len(w.World[0])

	if active {
		block = w.World[pt.X][pt.Y]
	}

	return actionSquare{
		X:      pt.X,
		Y:      pt.Y,
		active: active,
		block:  block,
	}
}

func (as actionSquare) toggle(n int) {
	if as.active {
		if w.World[as.X][as.Y] == 0 {
			w.World[as.X][as.Y] = n
		} else {
			w.World[as.X][as.Y] = 0
		}
	}
}

func (as actionSquare) set(n int) {
	selected = n
}

func (as actionSquare) execute() {
	if as.active {
		w.World[as.X][as.Y] = selected
	}
}

func run() {
	cfg := pixelgl.WindowConfig{
		Bounds:      pixel.R(0, 0, float64(width)*scale, float64(height)*scale),
		VSync:       true,
		Undecorated: false,
	}

	if fullscreen {
		cfg.Monitor = pixelgl.PrimaryMonitor()
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	c := win.Bounds().Center()

	last := time.Now()
	minimaprefresh := time.Now()
	minimapInit := true
	var lastM *pixel.PictureData

	mapRot := -1.6683362599999894

	for !win.Closed() {
		if win.JustPressed(pixelgl.KeyEscape) || win.JustPressed(pixelgl.KeyQ) {
			return
		}

		win.Clear(color.Black)

		dt := time.Since(last).Seconds()
		last = time.Now()

		as = getActionSquare()

		if win.Pressed(pixelgl.KeyUp) || win.Pressed(pixelgl.KeyW) {
			player.MoveForward(3.5*dt, w.World, &pos, plane, dir, wallDistance)
		}

		if win.Pressed(pixelgl.KeyLeft) || win.Pressed(pixelgl.KeyA) {
			player.MoveLeft(3.5*dt, w.World, &pos, plane)
		}

		if win.Pressed(pixelgl.KeyDown) || win.Pressed(pixelgl.KeyS) {
			player.MoveBackwards(3.5*dt, w.World, &pos, plane, dir)
		}

		if win.Pressed(pixelgl.KeyRight) || win.Pressed(pixelgl.KeyD) {
			player.MoveRight(3.5*dt, w.World, &pos, plane)
		}
		movedX := win.MousePosition().X - win.MousePreviousPosition().X

		if movedX > 0 {
			player.LookHorizontal(movedX*dt*0.5, &dir, &plane)
		}
		if movedX < 0 {
			player.LookHorizontal(movedX*dt*0.5, &dir, &plane)
		}

		if win.JustPressed(pixelgl.KeyM) {
			showMap = !showMap
		}
		if win.JustPressed(pixelgl.Key0) {
			as.set(0)
		}
		if win.JustPressed(pixelgl.Key1) {
			as.set(1)
		}

		if win.JustPressed(pixelgl.Key2) {
			as.set(2)
		}

		if win.JustPressed(pixelgl.Key3) {
			as.set(3)
		}

		if win.JustPressed(pixelgl.Key4) {
			as.set(4)
		}

		if win.JustPressed(pixelgl.Key5) {
			as.set(5)
		}

		if win.JustPressed(pixelgl.Key6) {
			as.set(6)
		}

		if win.JustPressed(pixelgl.Key0) {
			as.set(0)
		}
		if w.Enemies[int(pos.X)][int(pos.Y)] == 1 {
			fmt.Println("OUCH")
		}

		if win.JustPressed(pixelgl.KeySpace) {
			as.toggle(3)
		}
		if win.JustPressed(pixelgl.MouseButtonRight) {
			as.execute()
		}
		p := pixel.PictureDataFromImage(renderer.RenderFrame(width, height, dir, pos, plane))

		pixel.NewSprite(p, p.Bounds()).
			Draw(win, pixel.IM.Moved(c).Scaled(c, scale))

		if showMap {
			dt := time.Since(minimaprefresh).Seconds()
			m := pixel.PictureDataFromImage(minimap())
			if minimapInit {
				minimapInit = false
				lastM = m
			}
			if dt > .5 {
				lastM = m
				minimaprefresh = time.Now()
			}
			mc := m.Bounds().Min.Add(pixel.V(-lastM.Rect.W(), lastM.Rect.H()))

			pixel.NewSprite(lastM, lastM.Bounds()).
				Draw(win, pixel.IM.
					Moved(mc).
					Rotated(mc, mapRot).
					ScaledXY(pixel.ZV, pixel.V(-scale*2, scale*2)))
		}

		i := pixel.PictureDataFromImage(inv.getBitMap())
		ic := i.Bounds().Min.Add(pixel.V(-i.Rect.W(), i.Rect.H()))
		pixel.NewSprite(i, i.Bounds()).
			Draw(win, pixel.IM.
				Moved(ic).
				Rotated(ic, mapRot).
				ScaledXY(pixel.ZV, pixel.V(-scale*2, scale*2)))
		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("ShootingGame | FPS: %d", frames))
			fmt.Println("ShootingGame | FPS: ", frames)
			frames = 0
		default:
		}
		win.Update()
	}
}

func main() {
	flag.BoolVar(&fullscreen, "f", fullscreen, "fullscreen")
	flag.IntVar(&width, "w", width, "width")
	flag.IntVar(&height, "h", height, "height")
	flag.Float64Var(&scale, "s", scale, "scale")
	flag.Parse()
	fmt.Println("Done reading map.")
	fmt.Println("Done setting up input handler")
	fmt.Println("Window created")

	setup()
	pixelgl.Run(run)
}
