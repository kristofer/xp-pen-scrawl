package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"

	"xp-pen-controller/internal/drawing"
)

// DrawingArea is a custom widget that displays the drawing canvas
type DrawingArea struct {
	widget.BaseWidget
	canvas      *drawing.Canvas
	lines       []*canvas.Line
	needsUpdate bool
	isDragging  bool // Track if we're currently dragging
}

// NewDrawingArea creates a new drawing area widget
func NewDrawingArea(drawingCanvas *drawing.Canvas) *DrawingArea {
	area := &DrawingArea{
		canvas:      drawingCanvas,
		lines:       make([]*canvas.Line, 0),
		needsUpdate: true,
		isDragging:  false,
	}
	area.ExtendBaseWidget(area)
	return area
}

// Refresh updates the drawing area display
func (da *DrawingArea) Refresh() {
	da.needsUpdate = true
	da.BaseWidget.Refresh()
}

// Tapped handles tap events (single click to start drawing)
func (da *DrawingArea) Tapped(event *fyne.PointEvent) {
	// Note: This method handles single taps/clicks when the pen/mouse is pressed down
	// It does NOT fire for hover events - only actual press events
	// Single tap starts and immediately ends a stroke (for dots)
	size := da.Size()
	normalizedX := float64(event.Position.X) / float64(size.Width)
	normalizedY := float64(event.Position.Y) / float64(size.Height)

	point := drawing.Point{
		X:        normalizedX,
		Y:        normalizedY,
		Pressure: 0.5,
	}

	// Create a single point stroke
	da.canvas.StartStroke(point)
	da.canvas.FinishStroke()
	da.Refresh()
}

// MouseDown handles mouse press events (start drawing)
func (da *DrawingArea) MouseDown(event *fyne.PointEvent) {
	// This only fires when the mouse/stylus button is actually pressed down
	// NOT when hovering - ensures drawing only happens on contact
	da.isDragging = true

	// Convert screen coordinates to normalized coordinates
	size := da.Size()
	normalizedX := float64(event.Position.X) / float64(size.Width)
	normalizedY := float64(event.Position.Y) / float64(size.Height)

	// Create a new point with default pressure for mouse input
	point := drawing.Point{
		X:        normalizedX,
		Y:        normalizedY,
		Pressure: 0.5, // Default pressure for mouse/stylus without pressure sensitivity
	}

	// Start a new stroke
	da.canvas.StartStroke(point)
	da.Refresh()
}

// MouseUp handles mouse release events (stop drawing)
func (da *DrawingArea) MouseUp(event *fyne.PointEvent) {
	if da.isDragging {
		da.isDragging = false
		// Finish the current stroke
		da.canvas.FinishStroke()
		da.Refresh()
	}
}

// Dragged handles mouse drag events (continue drawing)
func (da *DrawingArea) Dragged(event *fyne.DragEvent) {
	if da.isDragging {
		// Convert screen coordinates to normalized coordinates
		size := da.Size()
		normalizedX := float64(event.Position.X) / float64(size.Width)
		normalizedY := float64(event.Position.Y) / float64(size.Height)

		// Create a new point
		point := drawing.Point{
			X:        normalizedX,
			Y:        normalizedY,
			Pressure: 0.5, // Default pressure for mouse/stylus
		}

		// Add point to current stroke
		da.canvas.AddPointToCurrentStroke(point)
		da.Refresh()
	}
}

// DragEnd handles end of drag events
func (da *DrawingArea) DragEnd() {
	if da.isDragging {
		da.isDragging = false
		// Finish the current stroke
		da.canvas.FinishStroke()
		da.Refresh()
	}
}

// CreateRenderer creates the renderer for this widget
func (da *DrawingArea) CreateRenderer() fyne.WidgetRenderer {
	return &drawingAreaRenderer{
		area: da,
	}
}

// drawingAreaRenderer renders the drawing area
type drawingAreaRenderer struct {
	area    *DrawingArea
	objects []fyne.CanvasObject
}

// Layout arranges the objects in the renderer
func (r *drawingAreaRenderer) Layout(size fyne.Size) {
	// Drawing area fills the entire available space
	for _, obj := range r.objects {
		obj.Resize(size)
		obj.Move(fyne.NewPos(0, 0))
	}
}

// MinSize returns the minimum size for the drawing area
func (r *drawingAreaRenderer) MinSize() fyne.Size {
	return fyne.NewSize(400, 300) // Minimum reasonable size
}

// Objects returns all the objects in the renderer
func (r *drawingAreaRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// Refresh updates the renderer display
func (r *drawingAreaRenderer) Refresh() {
	if !r.area.needsUpdate {
		return
	}

	// Clear existing objects
	r.objects = make([]fyne.CanvasObject, 0)

	// Create background
	bg := canvas.NewRectangle(color.RGBA{255, 255, 255, 255}) // White background
	r.objects = append(r.objects, bg)

	// Render all strokes
	strokes := r.area.canvas.GetAllStrokes()
	for _, stroke := range strokes {
		r.renderStroke(stroke)
	}

	r.area.needsUpdate = false
}

// renderStroke renders a single stroke as a series of lines
func (r *drawingAreaRenderer) renderStroke(stroke *drawing.Stroke) {
	if len(stroke.Points) < 2 {
		return // Need at least 2 points to draw a line
	}

	// Get the size of the drawing area
	size := r.area.Size()
	if size.Width == 0 || size.Height == 0 {
		return
	}

	// Draw lines between consecutive points
	for i := 0; i < len(stroke.Points)-1; i++ {
		p1 := stroke.Points[i]
		p2 := stroke.Points[i+1]

		// Convert normalized coordinates to screen coordinates
		x1 := float32(p1.X * float64(size.Width))
		y1 := float32(p1.Y * float64(size.Height))
		x2 := float32(p2.X * float64(size.Width))
		y2 := float32(p2.Y * float64(size.Height))

		// Create line
		line := canvas.NewLine(stroke.Color)
		line.Position1 = fyne.NewPos(x1, y1)
		line.Position2 = fyne.NewPos(x2, y2)

		// Set line width based on pressure (average of the two points)
		avgPressure := (p1.Pressure + p2.Pressure) / 2.0
		line.StrokeWidth = float32(stroke.GetWidth(avgPressure))

		r.objects = append(r.objects, line)
	}
}

// Destroy cleans up the renderer
func (r *drawingAreaRenderer) Destroy() {
	// Nothing special needed for cleanup
}
