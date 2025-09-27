package ui

import (
	"runtime"

	"bunny/render"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

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
uniform vec4 color;
void main() {
    FragColor = color;
}
`

func CreateWindow(width, height int, title string) *glfw.Window {
	// Initialize GLFW library.
	if err := glfw.Init(); err != nil {
		panic("Failed to initialize GLFW: " + err.Error())
	}

	// Set GLFW options for the OpenGL context.
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// Create a window.
	window, err := glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		panic("Failed to create GLFW window: " + err.Error())
	}

	// Make the OpenGL context current.
	window.MakeContextCurrent()

	// Initialize OpenGL.
	if err := gl.Init(); err != nil {
		panic("Failed to initialize OpenGL: " + err.Error())
	}

	// Set the viewport.
	gl.Viewport(0, 0, int32(width), int32(height))

	// Enable V-Sync.
	glfw.SwapInterval(1)

	return window
}

func Run(window *glfw.Window, root Component) {
	// Create shader program.
	program := render.CreateProgram(vertexShaderSource, fragmentShaderSource)

	// Main loop.
	for !window.ShouldClose() {
		// Clear the color buffer.
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// Draw the root component.
		root.Draw(program)

		// Swap buffers and poll events.
		window.SwapBuffers()
		glfw.PollEvents()
	}

	// Terminate GLFW when done.
	glfw.Terminate()
}
