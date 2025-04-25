package render

import (
	"shadxel/internal/render/mesh"
	"shadxel/internal/render/shader"
	"shadxel/internal/voxel"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Renderer struct {
	shader     shader.Program
	cube       *mesh.CubeMesh
	wirecube   *mesh.WireCubeMesh
	axis       *mesh.AxisMesh
	projection mgl32.Mat4
	scale      float32
	size       float32
}

var vertices = GenerateCubeVertices(1.0)

// Public constructor
func NewRenderer(scale float32, aspect float32) (*Renderer, error) {
	if err := gl.Init(); err != nil {
		return nil, err
	}
	gl.Enable(gl.DEPTH_TEST)

	shader, err := shader.NewShaderProgram()
	if err != nil {
		return nil, err
	}
	cube := mesh.NewCubeMesh(vertices)
	wirecube := mesh.NewWireCubeMesh(wireCubeVertices)
	axis := mesh.NewAxisMesh(axisVertices)
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), aspect, 0.1, 100.0)
	r := &Renderer{
		scale:      scale, // adjust based on grid size
		shader:     *shader,
		cube:       cube,
		wirecube:   wirecube,
		projection: projection,
		axis:       axis,
	}

	return r, nil
}

func (r *Renderer) Resize(w, h int32) {
	gl.Viewport(0, 0, w, h)
	aspect := float32(w) / float32(h)
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), aspect, 0.1, 100.0)
	r.projection = projection
}

func (r *Renderer) Draw(grid voxel.VoxelGrid, view mgl32.Mat4) {
	gl.ClearColor(0.9, 0.9, 0.9, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(r.shader.ID)

	// Set view and projection once per frame
	viewLoc := gl.GetUniformLocation(r.shader.ID, gl.Str("view\x00"))
	gl.UniformMatrix4fv(viewLoc, 1, false, &view[0])

	projLoc := gl.GetUniformLocation(r.shader.ID, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projLoc, 1, false, &r.projection[0])

	scale := r.scale
	center := float32(grid.Size) / 2
	bounds := float32(grid.Size) * scale
	model := mgl32.Scale3D(bounds/2, bounds/2, bounds/2) // scale from unit cube to voxel bounds
	lightDir := mgl32.Vec3{0.1, -0.1, 0.7}.Normalize()

	r.wirecube.Draw(r.shader, mgl32.Vec3{1, 1, 1}, model) // white cube
	arrowLen := bounds * 1.2
	cornerOffset := mgl32.Vec3{
		-bounds / 2, // shift from center to corner
		-bounds / 2,
		-bounds / 2,
	}
	cornerModel := mgl32.Translate3D(cornerOffset.X(), cornerOffset.Y(), cornerOffset.Z()).
		Mul4(mgl32.Scale3D(arrowLen, arrowLen, arrowLen))
	r.axis.Draw(r.shader, cornerModel)

	for z := range grid.Data {
		for y := range grid.Data[z] {
			for x := range grid.Data[z][y] {
				c := grid.At(x, y, z)
				if !c.Visible() {
					continue
				}

				color := mgl32.Vec3{
					float32(c.R()) / 255,
					float32(c.G()) / 255,
					float32(c.B()) / 255,
				}

				for _, dir := range dirs {
					nx, ny, nz := x+dir.dx, y+dir.dy, z+dir.dz
					if !grid.InBounds(nx, ny, nz) {
						continue
					}
					if !grid.At(nx, ny, nz).Visible() {
						continue
					}
					pos := mgl32.Vec3{
						(float32(x) - center) * scale,
						(float32(y) - center) * scale,
						(float32(z) - center) * scale,
					}
					model := mgl32.Translate3D(pos.Y(), pos.Z(), pos.X()).Mul4(
						mgl32.Scale3D(scale, scale, scale),
					)

					r.cube.DrawAt(r.shader, color, lightDir, model)
				}

			}
		}
	}
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
