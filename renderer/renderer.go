package renderer

import (
	"image"
	"image/color"
	"math"

	"github.com/deathowl/raycast-experiment-go/assets"
	w "github.com/deathowl/raycast-experiment-go/world"
	"github.com/faiface/pixel"
)

const (
	texSize = 64
)

var (
	textures = assets.LoadTextures()
)

func RenderFrame(width, height int, dir, pos, plane pixel.Vec) *image.RGBA {
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

			if w.World[worldX][worldY] > 0 {
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

		// if x == width/2 {
		// 	wallDistance = perpWallDist
		// }

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

		texNum := w.GetTexNum(worldX, worldY)

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

func RenderMinimap(pos pixel.Vec) *image.RGBA {
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
			c := w.GetColor(int(x), int(y))
			if c.A == 255 {
				c.A = 96
			}
			for rx := 0; rx <= 2; rx++ {
				for ry := 0; ry <= 2; ry++ {
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
