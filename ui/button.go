package ui

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

type Button struct {
	X, Y, Width, Height float32
	Color               [4]float32
	BackgroundColor     [4]float32
	Text                string
	OnClickFunc         func()
	VAO                 uint32
}

func NewButton(x, y, w, h float32, color [4]float32, text string) *Button {
	var VAO, VBO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)
	gl.BindVertexArray(VAO)
	// Calculate vertices based on x,y,w,h in -1 to 1 coords
	vertices := []float32{
		x, y, // bottom left
		x + w, y, // bottom right
		x + w, y + h, // top right
		x, y + h, // top left
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, nil)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	return &Button{
		X: x, Y: y, Width: w, Height: h,
		Color:           [4]float32{0.1, 0.1, 0.1, 1.0}, // Default text color
		BackgroundColor: color,
		Text:            text,
		VAO:             VAO,
	}
}

func (b *Button) Draw(program uint32) {
	gl.UseProgram(program)
	colorLoc := gl.GetUniformLocation(program, gl.Str("color\x00"))
	gl.Uniform4f(colorLoc, b.BackgroundColor[0], b.BackgroundColor[1], b.BackgroundColor[2], b.BackgroundColor[3])
	gl.BindVertexArray(b.VAO)
	gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)

	// Draw border
	borderColor := [4]float32{0.3, 0.3, 0.3, 1.0}
	gl.Uniform4fv(colorLoc, 1, &borderColor[0])
	gl.DrawArrays(gl.LINE_LOOP, 0, 4)

	// Draw text if present
	if len(b.Text) > 0 {
		charWidth := b.Width / float32(len(b.Text)+1) // Leave some padding
		if charWidth > b.Height*0.6 {
			charWidth = b.Height * 0.4 // Limit character width for buttons
		}
		charHeight := b.Height * 0.4

		// Center the text
		textWidth := float32(len(b.Text)) * charWidth
		textX := b.X + (b.Width-textWidth)/2
		textY := b.Y + b.Height*0.3

		DrawBitmapText(program, b.Text, textX, textY, charWidth, charHeight, b.Color)
	}

	gl.BindVertexArray(0)
}

func (b *Button) HandleMouse(x, y float64, action Action, width, height int) {
	// Convert mouse coords to OpenGL coords
	xNorm := float32((x/float64(width))*2 - 1)
	yNorm := float32(1 - (y/float64(height))*2)
	if action == Press && xNorm >= b.X && xNorm <= b.X+b.Width && yNorm >= b.Y && yNorm <= b.Y+b.Height {
		b.OnClick()
	}
}

func (b *Button) HandleCursorPos(x, y float64, width, height int) {
	// Buttons don't handle cursor position
}

func (b *Button) StopDragging() {
	// Buttons don't have dragging state
}

func (b *Button) HandleKey(key Key, action Action, mods ModifierKey) {
	// Buttons don't handle keys
}

func (b *Button) GetBounds() (x, y, w, h float32) {
	return b.X, b.Y, b.Width, b.Height
}

func (b *Button) OnClick() {
	if b.OnClickFunc != nil {
		b.OnClickFunc()
	}
}

func (b *Button) SetPosition(x, y float32) {
	b.X, b.Y = x, y
	// Recreate VAO with new vertices
	gl.DeleteVertexArrays(1, &b.VAO)
	var VAO, VBO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)
	gl.BindVertexArray(VAO)
	vertices := []float32{
		x, y,
		x + b.Width, y,
		x + b.Width, y + b.Height,
		x, y + b.Height,
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, nil)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	b.VAO = VAO
}
