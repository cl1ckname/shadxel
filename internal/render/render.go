package render

import (
	"shadxel/internal/voxel"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Renderer struct {
	shader ShaderProgram
	mesh   *CubeMesh
	scale  float32
}

// Public constructor
func NewRenderer() (*Renderer, error) {
	shader, err := NewShaderProgram()
	if err != nil {
		return nil, err
	}
	r := &Renderer{
		scale:  1. / 25, // adjust based on grid size
		shader: *shader,
		mesh:   NewCubeMesh(),
	}

	return r, nil
}

func (r *Renderer) Draw(grid voxel.VoxelGrid, view, projection mgl32.Mat4) {
	gl.UseProgram(r.shader.ID)

	// Set view and projection once per frame
	viewLoc := gl.GetUniformLocation(r.shader.ID, gl.Str("view\x00"))
	gl.UniformMatrix4fv(viewLoc, 1, false, &view[0])

	projLoc := gl.GetUniformLocation(r.shader.ID, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projLoc, 1, false, &projection[0])

	scale := r.scale
	center := float32(grid.Size) / 2

	for z := range grid.Data {
		for y := range grid.Data[z] {
			for x := range grid.Data[z][y] {
				c := grid.At(x, y, z)
				if c.R == 0 && c.G == 0 && c.B == 0 {
					continue
				}

				color := mgl32.Vec3{
					float32(c.R) / 255,
					float32(c.G) / 255,
					float32(c.B) / 255,
				}

				pos := mgl32.Vec3{
					(float32(x) - center) * scale,
					(float32(y) - center) * scale,
					(float32(z) - center) * scale,
				}
				model := mgl32.Translate3D(pos.X(), pos.Y(), pos.Z()).Mul4(
					mgl32.Scale3D(scale, scale, scale),
				)

				r.mesh.DrawAt(r.shader, color, model)
			}
		}
	}
}
