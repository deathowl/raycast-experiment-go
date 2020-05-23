package player

import (
	"math"

	"github.com/faiface/pixel"
)

func MoveLeft(s float64, world [][]int, pos *pixel.Vec, plane pixel.Vec) {
	if world[int(pos.X-plane.X*s)][int(pos.Y)] == 0 {
		pos.X -= plane.X * s
	}

	if world[int(pos.X)][int(pos.Y-plane.Y*s)] == 0 {
		pos.Y -= plane.Y * s
	}
}

func MoveBackwards(s float64, world [][]int, pos *pixel.Vec, plane, dir pixel.Vec) {
	if world[int(pos.X-dir.X*s)][int(pos.Y)] == 0 {
		pos.X -= dir.X * s
	}

	if world[int(pos.X)][int(pos.Y-dir.Y*s)] == 0 {
		pos.Y -= dir.Y * s
	}
}

func MoveRight(s float64, world [][]int, pos *pixel.Vec, plane pixel.Vec) {
	if world[int(pos.X+plane.X*s)][int(pos.Y)] == 0 {
		pos.X += plane.X * s
	}

	if world[int(pos.X)][int(pos.Y+plane.Y*s)] == 0 {
		pos.Y += plane.Y * s
	}
}
func MoveForward(s float64, world [][]int, pos *pixel.Vec, plane, dir pixel.Vec) {
	if true {
		if world[int(pos.X+dir.X*s)][int(pos.Y)] == 0 {
			pos.X += dir.X * s
		}

		if world[int(pos.X)][int(pos.Y+dir.Y*s)] == 0 {
			pos.Y += dir.Y * s
		}
	}
}

func LookHorizontal(s float64, dir, plane *pixel.Vec) {
	oldDirX := dir.X

	dir.X = dir.X*math.Cos(-s) - dir.Y*math.Sin(-s)
	dir.Y = oldDirX*math.Sin(-s) + dir.Y*math.Cos(-s)

	oldPlaneX := plane.X

	plane.X = plane.X*math.Cos(-s) - plane.Y*math.Sin(-s)
	plane.Y = oldPlaneX*math.Sin(-s) + plane.Y*math.Cos(-s)
}
