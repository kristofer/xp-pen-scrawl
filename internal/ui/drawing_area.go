package ui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"

	"xp-pen-controller/internal/drawing"
)

// DrawingArea is a custom widget that displays the drawing canvas and handles input
type DrawingArea struct {
	widget.BaseWidget
	canvas      *drawing.Canvas
	lines       []*canvas.Line
	needsUpdate bool
	isDragging  bool // Track if we're currently dragging
}

// Ensure DrawingArea implements the required interfaces
var _ fyne.Widget = (*DrawingArea)(nil)
var _ fyne.Tappable = (*DrawingArea)(nil)
var _ fyne.Draggable = (*DrawingArea)(nil)
var _ fyne.Focusable = (*DrawingArea)(nil)

// Note: MouseDown/MouseUp might not be standard Fyne interfaces

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

// Tapped handles tap events on the drawing area
func (da *DrawingArea) Tapped(event *fyne.PointEvent) {
	fmt.Printf("DEBUG: Tapped event at position: %v (ignoring - waiting for stylus input)\n", event.Position)
	// Do nothing - we only want to draw with stylus input, not mouse clicks
}

// MouseDown handles mouse press events (start drawing)
func (da *DrawingArea) MouseDown(event *fyne.PointEvent) {
	fmt.Printf("DEBUG: Mouse down at: %v\n", event.Position)

	// Set dragging flag
	da.isDragging = true

	// For testing: allow drawing with mouse (while we fix tablet permissions)
	// Convert screen coordinates to normalized coordinates
	normalizedX := float64(event.Position.X) / float64(da.Size().Width)
	normalizedY := float64(event.Position.Y) / float64(da.Size().Height)

	// Create a drawing point
	point := drawing.Point{
		X:        normalizedX,
		Y:        normalizedY,
		Pressure: 0.7, // Default pressure for mouse
	}

	// Start a new stroke
	fmt.Println("DEBUG: Starting new stroke with mouse")
	da.canvas.StartStroke(point)
	da.needsUpdate = true
	da.BaseWidget.Refresh()
}

// MouseUp handles mouse release events (stop drawing)
func (da *DrawingArea) MouseUp(event *fyne.PointEvent) {
	fmt.Printf("DEBUG: Mouse up at: %v (isDragging: %t)\n", event.Position, da.isDragging)

	if da.isDragging {
		da.isDragging = false
		// Finish the current stroke
		fmt.Println("DEBUG: Finishing stroke on mouse up")
		da.canvas.FinishStroke()
		da.needsUpdate = true
		da.BaseWidget.Refresh()
	}
}

// Dragged handles mouse drag events (continue drawing)
func (da *DrawingArea) Dragged(event *fyne.DragEvent) {
	fmt.Printf("DEBUG: Dragged event at: %v (isDragging: %t)\n", event.Position, da.isDragging)

	// If not already dragging, start a new stroke
	if !da.isDragging {
		fmt.Println("DEBUG: Starting drag - initializing stroke")
		da.isDragging = true

		// Convert screen coordinates to normalized coordinates for start point
		size := da.Size()
		normalizedX := float64(event.Position.X) / float64(size.Width)
		normalizedY := float64(event.Position.Y) / float64(size.Height)

		// Create the first point
		point := drawing.Point{
			X:        normalizedX,
			Y:        normalizedY,
			Pressure: 0.7, // Default pressure for mouse
		}

		// Start a new stroke
		fmt.Println("DEBUG: Starting new stroke with drag")
		da.canvas.StartStroke(point)
	} else {
		// Continue existing stroke
		// Convert screen coordinates to normalized coordinates
		size := da.Size()
		normalizedX := float64(event.Position.X) / float64(size.Width)
		normalizedY := float64(event.Position.Y) / float64(size.Height)

		// Create a new point
		point := drawing.Point{
			X:        normalizedX,
			Y:        normalizedY,
			Pressure: 0.7, // Higher pressure for visible line
		}

		// Add point to current stroke
		fmt.Println("DEBUG: Adding point to stroke during drag")
		da.canvas.AddPointToCurrentStroke(point)
	}

	// IMPORTANT: Force immediate refresh to show the stroke being drawn
	da.needsUpdate = true
	da.BaseWidget.Refresh() // This calls the renderer's Refresh() method
}

// DragEnd handles end of drag events
func (da *DrawingArea) DragEnd() {
	fmt.Printf("DEBUG: DragEnd called (isDragging: %t)\n", da.isDragging)

	if da.isDragging {
		da.isDragging = false
		// Finish the current stroke
		fmt.Println("DEBUG: Finishing stroke on drag end")
		da.canvas.FinishStroke()
		da.needsUpdate = true
		da.BaseWidget.Refresh()
	}
}

// FocusGained is called when the widget gains focus
func (da *DrawingArea) FocusGained() {
	// Nothing special needed when focus is gained
}

// FocusLost is called when the widget loses focus
func (da *DrawingArea) FocusLost() {
	// If we're in the middle of drawing, finish the stroke
	if da.isDragging {
		da.isDragging = false
		da.canvas.FinishStroke()
		da.BaseWidget.Refresh()
	}
}

// TypedRune handles typed characters
func (da *DrawingArea) TypedRune(rune) {
	// Not needed for drawing
}

// TypedKey handles special key presses
func (da *DrawingArea) TypedKey(*fyne.KeyEvent) {
	// Not needed for drawing
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
	// Only resize and position the background rectangle
	// Leave drawing objects (lines/circles) in their original positions
	if len(r.objects) > 0 {
		// First object should be the background
		bg := r.objects[0]
		bg.Resize(size)
		bg.Move(fyne.NewPos(0, 0))
	}
	// Drawing objects (lines, circles) keep their absolute positions
	// Don't modify them in Layout as they have specific coordinates
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
	fmt.Printf("DEBUG: Renderer Refresh called (needsUpdate: %t)\n", r.area.needsUpdate)

	if !r.area.needsUpdate {
		fmt.Println("DEBUG: Skipping refresh - needsUpdate is false")
		return
	}

	// Clear existing objects
	r.objects = make([]fyne.CanvasObject, 0)

	// Create background
	bg := canvas.NewRectangle(color.RGBA{255, 255, 255, 255}) // White background
	r.objects = append(r.objects, bg)

	// Render all strokes
	strokes := r.area.canvas.GetAllStrokes()
	fmt.Printf("DEBUG: Rendering %d strokes\n", len(strokes))

	for i, stroke := range strokes {
		fmt.Printf("DEBUG: Rendering stroke %d with %d points\n", i+1, len(stroke.Points))
		r.renderStroke(stroke)
	}

	r.area.needsUpdate = false
	fmt.Println("DEBUG: Refresh complete")
}

// renderStroke renders a single stroke as a series of lines
func (r *drawingAreaRenderer) renderStroke(stroke *drawing.Stroke) {
	fmt.Printf("DEBUG: renderStroke called with %d points\n", len(stroke.Points))

	if len(stroke.Points) == 0 {
		fmt.Println("DEBUG: No points to draw")
		return // No points to draw
	}

	// Get the size of the drawing area
	size := r.area.Size()
	fmt.Printf("DEBUG: Drawing area size: %v\n", size)

	if size.Width == 0 || size.Height == 0 {
		fmt.Println("DEBUG: Invalid drawing area size")
		return
	}

	// Handle single point (dot)
	if len(stroke.Points) == 1 {
		fmt.Println("DEBUG: Rendering single point as circle")
		point := stroke.Points[0]
		x := float32(point.X * float64(size.Width))
		y := float32(point.Y * float64(size.Height))

		// Create a small circle for single points
		circle := canvas.NewCircle(stroke.Color)
		radius := float32(stroke.GetWidth(point.Pressure)) / 2
		circle.Resize(fyne.NewSize(radius*2, radius*2))
		circle.Move(fyne.NewPos(x-radius, y-radius))
		circle.FillColor = stroke.Color

		r.objects = append(r.objects, circle)
		fmt.Printf("DEBUG: Added circle at (%f, %f) with radius %f\n", x, y, radius)
		return
	}

	// Draw lines between consecutive points
	fmt.Printf("DEBUG: Rendering %d line segments\n", len(stroke.Points)-1)
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
		fmt.Printf("DEBUG: Added line from (%f,%f) to (%f,%f) width %f\n",
			x1, y1, x2, y2, line.StrokeWidth)
	}
	fmt.Printf("DEBUG: renderStroke complete - total objects: %d\n", len(r.objects))
}

// Destroy cleans up the renderer
func (r *drawingAreaRenderer) Destroy() {
	// Nothing special needed for cleanup
}
