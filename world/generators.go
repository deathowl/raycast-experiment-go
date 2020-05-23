package world

import (
	"math"
	"math/rand"
	"time"

	"github.com/aquilax/go-perlin"
	wr "github.com/mroth/weightedrand"
)

const (
	alpha = 2.
	beta  = 2.
	n     = 3
)

func genWorld(size int64) *[][]int {
	rand.Seed(time.Now().UTC().UnixNano())
	world := make([][]int, size)
	p := perlin.NewPerlin(alpha, beta, n, rand.Int63n(100))
	caveMaterial := -1
	for x := 0.; x <= float64(size-int64(1)); x++ {
		world[int(x)] = make([]int, size)
		for y := 0.; y <= float64(size-int64(1)); y++ {
			intx := int64(x)
			inty := int64(y)
			if intx == 0 || inty == 0 || intx == size-int64(1) || inty == size-int64(1) {
				//borders
				world[intx][inty] = 1
			} else if (intx > size/2-int64(2) || intx > size/2+int64(2)) && (inty < size/2-int64(2) || inty > size/2+int64(2)) {
				//avoid spawning in blocks
				world[intx][inty] = 0
			} else {
				//caves
				noise := p.Noise2D(x/float64(10), y/float64(10))

				if math.Abs(noise) >= 0.2 {
					if caveMaterial == -1 {
						c := wr.NewChooser(
							wr.Choice{Item: 2, Weight: 1},
							wr.Choice{Item: 3, Weight: 1},
							wr.Choice{Item: 4, Weight: 1},
							wr.Choice{Item: 5, Weight: 1},
							wr.Choice{Item: 6, Weight: 1},
							wr.Choice{Item: 7, Weight: 1},
							//wr.Choice{Item: 8, Weight: 1000},

						)
						caveMaterial = c.Pick().(int)

					}
					world[intx][inty] = caveMaterial
				} else {
					//look ahead
					// noisex := p.Noise2D(x+1/float64(10), y/float64(10))
					// noisey := p.Noise2D(x/float64(10), y/float64(10))
					// noisexy := p.Noise2D(x+1/float64(10), y+1/float64(10))
					// if math.Abs(noisex) < .2 &&  math.Abs(noisey) < .2 &&  math.Abs(noisexy) < .2{
					//caveMaterial = -1
					// }
					world[intx][inty] = 0
				}
			}

		}
	}
	return &world
}

func genEnemies(world [][]int, probability int) *[][]int {
	rand.Seed(time.Now().UTC().UnixNano())
	monsters := make([][]int, len(world))
	for x := 0; x < len(world); x++ {
		monsters[x] = make([]int, len(world))
		for y := 0; y < len(world[x]); y++ {
			if world[x][y] == 0 {
				monsters[x][y] = 0
				roll := rand.Intn(100)
				if roll < probability {
					monsters[x][y] = 1
				}
			}
		}
	}
	return &monsters
}
