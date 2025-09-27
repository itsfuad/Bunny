package ui

import (
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Input struct {
	X, Y, Width, Height float32
	Text                string
	Color               [4]float32
	IsFocused           bool
	CursorPos           int
	VAO                 uint32
}

func NewInput(x, y, w, h float32, color [4]float32) *Input {
	var VAO, VBO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.GenBuffers(1, &VBO)
	gl.BindVertexArray(VAO)
	// Calculate vertices for the input rectangle
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
	return &Input{
		X: x, Y: y, Width: w, Height: h,
		Color:     color,
		CursorPos: 0,
		VAO:       VAO,
	}
}

func (i *Input) Draw(program uint32) {
	gl.UseProgram(program)
	colorUniform := gl.GetUniformLocation(program, gl.Str("color\x00"))

	// Draw rectangle background
	if i.IsFocused {
		focusedColor := [4]float32{1.0, 1.0, 1.0, 1.0} // White when focused
		gl.Uniform4fv(colorUniform, 1, &focusedColor[0])
	} else {
		gl.Uniform4fv(colorUniform, 1, &i.Color[0])
	}
	gl.BindVertexArray(i.VAO)
	gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)

	// Draw border
	borderColor := [4]float32{0.3, 0.3, 0.3, 1.0}
	gl.Uniform4fv(colorUniform, 1, &borderColor[0])
	gl.DrawArrays(gl.LINE_LOOP, 0, 4)

	// Draw text using the bitmap font renderer
	if len(i.Text) > 0 {
		charWidth := (i.Width - 0.04) / float32(len(i.Text)+1) // Leave padding
		if charWidth > i.Height*0.7 {
			charWidth = i.Height * 0.7 // Limit character width
		}
		charHeight := i.Height * 0.7
		textX := i.X + 0.02
		textY := i.Y + i.Height*0.15
		textColor := [4]float32{0.1, 0.1, 0.1, 1.0} // Dark text

		DrawBitmapText(program, i.Text, textX, textY, charWidth, charHeight, textColor)
	}

	// Draw cursor if focused
	if i.IsFocused {
		cursorColor := [4]float32{0.0, 0.0, 0.0, 1.0} // Black cursor
		gl.Uniform4fv(colorUniform, 1, &cursorColor[0])

		// Use the same character width calculation as the text
		var charWidth float32
		if len(i.Text) > 0 {
			charWidth = (i.Width - 0.04) / float32(len(i.Text)+1)
			if charWidth > i.Height*0.7 {
				charWidth = i.Height * 0.7
			}
		} else {
			// When no text, use a reasonable default character width
			charWidth = i.Height * 0.5
		}
		cursorX := i.X + 0.02 + float32(i.CursorPos)*charWidth
		cursorY := i.Y + i.Height*0.15 // Match text Y position
		cursorHeight := i.Height * 0.7 // Match text height
		cursorWidth := float32(0.005)

		cursorVertices := []float32{
			cursorX, cursorY,
			cursorX + cursorWidth, cursorY,
			cursorX + cursorWidth, cursorY + cursorHeight,
			cursorX, cursorY + cursorHeight,
		}

		var cursorVAO, cursorVBO uint32
		gl.GenVertexArrays(1, &cursorVAO)
		gl.GenBuffers(1, &cursorVBO)
		gl.BindVertexArray(cursorVAO)
		gl.BindBuffer(gl.ARRAY_BUFFER, cursorVBO)
		gl.BufferData(gl.ARRAY_BUFFER, len(cursorVertices)*4, gl.Ptr(cursorVertices), gl.STATIC_DRAW)
		gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, nil)
		gl.EnableVertexAttribArray(0)

		gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)

		gl.BindVertexArray(0)
		gl.DeleteVertexArrays(1, &cursorVAO)
		gl.DeleteBuffers(1, &cursorVBO)
	}

	gl.BindVertexArray(0)
}

func (i *Input) HandleMouse(x, y float64, action Action, width, height int) {
	xNorm := float32((x/float64(width))*2 - 1)
	yNorm := float32(1 - (y/float64(height))*2)
	if action == Press {
		if xNorm >= i.X && xNorm <= i.X+i.Width && yNorm >= i.Y && yNorm <= i.Y+i.Height {
			i.IsFocused = true
			// Calculate cursor position based on click
			relX := xNorm - i.X - 0.02 // Account for text padding
			var charWidth float32
			if len(i.Text) > 0 {
				charWidth = (i.Width - 0.04) / float32(len(i.Text)+1)
				if charWidth > i.Height*0.7 {
					charWidth = i.Height * 0.7
				}
			} else {
				// When no text, use a reasonable default character width
				charWidth = i.Height * 0.5
			}
			i.CursorPos = int(relX / charWidth)
			if i.CursorPos > len(i.Text) {
				i.CursorPos = len(i.Text)
			}
			if i.CursorPos < 0 {
				i.CursorPos = 0
			}
		} else {
			i.IsFocused = false // Lose focus when clicking outside
		}
	}
}

func (i *Input) HandleCursorPos(x, y float64, width, height int) {
	// Inputs don't handle cursor
}

func (i *Input) GetBounds() (x, y, w, h float32) {
	return i.X, i.Y, i.Width, i.Height
}

func (i *Input) SetPosition(x, y float32) {
	i.X, i.Y = x, y
}

func (i *Input) StopDragging() {
	// Inputs don't drag
}

func (i *Input) HandleKey(key Key, action Action, mods ModifierKey) {
	if !i.IsFocused || action != Press {
		return
	}
	if key >= KeyA && key <= KeyZ {
		char := string(rune(key))
		if mods&ModShift != 0 {
			char = strings.ToUpper(char)
		}
		i.Text = i.Text[:i.CursorPos] + char + i.Text[i.CursorPos:]
		i.CursorPos++
	} else if key == KeySpace {
		// Handle space key
		i.Text = i.Text[:i.CursorPos] + " " + i.Text[i.CursorPos:]
		i.CursorPos++
	} else if key >= 48 && key <= 57 {
		// Handle number keys (ASCII 48-57 = '0'-'9')
		char := string(rune(key))
		i.Text = i.Text[:i.CursorPos] + char + i.Text[i.CursorPos:]
		i.CursorPos++
	} else if key == KeyBackspace && i.CursorPos > 0 {
		i.Text = i.Text[:i.CursorPos-1] + i.Text[i.CursorPos:]
		i.CursorPos--
	} else if key == KeyLeft && i.CursorPos > 0 {
		i.CursorPos--
	} else if key == KeyRight && i.CursorPos < len(i.Text) {
		i.CursorPos++
	} else if key == KeyEnter {
		// For single line, perhaps submit or something, but for now, do nothing
	}
}
