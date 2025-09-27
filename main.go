package main

import (
	"fmt"
	"log"
	"runtime"

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

	// Create a button.
	button := NewButton(-0.5, -0.5, 1.0, 1.0, [4]float32{1.0, 1.0, 1.0, 1.0})
	button.OnClickFunc = func() {
		fmt.Println("Button clicked!")
	}

	// Create a slider.
	slider := NewSlider(-0.8, 0.6, 1.6, 0.1, [4]float32{0.5, 0.5, 0.5, 1.0}, [4]float32{0.8, 0.8, 0.8, 1.0})
	slider.OnValueChange = func(v float32) {
		fmt.Printf("Slider value: %.2f\n", v)
	}

	// Set mouse callback.
	window.SetMouseButtonCallback(func(w *glfw.Window, btn glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		if btn == glfw.MouseButtonLeft {
			x, y := w.GetCursorPos()
			button.HandleMouse(x, y, action, 800, 600)
			slider.HandleMouse(x, y, action, 800, 600)
		}
	})

	// Set cursor position callback for dragging.
	window.SetCursorPosCallback(func(w *glfw.Window, xpos, ypos float64) {
		slider.HandleCursorPos(xpos, ypos, 800, 600)
	})

	// Set the swap interval for the current OpenGL context.
	// 1 means V-Sync is enabled, which synchronizes the frame rate with the monitor's refresh rate.
	glfw.SwapInterval(1)

	// Main loop: runs until the window should be closed.
	for !window.ShouldClose() {
		// Clear the color buffer with a default color (black).
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// Draw the button.
		button.Draw(program)

		// Draw the slider.
		slider.Draw(program)

		// Swap the front and back buffers of the window.
		// This displays the rendered frame and prepares for the next one.
		window.SwapBuffers()

		// Poll for and process events (e.g., window close, keyboard input).
		glfw.PollEvents()
	}
}
