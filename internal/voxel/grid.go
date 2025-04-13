package voxel

type Grid = [][][]Color

type Color struct {
	R, G, B byte
}

type VoxelGrid struct {
	Data [][][]Color
	Size int
}

func NewVoxelGrid(size int) *VoxelGrid {
	data := make([][][]Color, size)
	for z := range data {
		data[z] = make([][]Color, size)
		for y := range data[z] {
			data[z][y] = make([]Color, size)
		}
	}
	return &VoxelGrid{Data: data, Size: size}
}

func (vg *VoxelGrid) At(x, y, z int) Color {
	return vg.Data[z][y][x]
}
