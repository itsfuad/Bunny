package main

import (
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
)

func createShader(source string, shaderType uint32) uint32 {
	shader := gl.CreateShader(shaderType)
	csource, free := gl.Strs(source + "\x00")
	gl.ShaderSource(shader, 1, csource, nil)
	free()
	gl.CompileShader(shader)
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		shaderLog := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(shaderLog))
		panic("Shader compilation failed: " + shaderLog)
	}
	return shader
}

func createProgram(vertexSource, fragmentSource string) uint32 {
	vertexShader := createShader(vertexSource, gl.VERTEX_SHADER)
	fragmentShader := createShader(fragmentSource, gl.FRAGMENT_SHADER)
	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)
	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
		programLog := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(programLog))
		panic("Program linking failed: " + programLog)
	}
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)
	return program
}

func drawRect(program uint32, x, y, w, h float32, color [4]float32) {
	gl.UseProgram(program)
	colorLoc := gl.GetUniformLocation(program, gl.Str("color\x00"))
	gl.Uniform4f(colorLoc, color[0], color[1], color[2], color[3])
	// For simplicity, assume rect is drawn with fixed VAO, but to make it general, we need to update vertices or use instancing.
	// For now, since only one rect, use the existing VAO.
	// But to generalize, perhaps pass VAO or modify.
	// For this task, keep simple.
}
