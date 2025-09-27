package main

import (
	"log"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
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
		log.Fatal("Shader compilation failed:", shaderLog)
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
		log.Fatal("Program linking failed:", programLog)
	}
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)
	return program
}

func init() {
	// GLFW requires that the main thread is the one that initializes GLFW and creates windows.
	// This is necessary for cross-platform compatibility, especially on macOS.
	runtime.LockOSThread()
}

const vertexShaderSource = `
#version 330 core
layout (location = 0) in vec2 aPos;
void main() {
    gl_Position = vec4(aPos, 0.0, 1.0);
}
`

const fragmentShaderSource = `
#version 330 core
out vec4 FragColor;
void main() {
    FragColor = vec4(1.0, 1.0, 1.0, 1.0);
}
`

func main() {
	// Initialize GLFW library. This must be done before any other GLFW functions.
	if err := glfw.Init(); err != nil {
		log.Fatal("Failed to initialize GLFW:", err)
	}
	defer glfw.Terminate() // Ensure GLFW is terminated when the program exits.

	// Set GLFW options for the OpenGL context.
	// We use OpenGL 3.3 core profile for modern OpenGL features.
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// Create a window with the specified dimensions and title.
	window, err := glfw.CreateWindow(800, 600, "Bunny UI Library", nil, nil)
	if err != nil {
		log.Fatal("Failed to create GLFW window:", err)
	}
	defer window.Destroy() // Ensure the window is destroyed when done.

	// Make the OpenGL context of the window current on the calling thread.
	window.MakeContextCurrent()

	// Initialize OpenGL. This loads all OpenGL function pointers.
	if err := gl.Init(); err != nil {
		log.Fatal("Failed to initialize OpenGL:", err)
	}

	// Set the viewport to match the window size.
	gl.Viewport(0, 0, 800, 600)

	// Create shader program.
	program := createProgram(vertexShaderSource, fragmentShaderSource)

	// Set up vertex data and buffers.
	var VAO, VBO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)
	gl.BindVertexArray(VAO)
	vertices := []float32{
		-0.5, -0.5, // bottom left
		 0.5, -0.5, // bottom right
		 0.5,  0.5, // top right
		-0.5,  0.5, // top left
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, nil)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	// Set the swap interval for the current OpenGL context.
	// 1 means V-Sync is enabled, which synchronizes the frame rate with the monitor's refresh rate.
	glfw.SwapInterval(1)

	// Main loop: runs until the window should be closed.
	for !window.ShouldClose() {
		// Clear the color buffer with a default color (black).
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// Draw the rectangle.
		gl.UseProgram(program)
		gl.BindVertexArray(VAO)
		gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)
		gl.BindVertexArray(0)

		// Swap the front and back buffers of the window.
		// This displays the rendered frame and prepares for the next one.
		window.SwapBuffers()

		// Poll for and process events (e.g., window close, keyboard input).
		glfw.PollEvents()
	}
}