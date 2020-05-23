package world

import (
	"image/color"
)

var World, Enemies [][]int

func InitWorld(size int64) {
	World = *genWorld(size)
	Enemies = *genEnemies(World, 10)

}

// func LoadWorld(path string) {
// 	f, err := os.Open("/tmp/dat")

// }

func GetTexNum(x, y int) int {
	return World[x][y]
}

func GetColor(x, y int) color.RGBA {
	switch GetTexNum(x, y) {
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
