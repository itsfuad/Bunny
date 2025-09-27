# Bunny UI Library

A cross-platform UI library for Go, built with GLFW for windowing and OpenGL for rendering.

## Description

Bunny is a simple, modular UI library designed for Go applications. It provides basic UI components like buttons, sliders, labels, panels, and stack layouts, rendered using OpenGL primitives. The library abstracts GLFW details for easier use, handles proper event bubbling, and focuses on cross-platform compatibility.

## Features

- **Cross-platform**: Works on Windows, macOS, and Linux.
- **Modular**: Components are structs with methods for drawing, updating, and event handling.
- **OpenGL-based**: Uses modern OpenGL 3.3 core for efficient rendering.
- **GLFW abstraction**: Provides abstracted input types and window management without exposing GLFW directly.
- **Proper event handling**: Supports event bubbling with targeted press events and broadcast release events.
- **Components**:
  - Button: Clickable rectangle with callback.
  - Slider: Draggable value selector with callbacks.
  - Label: Text display component (placeholder rendering).
  - RadioButton: Selectable button with callback.
  - Input: Single-line text input with keyboard handling.
  - Panel: Background container for grouping components.
  - StackLayout: Automatic vertical/horizontal arrangement of child components.

## Installation

Ensure you have Go 1.21+ and the necessary system dependencies for GLFW (e.g., on Ubuntu: `sudo apt install libgl1-mesa-dev xorg-dev`).

```bash
go get github.com/yourusername/bunny
```

## Usage

Import the library and create a window, components, and run the UI loop.

### Basic Example

```go
package main

import (
    "fmt"
    "bunny/ui"
)

func main() {
    // Create a window
    window := ui.CreateWindow(800, 600, "My App")

    // Create components
    label := ui.NewLabel(-0.8, 0.8, 1.6, 0.1, "My App", [4]float32{0, 0, 0, 1})
    input := ui.NewInput(-0.8, 0.4, 1.6, 0.1, [4]float32{0.9, 0.9, 0.9, 1})
    radio := ui.NewRadioButton(-0.5, 0.3, 0.2, 0.2, [4]float32{0, 1, 0, 1})
    radio.OnSelectFunc = func() { fmt.Println("Selected!") }
    button := ui.NewButton(-0.5, -0.5, 1.0, 1.0, [4]float32{1, 1, 1, 1})
    button.OnClickFunc = func() { fmt.Println("Clicked!") }

    slider := ui.NewSlider(-0.8, 0.6, 1.6, 0.1, [4]float32{0.5, 0.5, 0.5, 1}, [4]float32{0.8, 0.8, 0.8, 1})
    slider.OnValueChange = func(v float32) { fmt.Printf("Value: %.2f\n", v) }

    // Panel to group components
    root := ui.NewPanel(-0.9, -0.9, 1.8, 1.8, [4]float32{0.2, 0.2, 0.2, 1}, label, input, radio, button, slider)

    // Event callbacks
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

    // Run the UI
    ui.Run(window, root)
}
```

### Components API

- **Button**: `NewButton(x, y, w, h, color) *Button`
  - `OnClickFunc func()`: Set click callback.
  - Methods: `Draw(program)`, `HandleMouse(...)`, etc.

- **Slider**: `NewSlider(x, y, w, h, trackColor, knobColor) *Slider`
  - `Value float32`: Current value (0-1).
  - `OnValueChange func(float32)`: Callback on value change.
  - Methods: `Draw(program)`, `HandleMouse(...)`, `HandleCursorPos(...)`, `Update(value)`.

- **Label**: `NewLabel(x, y, w, h, text, color) *Label`
  - Displays text (placeholder rendering).
  - Methods: `Draw(program)`, etc.

- **RadioButton**: `NewRadioButton(x, y, w, h, color) *RadioButton`
  - `IsSelected bool`: Selection state.
  - `OnSelectFunc func()`: Set selection callback.
  - Methods: `Draw(program)`, `HandleMouse(...)`, etc.

- **Input**: `NewInput(x, y, w, h, color) *Input`
  - `Text string`: Current input text.
  - `CursorPos int`: Cursor position in text.
  - `IsFocused bool`: Focus state.
  - Methods: `Draw(program)`, `HandleMouse(...)`, `HandleKey(...)`, etc.

- **Panel**: `NewPanel(x, y, w, h, color, children...) *Panel`
  - Container for other components.

- **StackLayout**: `NewStackLayout(x, y, w, h, orientation, spacing, children...) *StackLayout`
  - `orientation`: `VERTICAL` or `HORIZONTAL`.
  - Automatically positions children.

All components implement the `Component` interface with `Draw`, `HandleMouse`, `HandleCursorPos`, `HandleKey`, `GetBounds`, `SetPosition`, `StopDragging`.

### Window Management

- `CreateWindow(width, height, title) *BunnyWindow`: Initializes GLFW/OpenGL and creates a window with abstracted callbacks.
- `BunnyWindow.SetMouseButtonCallback(func(btn MouseButton, action Action, mods ModifierKey))`: Set mouse button event handler.
- `BunnyWindow.SetCursorPosCallback(func(xpos, ypos float64))`: Set cursor position event handler.
- `BunnyWindow.SetKeyCallback(func(key Key, action Action, mods ModifierKey))`: Set keyboard event handler.
- `BunnyWindow.GetCursorPos() (float64, float64)`: Get current cursor position.
- `Run(window *BunnyWindow, rootComponent Component)`: Runs the render loop until window closes.

## Event Handling

Bunny handles events with proper bubbling:
- **Press events**: Targeted to the component under the mouse.
- **Release events**: Broadcast to all components to ensure dragging stops regardless of release position.
- **Cursor events**: Broadcast to all components for dragging updates.

This ensures intuitive behavior, such as sliders stopping drag when the mouse button is released anywhere.

## Coordinates

Components use OpenGL normalized device coordinates (-1 to 1). X/Y are bottom-left, width/height are relative.

## Dependencies

- `github.com/go-gl/glfw/v3.3/glfw`
- `github.com/go-gl/gl/v3.3-core/gl`

## License

MIT License. See LICENSE file.

## Contributing

Contributions welcome! Open issues or PRs on GitHub.