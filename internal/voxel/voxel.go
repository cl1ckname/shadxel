package voxel

type Voxel struct {
	r, g, b byte
	visible bool
}

func NewVoxel(r, g, b byte, v bool) Voxel {
	return Voxel{
		r:       r,
		g:       g,
		b:       b,
		visible: v,
	}
}

func (v Voxel) R() uint8 {
	return v.r
}

func (v Voxel) G() uint8 {
	return v.g
}

func (v Voxel) B() uint8 {
	return v.b
}

func (v Voxel) Visible() bool {
	return v.visible
}
