package voxel

type Grid = [][]Color

type Color struct {
	R, G, B byte
}

func GenerateGrid(width, height int) [][]Color {
	grid := make([][]Color, height)
	for y := 0; y < height; y++ {
		row := make([]Color, width)
		for x := 0; x < width; x++ {
			row[x] = Color{
				R: 255 - byte(x*5),
				G: byte(y * 5),
				B: byte(100),
			}
		}
		grid[y] = row
	}
	return grid
}
