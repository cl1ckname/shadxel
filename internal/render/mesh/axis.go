package mesh

import (
	"shadxel/internal/render/shader"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type AxisMesh struct {
	vao uint32
}

func NewAxisMesh(vert []float32) *AxisMesh {
	var vao, vbo uint32
	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)

	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vert)*4, gl.Ptr(vert), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	return &AxisMesh{vao: vao}
}

func (a *AxisMesh) Draw(shader shader.Program, model mgl32.Mat4) {
	gl.UseProgram(shader.ID)
	gl.BindVertexArray(a.vao)

	gl.UniformMatrix4fv(gl.GetUniformLocation(shader.ID, gl.Str("model\x00")), 1, false, &model[0])

	// X - red
	gl.Uniform3f(gl.GetUniformLocation(shader.ID, gl.Str("color\x00")), 1, 0, 0)
	gl.DrawArrays(gl.LINES, 0, 2)

	// Y - green
	gl.Uniform3f(gl.GetUniformLocation(shader.ID, gl.Str("color\x00")), 0, 1, 0)
	gl.DrawArrays(gl.LINES, 2, 2)

	// Z - blue
	gl.Uniform3f(gl.GetUniformLocation(shader.ID, gl.Str("color\x00")), 0, 0, 1)
	gl.DrawArrays(gl.LINES, 4, 2)
}
