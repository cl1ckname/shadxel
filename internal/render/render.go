package render

import (
	"shadxel/internal/render/mesh"
	"shadxel/internal/render/shader"
	"shadxel/internal/voxel"
	"unsafe"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/go-gl/gl/v3.3-core/gl"
)

var lightDir = mgl32.Vec3{0.1, -0.1, 0.7}.Normalize()

type Renderer struct {
	shader     shader.Program
	wirecube   *mesh.WireCubeMesh
	axis       *mesh.AxisMesh
	projection mgl32.Mat4
	scale      float32
	size       float32

	instanceVAO uint32
	instanceVBO uint32
	colorVBO    uint32
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
	wirecube := mesh.NewWireCubeMesh(wireCubeVertices)
	axis := mesh.NewAxisMesh(axisVertices)
	projection := mgl32.Perspective(mgl32.DegToRad(45.0), aspect, 0.1, 100.0)
	r := &Renderer{
		scale:      scale, // adjust based on grid size
		shader:     *shader,
		wirecube:   wirecube,
		projection: projection,
		axis:       axis,
	}

	gl.GenVertexArrays(1, &r.instanceVAO)
	gl.BindVertexArray(r.instanceVAO)

	// Upload static cube geometry
	var cubeVBO uint32
	gl.GenBuffers(1, &cubeVBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, cubeVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// Vertex positions (location = 0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// Normals (location = 1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	// Instance model matrix (locations 3,4,5,6)
	gl.GenBuffers(1, &r.instanceVBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, r.instanceVBO)
	for i := 0; i < 4; i++ {
		loc := uint32(3 + i)
		gl.EnableVertexAttribArray(loc)
		gl.VertexAttribPointer(loc, 4, gl.FLOAT, false, int32(unsafe.Sizeof(mgl32.Mat4{})), gl.PtrOffset(i*16))
		gl.VertexAttribDivisor(loc, 1)
	}

	// Instance color (location = 7)
	gl.GenBuffers(1, &r.colorVBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, r.colorVBO)
	gl.EnableVertexAttribArray(7)
	gl.VertexAttribPointer(7, 3, gl.FLOAT, false, 0, nil)
	gl.VertexAttribDivisor(7, 1)

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
	gl.Uniform3f(gl.GetUniformLocation(r.shader.ID, gl.Str("worldLightDir\x00")), lightDir.X(), lightDir.Y(), lightDir.Z())

	scale := r.scale
	center := float32(grid.Size) / 2
	bounds := float32(grid.Size) * scale
	model := mgl32.Scale3D(bounds/2, bounds/2, bounds/2) // scale from unit cube to voxel bounds
	transforms := make([]mgl32.Mat4, 0)
	colors := make([]mgl32.Vec3, 0)

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
				color := c.Color()

				if !grid.Hit(x, y, z) {
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
				transforms = append(transforms, model)
				colors = append(colors, color)
			}
		}
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, r.instanceVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(transforms)*int(unsafe.Sizeof(mgl32.Mat4{})), gl.Ptr(transforms), gl.DYNAMIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, r.colorVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(colors)*int(unsafe.Sizeof(mgl32.Vec3{})), gl.Ptr(colors), gl.DYNAMIC_DRAW)

	// Bind VAO and draw instances
	gl.BindVertexArray(r.instanceVAO)
	gl.DrawArraysInstanced(gl.TRIANGLES, 0, 36, int32(len(transforms)))
}
