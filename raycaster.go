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

	"github.com/deathowl/raycast-experiment-go/player"
	wgen "github.com/deathowl/raycast-experiment-go/world"
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

	textures = loadTextures()
	frames   = 0
	second   = time.Tick(time.Second)
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

var world = *wgen.GenWorld(1000)
var monsters = *wgen.GenEnemies(world, 10)

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

		drawStart := -lineHeight/2 + height/2

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
	active := pt.X > 0 && pt.X < len(world) && pt.Y > 0 && pt.Y < len(world[0])

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
			player.MoveForward(3.5*dt, world, &pos, plane, dir, wallDistance)
		}

		if win.Pressed(pixelgl.KeyLeft) || win.Pressed(pixelgl.KeyA) {
			player.MoveLeft(3.5*dt, world, &pos, plane)
		}

		if win.Pressed(pixelgl.KeyDown) || win.Pressed(pixelgl.KeyS) {
			player.MoveBackwards(3.5*dt, world, &pos, plane, dir)
		}

		if win.Pressed(pixelgl.KeyRight) || win.Pressed(pixelgl.KeyD) {
			player.MoveRight(3.5*dt, world, &pos, plane)
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

		if win.JustPressed(pixelgl.Key7) {
			as.set(7)
		}

		if win.JustPressed(pixelgl.Key0) {
			as.set(0)
		}
		if monsters[int(pos.X)][int(pos.Y)] == 1 {
			fmt.Println("OUCH")
		}

		if win.JustPressed(pixelgl.KeySpace) {
			as.toggle(3)
		}
		if win.JustPressed(pixelgl.MouseButtonRight) {
			as.execute()
		}
		p := pixel.PictureDataFromImage(frame())

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
