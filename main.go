package main

import (
	"fmt"

	"bunny/ui"
)

func main() {
	// Create a window.
	window := ui.CreateWindow(800, 600, "Bunny UI Library")

	// Create UI components.
	label := ui.NewLabel(-0.8, 0.8, 1.6, 0.1, "Bunny UI Demo", [4]float32{1.0, 1.0, 1.0, 1.0})
	input := ui.NewInput(-0.8, 0.4, 1.6, 0.1, [4]float32{0.9, 0.9, 0.9, 1.0})
	radio := ui.NewRadioButton(-0.1, 0.1, 0.08, 0.08, [4]float32{0.8, 0.8, 0.8, 1.0})
	radio.OnSelectFunc = func() {
		fmt.Println("Radio button selected!")
	}
	button := ui.NewButton(-0.3, -0.5, 0.6, 0.2, [4]float32{0.4, 0.6, 1.0, 1.0}, "Click Me")
	button.OnClickFunc = func() {
		fmt.Println("Button clicked!")
	}

	slider := ui.NewSlider(-0.8, 0.6, 1.6, 0.1, [4]float32{0.5, 0.5, 0.5, 1.0}, [4]float32{0.8, 0.8, 0.8, 1.0})
	slider.OnValueChange = func(v float32) {
		fmt.Printf("Slider value: %.2f\n", v)
	}

	// Create a panel containing the components.
	root := ui.NewPanel(-0.9, -0.9, 1.8, 1.8, [4]float32{0.2, 0.2, 0.2, 1.0}, label, input, radio, button, slider)

	// Set up event callbacks.
	window.SetMouseButtonCallback(func(btn ui.MouseButton, action ui.Action, mods ui.ModifierKey) {
		if btn == ui.MouseButtonLeft {
			x, y := window.GetCursorPos()
			root.HandleMouse(x, y, action, 800, 600)
		}
	})

	window.SetCursorPosCallback(func(xpos, ypos float64) {
		root.HandleCursorPos(xpos, ypos, 800, 600)
	})

	window.SetKeyCallback(func(key ui.Key, action ui.Action, mods ui.ModifierKey) {
		root.HandleKey(key, action, mods)
	})

	// Run the UI loop.
	ui.Run(window, root)
}
