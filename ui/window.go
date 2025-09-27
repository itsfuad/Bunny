package ui

import (
	"fmt"
	"runtime"

	"bunny/render"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

// Input types to abstract GLFW
type MouseButton int

const (
	MouseButtonLeft   MouseButton = 0
	MouseButtonRight  MouseButton = 1
	MouseButtonMiddle MouseButton = 2
	MouseButton4      MouseButton = 3
	MouseButton5      MouseButton = 4
	MouseButton6      MouseButton = 5
	MouseButton7      MouseButton = 6
	MouseButton8      MouseButton = 7
)

type Action int

const (
	Release Action = 0
	Press   Action = 1
	Repeat  Action = 2
)

type ModifierKey int

const (
	ModShift    ModifierKey = 0x0001
	ModControl  ModifierKey = 0x0002
	ModAlt      ModifierKey = 0x0004
	ModSuper    ModifierKey = 0x0008
	ModCapsLock ModifierKey = 0x0010
	ModNumLock  ModifierKey = 0x0020
)

// BunnyWindow wraps the GLFW window and provides abstracted callbacks
type BunnyWindow struct {
	window              *glfw.Window
	mouseButtonCallback func(btn MouseButton, action Action, mods ModifierKey)
	cursorPosCallback   func(xpos, ypos float64)
	root                Component
}

func (w *BunnyWindow) SetMouseButtonCallback(cb func(btn MouseButton, action Action, mods ModifierKey)) {
	w.mouseButtonCallback = cb
	w.window.SetMouseButtonCallback(func(gw *glfw.Window, btn glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		if w.mouseButtonCallback != nil {
			w.mouseButtonCallback(MouseButton(btn), Action(action), ModifierKey(mods))
		}
	})
}

func (w *BunnyWindow) SetCursorPosCallback(cb func(xpos, ypos float64)) {
	w.cursorPosCallback = cb
	w.window.SetCursorPosCallback(func(gw *glfw.Window, xpos, ypos float64) {
		if w.cursorPosCallback != nil {
			w.cursorPosCallback(xpos, ypos)
		}
	})
}

func (w *BunnyWindow) GetCursorPos() (float64, float64) {
	return w.window.GetCursorPos()
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
uniform vec4 color;
void main() {
    FragColor = color;
}
`

func CreateWindow(width, height int, title string) *BunnyWindow {
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

	return &BunnyWindow{window: window}
}

func Run(window *BunnyWindow, root Component) {
	window.root = root
	// Set cursor enter callback to stop dragging when mouse enters the window
	window.window.SetCursorEnterCallback(func(gw *glfw.Window, entered bool) {
		fmt.Printf("Cursor enter: entered=%v\n", entered)
		if entered && window.root != nil {
			window.root.StopDragging()
		}
	})

	// Create shader program.
	program := render.CreateProgram(vertexShaderSource, fragmentShaderSource)

	// Main loop.
	for !window.window.ShouldClose() {
		// Clear the color buffer.
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// Draw the root component.
		root.Draw(program)

		// Swap buffers and poll events.
		window.window.SwapBuffers()
		glfw.PollEvents()
	}

	// Terminate GLFW when done.
	glfw.Terminate()
}
