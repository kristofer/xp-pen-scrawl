package drawing

import (
	"image/color"
)

// Point represents a point in the drawing with coordinates and pressure
type Point struct {
	X, Y     float64 // Normalized coordinates (0.0 to 1.0)
	Pressure float64 // Pressure value (0.0 to 1.0)
}

// Stroke represents a continuous drawing stroke
type Stroke struct {
	Points    []Point
	Color     color.Color // Always black per specification
	MinWidth  float64     // Minimum line width
	MaxWidth  float64     // Maximum line width based on pressure
	Completed bool        // Whether the stroke is finished
}

// NewStroke creates a new stroke with the specified color and width range
func NewStroke() *Stroke {
	return &Stroke{
		Points:    make([]Point, 0),
		Color:     color.RGBA{0, 0, 0, 255}, // Black
		MinWidth:  1.0,
		MaxWidth:  8.0,
		Completed: false,
	}
}

// AddPoint adds a new point to the stroke
func (s *Stroke) AddPoint(point Point) {
	s.Points = append(s.Points, point)
}

// GetWidth calculates the line width at a given pressure
func (s *Stroke) GetWidth(pressure float64) float64 {
	// Linear interpolation between min and max width based on pressure
	return s.MinWidth + (s.MaxWidth-s.MinWidth)*pressure
}

// Complete marks the stroke as finished
func (s *Stroke) Complete() {
	s.Completed = true
}

// IsEmpty returns true if the stroke has no points
func (s *Stroke) IsEmpty() bool {
	return len(s.Points) == 0
}

// Canvas represents the drawing surface
type Canvas struct {
	Strokes       []*Stroke
	CurrentStroke *Stroke
	Width         float64
	Height        float64
	Background    color.Color
}

// NewCanvas creates a new canvas with the specified dimensions
func NewCanvas(width, height float64) *Canvas {
	return &Canvas{
		Strokes:       make([]*Stroke, 0),
		CurrentStroke: nil,
		Width:         width,
		Height:        height,
		Background:    color.RGBA{255, 255, 255, 255}, // White background
	}
}

// StartStroke begins a new stroke at the given point
func (c *Canvas) StartStroke(point Point) {
	c.CurrentStroke = NewStroke()
	c.CurrentStroke.AddPoint(point)
}

// AddPointToCurrentStroke adds a point to the current stroke
func (c *Canvas) AddPointToCurrentStroke(point Point) {
	if c.CurrentStroke != nil {
		c.CurrentStroke.AddPoint(point)
	}
}

// FinishStroke completes the current stroke and adds it to the canvas
func (c *Canvas) FinishStroke() {
	if c.CurrentStroke != nil && !c.CurrentStroke.IsEmpty() {
		c.CurrentStroke.Complete()
		c.Strokes = append(c.Strokes, c.CurrentStroke)
	}
	c.CurrentStroke = nil
}

// Clear removes all strokes from the canvas
func (c *Canvas) Clear() {
	c.Strokes = make([]*Stroke, 0)
	c.CurrentStroke = nil
}

// GetAllStrokes returns all completed strokes plus the current stroke if active
func (c *Canvas) GetAllStrokes() []*Stroke {
	allStrokes := make([]*Stroke, len(c.Strokes))
	copy(allStrokes, c.Strokes)

	if c.CurrentStroke != nil && !c.CurrentStroke.IsEmpty() {
		allStrokes = append(allStrokes, c.CurrentStroke)
	}

	return allStrokes
}
