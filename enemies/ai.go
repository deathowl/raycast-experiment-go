package enemies

import (
	"fmt"

	w "github.com/deathowl/raycast-experiment-go/world"
	"github.com/faiface/pixel"
)

func CanSee(x, y, pX, pY int) bool {
	fmt.Println(x, y)
	fmt.Println(pX, pY)
	if w.Enemies[x][y] != 0 {
		var startX, endX, dirX, startY, endY, dirY int
		if x <= pX {
			dirX = 1
			startX = x
			endX = pX
		} else {
			dirX = -1
			startX = pX
			endX = x
		}
		if y <= pY {
			dirY = 1
			startY = y
			endY = pY
		} else {

			dirY = -1
			startY = pY
			endY = y
		}
		step := 0
		stepY := startY
		for ex := startX; x != endX; ex += dirX {
			step++
			if (dirY == -1 && (stepY+(dirY*step)) > endY) || (dirY == 1 && (stepY+(dirY*step)) < endY) {
				stepY = stepY + (dirY * step)
			}
			fmt.Println(ex)
			fmt.Println(stepY)
			fmt.Println(w.World[ex][stepY])
			if w.World[ex][stepY] != 0 {
				fmt.Println("-----")
				return false
			}

		}
		fmt.Println("-----")

		return true
	}
	fmt.Println("-----")

	return false

}

func Tick(pos pixel.Vec) {
	newEnemies := make([][]int, len(w.Enemies))
	copy(newEnemies, w.Enemies)

	for x := 0; x < len(w.Enemies); x++ {
		newEnemies[x] = make([]int, len(w.Enemies[x]))
		for y := 0; y < len(w.Enemies[x]); y++ {
			if w.Enemies[x][y] != 0 {
				if CanSee(x, y, int(pos.X), int(pos.Y)) {
					newX := x
					newY := y
					if int(pos.X) > x {
						newX = x + 1
					}
					if int(pos.Y) > x {
						newY = y + 1
					}
					if int(pos.X) < y {
						newX = x - 1
					}
					if int(pos.Y) < y {
						newY = y - 1
					}
					if w.Enemies[newX][newY] == 0 {
						newEnemies[newX][newY] = w.Enemies[x][y]
						newEnemies[x][y] = 0
					}
				}
			}
		}
	}
	w.Enemies = newEnemies
	fmt.Println(newEnemies)
}
