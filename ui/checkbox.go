package ui

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

type Checkbox struct {
	X, Y, Width, Height float32
	Color               [4]float32
	BackgroundColor     [4]float32
	IsChecked           bool
	OnCheckFunc         func()
	VAO                 uint32
}

func NewCheckbox(x, y, w, h float32, color [4]float32) *Checkbox {
	return &Checkbox{
		X: x, Y: y, Width: w, Height: h,
		Color:           [4]float32{0.2, 0.6, 1.0, 1.0}, // Default checked color
		BackgroundColor: color,
	}
}

func (cb *Checkbox) Draw(program uint32) {
	gl.UseProgram(program)
	colorUniform := gl.GetUniformLocation(program, gl.Str("color\x00"))

	// Draw outer square border
	vertices := []float32{
		cb.X, cb.Y,
		cb.X + cb.Width, cb.Y,
		cb.X + cb.Width, cb.Y + cb.Height,
		cb.X, cb.Y + cb.Height,
	}

	var VAO, VBO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)
	gl.BindVertexArray(VAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, nil)
	gl.EnableVertexAttribArray(0)

	// Draw border
	gl.Uniform4fv(colorUniform, 1, &cb.BackgroundColor[0])
	gl.DrawArrays(gl.LINE_LOOP, 0, 4)

	// If checked, draw filled inner square
	if cb.IsChecked {
		innerMargin := cb.Width * 0.1
		innerVertices := []float32{
			cb.X + innerMargin, cb.Y + innerMargin,
			cb.X + cb.Width - innerMargin, cb.Y + innerMargin,
			cb.X + cb.Width - innerMargin, cb.Y + cb.Height - innerMargin,
			cb.X + innerMargin, cb.Y + cb.Height - innerMargin,
		}

		var innerVAO, innerVBO uint32
		gl.GenVertexArrays(1, &innerVAO)
		gl.GenBuffers(1, &innerVBO)
		gl.BindVertexArray(innerVAO)
		gl.BindBuffer(gl.ARRAY_BUFFER, innerVBO)
		gl.BufferData(gl.ARRAY_BUFFER, len(innerVertices)*4, gl.Ptr(innerVertices), gl.STATIC_DRAW)
		gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, nil)
		gl.EnableVertexAttribArray(0)

		gl.Uniform4fv(colorUniform, 1, &cb.Color[0])
		gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)

		gl.BindVertexArray(0)
		gl.DeleteVertexArrays(1, &innerVAO)
		gl.DeleteBuffers(1, &innerVBO)
	}

	gl.BindVertexArray(0)
	gl.DeleteVertexArrays(1, &VAO)
	gl.DeleteBuffers(1, &VBO)
}

func (cb *Checkbox) HandleMouse(x, y float64, action Action, width, height int) {
	xNorm := float32((x/float64(width))*2 - 1)
	yNorm := float32(1 - (y/float64(height))*2)
	if action == Press && xNorm >= cb.X && xNorm <= cb.X+cb.Width && yNorm >= cb.Y && yNorm <= cb.Y+cb.Height {
		cb.IsChecked = !cb.IsChecked
		if cb.OnCheckFunc != nil {
			cb.OnCheckFunc()
		}
	}
}

func (cb *Checkbox) HandleCursorPos(x, y float64, width, height int) {
	// Checkboxes don't handle cursor
}

func (cb *Checkbox) GetBounds() (x, y, w, h float32) {
	return cb.X, cb.Y, cb.Width, cb.Height
}

func (cb *Checkbox) SetPosition(x, y float32) {
	cb.X, cb.Y = x, y
}

func (cb *Checkbox) StopDragging() {
	// Checkboxes don't drag
}

func (cb *Checkbox) HandleKey(key Key, action Action, mods ModifierKey) {
	// Checkboxes don't handle keys
}
