package ui

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Panel struct {
	X, Y, Width, Height float32
	Color               [4]float32
	Children            []Component
	VAO                 uint32
}

func NewPanel(x, y, w, h float32, color [4]float32, children ...Component) *Panel {
	var VAO, VBO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)
	gl.BindVertexArray(VAO)
	vertices := []float32{
		x, y,
		x + w, y,
		x + w, y + h,
		x, y + h,
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, nil)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	return &Panel{
		X: x, Y: y, Width: w, Height: h,
		Color:    color,
		Children: children,
		VAO:      VAO,
	}
}

func (p *Panel) Draw(program uint32) {
	gl.UseProgram(program)
	colorLoc := gl.GetUniformLocation(program, gl.Str("color\x00"))
	gl.Uniform4f(colorLoc, p.Color[0], p.Color[1], p.Color[2], p.Color[3])
	gl.BindVertexArray(p.VAO)
	gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)
	gl.BindVertexArray(0)
	for _, child := range p.Children {
		child.Draw(program)
	}
}

func (p *Panel) HandleMouse(x, y float64, action glfw.Action, width, height int) {
	for _, child := range p.Children {
		cx, cy, cw, ch := child.GetBounds()
		xNorm := float32((x/float64(width))*2 - 1)
		yNorm := float32(1 - (y/float64(height))*2)
		if xNorm >= cx && xNorm <= cx+cw && yNorm >= cy && yNorm <= cy+ch {
			child.HandleMouse(x, y, action, width, height)
			break
		}
	}
}

func (p *Panel) HandleCursorPos(x, y float64, width, height int) {
	for _, child := range p.Children {
		child.HandleCursorPos(x, y, width, height)
	}
}

func (p *Panel) GetBounds() (x, y, w, h float32) {
	return p.X, p.Y, p.Width, p.Height
}

func (p *Panel) SetPosition(x, y float32) {
	p.X, p.Y = x, y
	// Recreate VAO
	gl.DeleteVertexArrays(1, &p.VAO)
	var VAO, VBO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)
	gl.BindVertexArray(VAO)
	vertices := []float32{
		x, y,
		x + p.Width, y,
		x + p.Width, y + p.Height,
		x, y + p.Height,
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, nil)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	p.VAO = VAO
}
