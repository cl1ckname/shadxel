package voxel

type Grid = [][][]Voxel

type VoxelGrid struct {
	Data Grid
	Size int
}

func NewVoxelGrid(size int) *VoxelGrid {
	data := make([][][]Voxel, size)
	for z := range data {
		data[z] = make([][]Voxel, size)
		for y := range data[z] {
			data[z][y] = make([]Voxel, size)
		}
	}
	return &VoxelGrid{Data: data, Size: size}
}

func (vg *VoxelGrid) At(x, y, z int) Voxel {
	return vg.Data[z][y][x]
}
