package ui

type Component interface {
	Draw(program uint32)
	HandleMouse(x, y float64, action Action, width, height int)
	HandleCursorPos(x, y float64, width, height int)
	GetBounds() (x, y, w, h float32)
	SetPosition(x, y float32)
	StopDragging()
}
