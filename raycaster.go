package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const texSize = 64

var (
	fullscreen   = false
	showMap      = true
	width        = 320
	height       = 200
	scale        = 3.0
	wallDistance = 8.0

	as       actionSquare
	inv      Inventory
	selected int

	pos, dir, plane pixel.Vec

	textures = loadTextures()
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
}

var world = [24][24]int{
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 2, 2, 2, 2, 2, 2, 2, 0, 0, 0, 0, 3, 0, 3, 0, 3, 0, 0, 0, 1},
	{1, 0, 0, 0, 2, 7, 2, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 2, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 3, 0, 7, 0, 3, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 2, 2, 2, 2, 0, 2, 2, 0, 0, 0, 0, 3, 0, 3, 0, 3, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 0, 6, 0, 4, 0, 0, 0, 4, 0, 0, 0, 0, 5, 0, 0, 0, 0, 0, 5, 0, 0, 0, 1},
	{1, 0, 6, 0, 4, 0, 7, 0, 4, 0, 0, 0, 0, 0, 5, 0, 0, 0, 5, 0, 0, 0, 0, 1},
	{1, 0, 0, 0, 4, 0, 0, 0, 4, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5, 5, 0, 0, 0, 1},
	{1, 4, 4, 4, 4, 4, 4, 0, 4, 0, 0, 0, 5, 5, 0, 5, 5, 5, 0, 5, 5, 0, 0, 1},
	{1, 4, 0, 0, 0, 0, 0, 0, 4, 0, 0, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 0, 1},
	{1, 4, 0, 4, 0, 0, 0, 0, 4, 0, 0, 5, 0, 5, 5, 5, 5, 5, 5, 5, 0, 5, 0, 1},
	{1, 4, 0, 4, 4, 4, 4, 4, 4, 0, 0, 5, 0, 5, 0, 0, 0, 0, 0, 5, 0, 5, 0, 1},
	{1, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 5, 0, 5, 5, 0, 0, 0, 0, 1},
	{1, 4, 4, 4, 4, 4, 4, 4, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
}

func loadTextures() *image.RGBA {
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

func getTexNum(x, y int) int {
	return world[x][y]
}

func getColor(x, y int) color.RGBA {
	switch world[x][y] {
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

func frame() *image.RGBA {
	m := image.NewRGBA(image.Rect(0, 0, width, height))

	for x := 0; x < width; x++ {
		var (
			step         image.Point
			sideDist     pixel.Vec
			perpWallDist float64
			hit, side    bool

			rayPos, worldX, worldY = pos, int(pos.X), int(pos.Y)

			cameraX = 2*float64(x)/float64(width) - 1

			rayDir = pixel.V(
				dir.X+plane.X*cameraX,
				dir.Y+plane.Y*cameraX,
			)

			deltaDist = pixel.V(
				math.Sqrt(1.0+(rayDir.Y*rayDir.Y)/(rayDir.X*rayDir.X)),
				math.Sqrt(1.0+(rayDir.X*rayDir.X)/(rayDir.Y*rayDir.Y)),
			)
		)

		if rayDir.X < 0 {
			step.X = -1
			sideDist.X = (rayPos.X - float64(worldX)) * deltaDist.X
		} else {
			step.X = 1
			sideDist.X = (float64(worldX) + 1.0 - rayPos.X) * deltaDist.X
		}

		if rayDir.Y < 0 {
			step.Y = -1
			sideDist.Y = (rayPos.Y - float64(worldY)) * deltaDist.Y
		} else {
			step.Y = 1
			sideDist.Y = (float64(worldY) + 1.0 - rayPos.Y) * deltaDist.Y
		}

		for !hit {
			if sideDist.X < sideDist.Y {
				sideDist.X += deltaDist.X
				worldX += step.X
				side = false
			} else {
				sideDist.Y += deltaDist.Y
				worldY += step.Y
				side = true
			}

			if world[worldX][worldY] > 0 {
				hit = true
			}
		}

		var wallX float64

		if side {
			perpWallDist = (float64(worldY) - rayPos.Y + (1-float64(step.Y))/2) / rayDir.Y
			wallX = rayPos.X + perpWallDist*rayDir.X
		} else {
			perpWallDist = (float64(worldX) - rayPos.X + (1-float64(step.X))/2) / rayDir.X
			wallX = rayPos.Y + perpWallDist*rayDir.Y
		}

		if x == width/2 {
			wallDistance = perpWallDist
		}

		wallX -= math.Floor(wallX)

		texX := int(wallX * float64(texSize))

		lineHeight := int(float64(height) / perpWallDist)

		if lineHeight < 1 {
			lineHeight = 1
		}

		drawStart := -lineHeight/2 + height/2
		if drawStart < 0 {
			drawStart = 0
		}

		drawEnd := lineHeight/2 + height/2
		if drawEnd >= height {
			drawEnd = height - 1
		}

		if !side && rayDir.X > 0 {
			texX = texSize - texX - 1
		}

		if side && rayDir.Y < 0 {
			texX = texSize - texX - 1
		}

		texNum := getTexNum(worldX, worldY)

		for y := drawStart; y < drawEnd+1; y++ {
			d := y*256 - height*128 + lineHeight*128
			texY := ((d * texSize) / lineHeight) / 256

			c := textures.RGBAAt(
				texX+texSize*(texNum),
				texY%texSize,
			)

			if side {
				c.R = c.R / 2
				c.G = c.G / 2
				c.B = c.B / 2
			}

			m.Set(x, y, c)
		}

		var floorWall pixel.Vec

		if !side && rayDir.X > 0 {
			floorWall.X = float64(worldX)
			floorWall.Y = float64(worldY) + wallX
		} else if !side && rayDir.X < 0 {
			floorWall.X = float64(worldX) + 1.0
			floorWall.Y = float64(worldY) + wallX
		} else if side && rayDir.Y > 0 {
			floorWall.X = float64(worldX) + wallX
			floorWall.Y = float64(worldY)
		} else {
			floorWall.X = float64(worldX) + wallX
			floorWall.Y = float64(worldY) + 1.0
		}

		distWall, distPlayer := perpWallDist, 0.0

		for y := drawEnd + 1; y < height; y++ {
			currentDist := float64(height) / (2.0*float64(y) - float64(height))

			weight := (currentDist - distPlayer) / (distWall - distPlayer)

			currentFloor := pixel.V(
				weight*floorWall.X+(1.0-weight)*pos.X,
				weight*floorWall.Y+(1.0-weight)*pos.Y,
			)

			fx := int(currentFloor.X*float64(texSize)) % texSize
			fy := int(currentFloor.Y*float64(texSize)) % texSize

			m.Set(x, y, textures.At(fx, fy))

			m.Set(x, height-y-1, textures.At(fx+(4*texSize), fy))
			m.Set(x, height-y, textures.At(fx+(4*texSize), fy))
		}
	}

	cursor := textures.RGBAAt(
		200,
		200,
	)
	cursor.R = 30
	cursor.G = 30
	cursor.B = 30

	for i := 2; i < 5; i++ {
		m.Set(width/2-i, height/2, cursor)
		m.Set(width/2+i, height/2, cursor)
		m.Set(width/2, height/2+i, cursor)
		m.Set(width/2, height/2-i, cursor)

	}

	return m
}

func minimap() *image.RGBA {
	m := image.NewRGBA(image.Rect(0, 0, 24, 26))

	for x, row := range world {
		for y, _ := range row {
			c := getColor(x, y)
			if c.A == 255 {
				c.A = 96
			}
			m.Set(x, y, c)
		}
	}

	m.Set(int(pos.X), int(pos.Y), color.RGBA{255, 0, 0, 255})

	if as.active {
		m.Set(as.X, as.Y, color.RGBA{255, 255, 255, 255})
	} else {
		m.Set(as.X, as.Y, color.RGBA{64, 64, 64, 255})
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
	active := pt.X > 0 && pt.X < 23 && pt.Y > 0 && pt.Y < 23

	if active {
		block = world[pt.X][pt.Y]
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
		if world[as.X][as.Y] == 0 {
			world[as.X][as.Y] = n
		} else {
			world[as.X][as.Y] = 0
		}
	}
}

func (as actionSquare) set(n int) {
	selected = n
}

func (as actionSquare) execute() {
	if as.active {
		world[as.X][as.Y] = selected
	}
}

func run() {
	cfg := pixelgl.WindowConfig{
		Bounds:      pixel.R(0, 0, float64(width)*scale, float64(height)*scale),
		VSync:       true,
		Undecorated: true,
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
			moveForward(3.5 * dt)
		}

		if win.Pressed(pixelgl.KeyA) {
			moveLeft(3.5 * dt)
		}

		if win.Pressed(pixelgl.KeyDown) || win.Pressed(pixelgl.KeyS) {
			moveBackwards(3.5 * dt)
		}

		if win.Pressed(pixelgl.KeyD) {
			moveRight(3.5 * dt)
		}

		movedX := win.MousePosition().X - win.MousePreviousPosition().X

		if movedX > 0 {
			lookRight(movedX * dt * 0.5)
		}
		if movedX < 0 {
			lookLeft(movedX * dt * -0.5)
		}

		/*
		movedY := win.MousePosition().X - win.MousePreviousPosition().X

		if movedY > 0 {
			lookUp(movedY * dt * 0.5)
		}
		if movedY < 0 {
			lookDown(movedY * dt * -0.5)
		}*/

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

		if win.JustPressed(pixelgl.Key7) {
			as.set(7)
		}

		if win.JustPressed(pixelgl.Key0) {
			as.set(0)
		}

		if win.JustPressed(pixelgl.KeySpace) {
			as.toggle(3)
		}
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			as.execute()
		}
		p := pixel.PictureDataFromImage(frame())

		pixel.NewSprite(p, p.Bounds()).
			Draw(win, pixel.IM.Moved(c).Scaled(c, scale))

		if showMap {
			m := pixel.PictureDataFromImage(minimap())

			mc := m.Bounds().Min.Add(pixel.V(-m.Rect.W(), m.Rect.H()))

			pixel.NewSprite(m, m.Bounds()).
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

		win.Update()
	}
}

func moveForward(s float64) {
	if wallDistance > 0.3 {
		if world[int(pos.X+dir.X*s)][int(pos.Y)] == 0 {
			pos.X += dir.X * s
		}

		if world[int(pos.X)][int(pos.Y+dir.Y*s)] == 0 {
			pos.Y += dir.Y * s
		}
	}
}

func moveLeft(s float64) {
	if world[int(pos.X-plane.X*s)][int(pos.Y)] == 0 {
		pos.X -= plane.X * s
	}

	if world[int(pos.X)][int(pos.Y-plane.Y*s)] == 0 {
		pos.Y -= plane.Y * s
	}
}

func moveBackwards(s float64) {
	if world[int(pos.X-dir.X*s)][int(pos.Y)] == 0 {
		pos.X -= dir.X * s
	}

	if world[int(pos.X)][int(pos.Y-dir.Y*s)] == 0 {
		pos.Y -= dir.Y * s
	}
}

func moveRight(s float64) {
	if world[int(pos.X+plane.X*s)][int(pos.Y)] == 0 {
		pos.X += plane.X * s
	}

	if world[int(pos.X)][int(pos.Y+plane.Y*s)] == 0 {
		pos.Y += plane.Y * s
	}
}

func lookRight(s float64) {
	oldDirX := dir.X

	dir.X = dir.X*math.Cos(-s) - dir.Y*math.Sin(-s)
	dir.Y = oldDirX*math.Sin(-s) + dir.Y*math.Cos(-s)

	oldPlaneX := plane.X

	plane.X = plane.X*math.Cos(-s) - plane.Y*math.Sin(-s)
	plane.Y = oldPlaneX*math.Sin(-s) + plane.Y*math.Cos(-s)
}

func lookLeft(s float64) {
	oldDirX := dir.X

	dir.X = dir.X*math.Cos(s) - dir.Y*math.Sin(s)
	dir.Y = oldDirX*math.Sin(s) + dir.Y*math.Cos(s)

	oldPlaneX := plane.X

	plane.X = plane.X*math.Cos(s) - plane.Y*math.Sin(s)
	plane.Y = oldPlaneX*math.Sin(s) + plane.Y*math.Cos(s)
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
