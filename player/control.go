package player

import "github.com/faiface/pixel"

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
func MoveForward(s float64, world [][]int, pos *pixel.Vec, plane, dir pixel.Vec, wallDistance float64) {
	if wallDistance > 0.3 {
		if world[int(pos.X+dir.X*s)][int(pos.Y)] == 0 {
			pos.X += dir.X * s
		}

		if world[int(pos.X)][int(pos.Y+dir.Y*s)] == 0 {
			pos.Y += dir.Y * s
		}
	}
}