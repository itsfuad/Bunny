package ui

import (
	"math"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type RadioButton struct {
	X, Y, Width, Height float32
	Color               [4]float32
	BackgroundColor     [4]float32
	IsSelected          bool
	OnSelectFunc        func()
	VAO                 uint32
}

func NewRadioButton(x, y, w, h float32, color [4]float32) *RadioButton {
	return &RadioButton{
		X: x, Y: y, Width: w, Height: h,
		Color:           [4]float32{0.2, 0.6, 1.0, 1.0}, // Default selected color
		BackgroundColor: color,
	}
}

func (rb *RadioButton) Draw(program uint32) {
	gl.UseProgram(program)
	colorUniform := gl.GetUniformLocation(program, gl.Str("color\x00"))

	centerX := rb.X + rb.Width/2
	centerY := rb.Y + rb.Height/2
	radius := rb.Width / 2
	segments := 32

	// Draw outer circle border (unfilled)
	borderVertices := make([]float32, 0, segments*2)
	for i := 0; i <= segments; i++ {
		angle := float32(i) * 2 * 3.14159 / float32(segments)
		vx := centerX + radius*float32(math.Cos(float64(angle)))
		vy := centerY + radius*float32(math.Sin(float64(angle)))
		borderVertices = append(borderVertices, vx, vy)
	}

	var borderVAO, borderVBO uint32
	gl.GenVertexArrays(1, &borderVAO)
	gl.GenBuffers(1, &borderVBO)
	gl.BindVertexArray(borderVAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, borderVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(borderVertices)*4, gl.Ptr(borderVertices), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, nil)
	gl.EnableVertexAttribArray(0)

	// Draw border with background color
	gl.Uniform4fv(colorUniform, 1, &rb.BackgroundColor[0])
	gl.DrawArrays(gl.LINE_LOOP, 0, int32(len(borderVertices)/2))

	// If selected, draw filled inner circle
	if rb.IsSelected {
		innerRadius := radius * 0.6 // Make inner circle 60% of outer radius
		innerVertices := make([]float32, 0, segments*2+2)
		innerVertices = append(innerVertices, centerX, centerY) // Center point

		for i := 0; i <= segments; i++ {
			angle := float32(i) * 2 * 3.14159 / float32(segments)
			vx := centerX + innerRadius*float32(math.Cos(float64(angle)))
			vy := centerY + innerRadius*float32(math.Sin(float64(angle)))
			innerVertices = append(innerVertices, vx, vy)
		}

		var innerVAO, innerVBO uint32
		gl.GenVertexArrays(1, &innerVAO)
		gl.GenBuffers(1, &innerVBO)
		gl.BindVertexArray(innerVAO)
		gl.BindBuffer(gl.ARRAY_BUFFER, innerVBO)
		gl.BufferData(gl.ARRAY_BUFFER, len(innerVertices)*4, gl.Ptr(innerVertices), gl.STATIC_DRAW)
		gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, nil)
		gl.EnableVertexAttribArray(0)

		selectedColor := rb.Color // Use component's color for selected inner
		gl.Uniform4fv(colorUniform, 1, &selectedColor[0])
		gl.DrawArrays(gl.TRIANGLE_FAN, 0, int32(len(innerVertices)/2))

		gl.BindVertexArray(0)
		gl.DeleteVertexArrays(1, &innerVAO)
		gl.DeleteBuffers(1, &innerVBO)
	}

	gl.BindVertexArray(0)
	gl.DeleteVertexArrays(1, &borderVAO)
	gl.DeleteBuffers(1, &borderVBO)
}

func (rb *RadioButton) HandleMouse(x, y float64, action Action, width, height int) {
	if rb.IsSelected {
		return // Already selected, do nothing
	}
	xNorm := float32((x/float64(width))*2 - 1)
	yNorm := float32(1 - (y/float64(height))*2)
	if action == Press && xNorm >= rb.X && xNorm <= rb.X+rb.Width && yNorm >= rb.Y && yNorm <= rb.Y+rb.Height {
		rb.IsSelected = true
		if rb.OnSelectFunc != nil {
			rb.OnSelectFunc()
		}
	}
}

func (rb *RadioButton) HandleCursorPos(x, y float64, width, height int) {
	// Radio buttons don't handle cursor
}

func (rb *RadioButton) GetBounds() (x, y, w, h float32) {
	return rb.X, rb.Y, rb.Width, rb.Height
}

func (rb *RadioButton) SetPosition(x, y float32) {
	rb.X, rb.Y = x, y
}

func (rb *RadioButton) StopDragging() {
	// Radio buttons don't drag
}

func (rb *RadioButton) HandleKey(key Key, action Action, mods ModifierKey) {
	// Radio buttons don't handle keys
}
