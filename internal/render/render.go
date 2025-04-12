package render

import (
	"fmt"
	"shadxel/internal/voxel"

	"github.com/go-gl/gl/v3.3-core/gl"
)

const vertexShaderSource = `
#version 330 core
layout (location = 0) in vec2 aPos;

uniform float rotation;
uniform float scale;
uniform vec2 offset;

void main() {
    float s = sin(rotation);
    float c = cos(rotation);

    vec2 worldPos = (aPos * scale) + offset;
    vec2 rotated = vec2(
        c * worldPos.x - s * worldPos.y,
        s * worldPos.x + c * worldPos.y
    );

    gl_Position = vec4(rotated, 0.0, 1.0);
}
` + "\x00"

const fragmentShaderSource = `
#version 330 core
out vec4 FragColor;
uniform vec3 color;
void main() {
    FragColor = vec4(color, 1.0);
}
` + "\x00"

type Renderer struct {
	program  uint32
	vao      uint32
	rotation float32
	scale    float32
}

// Public constructor
func NewRenderer() (*Renderer, error) {
	r := &Renderer{
		scale: 1.5 / 50.0, // adjust based on grid size
	}

	if err := r.initGL(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Renderer) initGL() error {
	var err error
	r.program, err = makeProgram(vertexShaderSource, fragmentShaderSource)
	if err != nil {
		return err
	}

	// Single square (2 triangles)
	vertices := []float32{
		-0.5, -0.5,
		0.5, -0.5,
		0.5, 0.5,
		-0.5, -0.5,
		0.5, 0.5,
		-0.5, 0.5,
	}

	var vbo uint32
	gl.GenVertexArrays(1, &r.vao)
	gl.BindVertexArray(r.vao)

	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.BindVertexArray(0)
	return nil
}

func (r *Renderer) DrawGrid(grid voxel.Grid, rotation float32) {
	gl.ClearColor(0.1, 0.1, 0.1, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.UseProgram(r.program)
	gl.BindVertexArray(r.vao)

	locRot := gl.GetUniformLocation(r.program, gl.Str("rotation\x00"))
	gl.Uniform1f(locRot, rotation)

	scale := float32(1 / 50.0)
	locColor := gl.GetUniformLocation(r.program, gl.Str("color\x00"))
	locOffset := gl.GetUniformLocation(r.program, gl.Str("offset\x00"))
	locScale := gl.GetUniformLocation(r.program, gl.Str("scale\x00"))

	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			color := grid[y][x]
			gl.Uniform3f(locColor,
				float32(color.R)/255.0,
				float32(color.G)/255.0,
				float32(color.B)/255.0)

			cx := float32(x)*scale - 0.5 + scale/2
			cy := float32(y)*scale - 0.5 + scale/2
			gl.Uniform2f(locOffset, cx, cy)
			gl.Uniform1f(locScale, scale)

			gl.DrawArrays(gl.TRIANGLES, 0, 6)
		}
	}

	gl.BindVertexArray(0)
}

func makeProgram(vertexSrc, fragmentSrc string) (uint32, error) {
	vs, err := compileShader(vertexSrc, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}
	fs, err := compileShader(fragmentSrc, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vs)
	gl.AttachShader(program, fs)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
		logMsg := make([]byte, logLength)
		gl.GetProgramInfoLog(program, logLength, nil, &logMsg[0])
		return 0, fmt.Errorf("program link error: %s", string(logMsg))
	}

	gl.DeleteShader(vs)
	gl.DeleteShader(fs)

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		logMsg := make([]byte, logLength)
		gl.GetShaderInfoLog(shader, logLength, nil, &logMsg[0])
		return 0, fmt.Errorf("shader compile error: %s", string(logMsg))
	}
	return shader, nil
}
