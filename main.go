package main

import (
	"fmt"

	"bunny/ui"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func main() {
	// Create a window.
	window := ui.CreateWindow(800, 600, "Bunny UI Library")

	// Create UI components.
	button := ui.NewButton(-0.5, -0.5, 1.0, 1.0, [4]float32{1.0, 1.0, 1.0, 1.0})
	button.OnClickFunc = func() {
		fmt.Println("Button clicked!")
	}

	slider := ui.NewSlider(-0.8, 0.6, 1.6, 0.1, [4]float32{0.5, 0.5, 0.5, 1.0}, [4]float32{0.8, 0.8, 0.8, 1.0})
	slider.OnValueChange = func(v float32) {
		fmt.Printf("Slider value: %.2f\n", v)
	}

	// Create a panel containing the components.
	root := ui.NewPanel(-0.9, -0.9, 1.8, 1.8, [4]float32{0.2, 0.2, 0.2, 1.0}, button, slider)

	// Set up event callbacks.
	window.SetMouseButtonCallback(func(w *glfw.Window, btn glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		if btn == glfw.MouseButtonLeft {
			x, y := w.GetCursorPos()
			root.HandleMouse(x, y, action, 800, 600)
		}
	})

	window.SetCursorPosCallback(func(w *glfw.Window, xpos, ypos float64) {
		root.HandleCursorPos(xpos, ypos, 800, 600)
	})

	// Run the UI loop.
	ui.Run(window, root)
}
