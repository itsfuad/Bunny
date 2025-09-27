package main

import "github.com/go-gl/glfw/v3.3/glfw"

type Component interface {
	Draw(program uint32)
	HandleMouse(x, y float64, action glfw.Action, width, height int)
	HandleCursorPos(x, y float64, width, height int)
	GetBounds() (x, y, w, h float32)
}