package main

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Orientation int

const (
	VERTICAL Orientation = iota
	HORIZONTAL
)

type StackLayout struct {
	X, Y, Width, Height float32
	Orientation         Orientation
	Children            []Component
	Spacing             float32
}

func NewStackLayout(x, y, w, h float32, orientation Orientation, spacing float32, children ...Component) *StackLayout {
	layout := &StackLayout{
		X: x, Y: y, Width: w, Height: h,
		Orientation: orientation,
		Children:    children,
		Spacing:     spacing,
	}
	layout.arrangeChildren()
	return layout
}

func (sl *StackLayout) arrangeChildren() {
	switch sl.Orientation {
	case VERTICAL:
		currentY := sl.Y
		for _, child := range sl.Children {
			_, _, _, ch := child.GetBounds()
			child.SetPosition(sl.X, currentY)
			currentY += ch + sl.Spacing
		}
	case HORIZONTAL:
		currentX := sl.X
		for _, child := range sl.Children {
			_, _, cw, _ := child.GetBounds()
			child.SetPosition(currentX, sl.Y)
			currentX += cw + sl.Spacing
		}
	}
}

func (sl *StackLayout) Draw(program uint32) {
	for _, child := range sl.Children {
		child.Draw(program)
	}
}

func (sl *StackLayout) HandleMouse(x, y float64, action glfw.Action, width, height int) {
	for _, child := range sl.Children {
		cx, cy, cw, ch := child.GetBounds()
		xNorm := float32((x/float64(width))*2 - 1)
		yNorm := float32(1 - (y/float64(height))*2)
		if xNorm >= cx && xNorm <= cx+cw && yNorm >= cy && yNorm <= cy+ch {
			child.HandleMouse(x, y, action, width, height)
			break // or handle all
		}
	}
}

func (sl *StackLayout) HandleCursorPos(x, y float64, width, height int) {
	for _, child := range sl.Children {
		child.HandleCursorPos(x, y, width, height)
	}
}

func (sl *StackLayout) GetBounds() (x, y, w, h float32) {
	return sl.X, sl.Y, sl.Width, sl.Height
}

func (sl *StackLayout) SetPosition(x, y float32) {
	sl.X, sl.Y = x, y
	sl.arrangeChildren() // re-arrange children
}
