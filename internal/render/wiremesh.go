package render

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type WireCubeMesh struct {
	vao uint32
}

func NewWireCubeMesh() *WireCubeMesh {
	var vao, vbo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)

	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(wireCubeVertices)*4, gl.Ptr(wireCubeVertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	return &WireCubeMesh{vao: vao}
}

func (c *WireCubeMesh) Draw(shader ShaderProgram, color mgl32.Vec3, model mgl32.Mat4) {
	gl.UseProgram(shader.ID)

	gl.BindVertexArray(c.vao)
	gl.Uniform3f(gl.GetUniformLocation(shader.ID, gl.Str("color\x00")), color.X(), color.Y(), color.Z())
	gl.UniformMatrix4fv(gl.GetUniformLocation(shader.ID, gl.Str("model\x00")), 1, false, &model[0])

	gl.DrawArrays(gl.LINES, 0, 24)
}
