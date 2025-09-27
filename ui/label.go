package ui

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

type Label struct {
	X, Y, Width, Height float32
	Text                string
	Color               [4]float32
	BackgroundColor     [4]float32
}

func NewLabel(x, y, w, h float32, text string, color [4]float32) *Label {
	return &Label{X: x, Y: y, Width: w, Height: h, Text: text, Color: [4]float32{0.1, 0.1, 0.1, 1.0}, BackgroundColor: color}
}

func (l *Label) Draw(program uint32) {
	// For now, draw a colored rectangle as placeholder
	vertices := []float32{
		l.X, l.Y,
		l.X + l.Width, l.Y,
		l.X + l.Width, l.Y + l.Height,
		l.X, l.Y + l.Height,
	}
	var VAO, VBO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)
	gl.BindVertexArray(VAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(0)
	gl.UseProgram(program)
	colorUniform := gl.GetUniformLocation(program, gl.Str("color\x00"))
	gl.Uniform4fv(colorUniform, 1, &l.BackgroundColor[0])
	gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)
	gl.BindVertexArray(0)
	gl.DeleteVertexArrays(1, &VAO)
	gl.DeleteBuffers(1, &VBO)
}

func (l *Label) HandleMouse(x, y float64, action Action, width, height int) {
	// Labels don't handle mouse
}

func (l *Label) HandleCursorPos(x, y float64, width, height int) {
	// Labels don't handle cursor
}

func (l *Label) GetBounds() (x, y, w, h float32) {
	return l.X, l.Y, l.Width, l.Height
}

func (l *Label) SetPosition(x, y float32) {
	l.X, l.Y = x, y
}

func (l *Label) StopDragging() {
	// Labels don't drag
}

func (l *Label) HandleKey(key Key, action Action, mods ModifierKey) {
	// Labels don't handle keys
}
