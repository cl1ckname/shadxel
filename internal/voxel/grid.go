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

func (vg *VoxelGrid) InBounds(x, y, z int) bool {
	if x < 0 || x >= vg.Size {
		return false
	}
	if y < 0 || y >= vg.Size {
		return false
	}
	if z < 0 || z >= vg.Size {
		return false
	}
	return true
}
