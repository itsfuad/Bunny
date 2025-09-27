package ui

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Slider struct {
	X, Y, Width, Height float32
	Value               float32 // 0 to 1
	TrackColor          [4]float32
	KnobColor           [4]float32
	OnValueChange       func(float32)
	TrackVAO            uint32
	KnobVAO             uint32
	Dragging            bool
	KnobWidth           float32
}

func NewSlider(x, y, w, h float32, trackColor, knobColor [4]float32) *Slider {
	knobW := h * 2 // knob width is 2x height for visibility
	var trackVAO, trackVBO uint32
	gl.GenVertexArrays(1, &trackVAO)
	gl.GenBuffers(1, &trackVBO)
	gl.BindVertexArray(trackVAO)
	vertices := []float32{
		x, y,
		x + w, y,
		x + w, y + h,
		x, y + h,
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, trackVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, nil)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	var knobVAO, knobVBO uint32
	gl.GenVertexArrays(1, &knobVAO)
	gl.GenBuffers(1, &knobVBO)
	gl.BindVertexArray(knobVAO)
	// Knob vertices will be updated in Draw or Update
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	return &Slider{
		X: x, Y: y, Width: w, Height: h,
		Value:      0.5,
		TrackColor: trackColor,
		KnobColor:  knobColor,
		TrackVAO:   trackVAO,
		KnobVAO:    knobVAO,
		KnobWidth:  knobW,
	}
}

func (s *Slider) Draw(program uint32) {
	// Draw track
	gl.UseProgram(program)
	colorLoc := gl.GetUniformLocation(program, gl.Str("color\x00"))
	gl.Uniform4f(colorLoc, s.TrackColor[0], s.TrackColor[1], s.TrackColor[2], s.TrackColor[3])
	gl.BindVertexArray(s.TrackVAO)
	gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)
	gl.BindVertexArray(0)

	// Draw knob
	knobX := s.X + s.Value*(s.Width-s.KnobWidth)
	vertices := []float32{
		knobX, s.Y,
		knobX + s.KnobWidth, s.Y,
		knobX + s.KnobWidth, s.Y + s.Height,
		knobX, s.Y + s.Height,
	}
	gl.BindVertexArray(s.KnobVAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0) // Assuming VBO is bound, but to update, need to bind and update data
	// For simplicity, regenerate VBO each time, but better to update buffer
	var knobVBO uint32
	gl.GenBuffers(1, &knobVBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, knobVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.DYNAMIC_DRAW)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, nil)
	gl.EnableVertexAttribArray(0)
	gl.Uniform4f(colorLoc, s.KnobColor[0], s.KnobColor[1], s.KnobColor[2], s.KnobColor[3])
	gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	gl.DeleteBuffers(1, &knobVBO) // Clean up, but inefficient
}

func (s *Slider) Update(value float32) {
	if value < 0 {
		value = 0
	} else if value > 1 {
		value = 1
	}
	s.Value = value
	if s.OnValueChange != nil {
		s.OnValueChange(s.Value)
	}
}

func (s *Slider) HandleMouse(x, y float64, action glfw.Action, width, height int) {
	xNorm := float32((x/float64(width))*2 - 1)
	yNorm := float32(1 - (y/float64(height))*2)
	knobX := s.X + s.Value*(s.Width-s.KnobWidth)
	if action == glfw.Press && xNorm >= knobX && xNorm <= knobX+s.KnobWidth && yNorm >= s.Y && yNorm <= s.Y+s.Height {
		s.Dragging = true
	} else if action == glfw.Release {
		s.Dragging = false
	}
}

func (s *Slider) GetBounds() (x, y, w, h float32) {
	return s.X, s.Y, s.Width, s.Height
}

func (s *Slider) SetPosition(x, y float32) {
	s.X, s.Y = x, y
	// Recreate TrackVAO with new vertices
	gl.DeleteVertexArrays(1, &s.TrackVAO)
	var trackVAO, trackVBO uint32
	gl.GenVertexArrays(1, &trackVAO)
	gl.GenBuffers(1, &trackVBO)
	gl.BindVertexArray(trackVAO)
	vertices := []float32{
		x, y,
		x + s.Width, y,
		x + s.Width, y + s.Height,
		x, y + s.Height,
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, trackVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 2*4, nil)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	s.TrackVAO = trackVAO
}

func (s *Slider) HandleCursorPos(x, y float64, width, height int) {
	if s.Dragging {
		xNorm := float32((x/float64(width))*2 - 1)
		relX := xNorm - s.X
		if relX < 0 {
			relX = 0
		} else if relX > s.Width-s.KnobWidth {
			relX = s.Width - s.KnobWidth
		}
		newValue := relX / (s.Width - s.KnobWidth)
		s.Update(newValue)
	}
}
