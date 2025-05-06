package shader

import (
	"fmt"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Program struct {
	ID uint32
}

func NewShaderProgram() (*Program, error) {
	program, err := makeProgram(vertexShaderSource, fragmentShaderSource)
	if err != nil {
		return nil, err
	}
	return &Program{ID: program}, nil
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

const vertexShaderSource = `
#version 330 core

layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aNormal;

// mat4 takes up locations 3, 4, 5, 6
layout (location = 3) in mat4 instanceModel;

// vec3 color at location 7
layout (location = 7) in vec3 instanceColor;

uniform mat4 view;
uniform mat4 projection;

out vec3 FragPos;
out vec3 Normal;
out vec3 Color;

void main() {
    mat4 model = instanceModel;

    FragPos = vec3(model * vec4(aPos, 1.0));
    Normal = mat3(transpose(inverse(model * view))) * aNormal;
    Color = instanceColor;

    gl_Position = projection * view * vec4(FragPos, 1.0);
}
` + "\x00"

const fragmentShaderSource = `
#version 330 core

in vec3 FragPos;
in vec3 Normal;
in vec3 Color;

out vec4 FragColor;

uniform vec3 worldLightDir;

void main() {
    float diff = max(dot(normalize(Normal), normalize(worldLightDir)), 0.0);
    vec3 ambient = 0.2 * Color;
    vec3 shadedColor = ambient + Color * diff;

    FragColor = vec4(shadedColor, 1.0);
	// FragColor = vec4(normalize(Normal) * 0.5 + 0.5, 1.0);
	// FragColor = vec4(normalize(worldLightDir) * 0.5 + 0.5, 1.0);
	
}
` + "\x00"
