package world

var World, Enemies [][]int

func InitWorld(size int64) {
	World = *genWorld(size)
	Enemies = *genEnemies(World, 10)

}

func GetTexNum(x, y int) int {
	return World[x][y]
}
