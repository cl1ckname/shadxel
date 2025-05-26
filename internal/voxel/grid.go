package voxel

type Grid = [][][]Voxel

type VoxelGrid struct {
	Data Grid
	Mask [][][]bool
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
	mask := make([][][]bool, size)
	for z := range data {
		mask[z] = make([][]bool, size)
		for y := range mask[z] {
			mask[z][y] = make([]bool, size)
		}
	}
	return &VoxelGrid{Data: data, Size: size, Mask: mask}
}

func (v *VoxelGrid) Set(g Grid) {
	v.Data = g
	for z := range g {
		for y := range g[z] {
			for x := range g[z][y] {
				c := v.At(x, y, z)
				if !c.Visible() {
					continue
				}

				for _, dir := range dirs {
					nx, ny, nz := x+dir.dx, y+dir.dy, z+dir.dz
					if !v.InBounds(nx, ny, nz) {
						v.Mask[z][y][x] = true
						continue
					}
					if v.At(nx, ny, nz).Visible() {
						continue
					}
					v.Mask[z][y][x] = true
				}

			}
		}
	}
}

func (vg *VoxelGrid) At(x, y, z int) Voxel {
	return vg.Data[z][y][x]
}

func (vg *VoxelGrid) Hit(x, y, z int) bool {
	return vg.Mask[z][y][x]
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

var dirs = []struct {
	dx, dy, dz int
}{
	{1, 0, 0},
	{0, 1, 0},
	{0, 0, 1},
	{-1, 0, 0},
	{0, -1, 0},
	{0, 0, -1},
}
