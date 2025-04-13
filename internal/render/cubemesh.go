package render

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type CubeMesh struct {
	vao uint32
}

func NewCubeMesh() *CubeMesh {
	var vao, vbo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)

	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4)) // normal
	gl.EnableVertexAttribArray(1)

	return &CubeMesh{vao: vao}
}

func (c *CubeMesh) DrawAt(shader ShaderProgram, color mgl32.Vec3, model mgl32.Mat4) {
	gl.UseProgram(shader.ID)

	gl.BindVertexArray(c.vao)
	gl.Uniform3f(gl.GetUniformLocation(shader.ID, gl.Str("color\x00")), color.X(), color.Y(), color.Z())
	gl.UniformMatrix4fv(gl.GetUniformLocation(shader.ID, gl.Str("model\x00")), 1, false, &model[0])
	gl.DrawArrays(gl.TRIANGLES, 0, 36)
}

var vertices = GenerateCubeVertices(1.0)
