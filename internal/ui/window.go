package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	"xp-pen-controller/internal/drawing"
	"xp-pen-controller/internal/tablet"
)

// WhiteboardWindow represents the main application window
type WhiteboardWindow struct {
	app         fyne.App
	window      fyne.Window
	canvas      *drawing.Canvas
	drawingArea *DrawingArea
	tablet      *tablet.TabletController
	mapper      *tablet.CoordinateMapper
}

// NewWhiteboardWindow creates a new whiteboard window
func NewWhiteboardWindow() *WhiteboardWindow {
	app := app.New()
	app.SetIcon(nil) // No icon for minimal design

	window := app.NewWindow("XP-Pen Whiteboard")
	window.Resize(fyne.NewSize(1200, 900)) // Set window size to 1200x900
	window.CenterOnScreen()                // Center the window on screen

	// Create canvas with specified dimensions
	drawingCanvas := drawing.NewCanvas(1200, 900)

	// Create tablet controller
	tabletController := tablet.NewTabletController()

	// Create coordinate mapper (will be updated when tablet connects)
	mapper := tablet.NewCoordinateMapper(32767, 32767, 1200, 900)

	ww := &WhiteboardWindow{
		app:    app,
		window: window,
		canvas: drawingCanvas,
		tablet: tabletController,
		mapper: mapper,
	}

	// Create custom drawing area
	ww.drawingArea = NewDrawingArea(drawingCanvas)

	// Setup UI
	ww.setupUI()

	return ww
}

// setupUI creates the user interface layout
func (ww *WhiteboardWindow) setupUI() {
	// Create minimal toolbar
	clearButton := widget.NewButton("Test Draw", func() {
		// Test: Add a stroke manually to verify the drawing system works
		testPoint1 := drawing.Point{X: 0.5, Y: 0.2, Pressure: 0.7}
		testPoint2 := drawing.Point{X: 0.7, Y: 0.6, Pressure: 0.7}
		ww.canvas.StartStroke(testPoint1)
		ww.canvas.AddPointToCurrentStroke(testPoint2)
		ww.canvas.FinishStroke()
		ww.drawingArea.Refresh()
		println("Test stroke added manually")
	})

	saveButton := widget.NewButton("Save", func() {
		// TODO: Implement save functionality
		// For now, just show a placeholder
		dialog.ShowInformation("Save", "Save functionality coming soon!", ww.window)
	})

	clearButton2 := widget.NewButton("Clear", func() {
		ww.canvas.Clear()
		ww.drawingArea.Refresh()
	})

	quitButton := widget.NewButton("Quit", func() {
		ww.Close()
	})

	// Create toolbar with minimal buttons
	toolbar := container.NewHBox(
		clearButton,
		clearButton2,
		saveButton,
		quitButton,
		widget.NewSeparator(),
		widget.NewLabel("XP-Pen Whiteboard"),
	)

	// Create main layout with toolbar at top and drawing area filling the rest
	content := container.NewBorder(
		toolbar,        // top
		nil,            // bottom
		nil,            // left
		nil,            // right
		ww.drawingArea, // center - drawing area fills remaining space
	)

	ww.window.SetContent(content)

	// Make the drawing area focusable and focus it
	ww.window.Canvas().Focus(ww.drawingArea)

	// Setup keyboard shortcuts
	ww.setupKeyboardShortcuts()
}

// setupKeyboardShortcuts configures keyboard shortcuts
func (ww *WhiteboardWindow) setupKeyboardShortcuts() {
	// Clear canvas shortcut (Ctrl+N)
	ww.window.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		if key.Name == fyne.KeyN && key.Physical.ScanCode != 0 {
			// Check for Ctrl modifier (this is simplified, real implementation would need proper modifier detection)
			ww.canvas.Clear()
			ww.drawingArea.Refresh()
		}
	})
}

// ConnectTablet attempts to connect to the XP-Pen tablet
func (ww *WhiteboardWindow) ConnectTablet() error {
	err := ww.tablet.Connect()
	if err != nil {
		return err
	}

	// Update coordinate mapper with actual tablet dimensions
	maxX, maxY := ww.tablet.GetTabletDimensions()
	ww.mapper = tablet.NewCoordinateMapper(maxX, maxY, 1200, 900)

	// Start tablet input processing
	go ww.processTabletInput()

	return nil
}

// processTabletInput continuously reads tablet input
func (ww *WhiteboardWindow) processTabletInput() {
	fmt.Println("DEBUG: Starting tablet input processing...")

	for ww.tablet.IsConnected() {
		penData, err := ww.tablet.ReadPenData()
		if err != nil {
			continue // Skip errors and keep trying
		}

		// Debug output for pen data
		fmt.Printf("DEBUG: Pen data - X:%d Y:%d Pressure:%d PenDown:%t InRange:%t Button1:%t Button2:%t\n",
			penData.X, penData.Y, penData.Pressure, penData.PenDown, penData.InRange, penData.Button1, penData.Button2)

		// Special debug output for button presses
		if penData.Button1 {
			fmt.Println("DEBUG: *** BUTTON 1 PRESSED ***")
		}
		if penData.Button2 {
			fmt.Println("DEBUG: *** BUTTON 2 PRESSED ***")
		}

		// Convert to drawing point
		point := ww.mapper.PenDataToPoint(penData)

		// Handle pen input - only draw when pen is down AND button 1 is pressed
		if penData.PenDown && penData.Button1 {
			fmt.Printf("DEBUG: Drawing! Stylus touching tablet at position (%d, %d) with pressure %d and button pressed\n",
				penData.X, penData.Y, penData.Pressure)

			if ww.canvas.CurrentStroke == nil {
				fmt.Println("DEBUG: Starting new stroke")
				ww.canvas.StartStroke(point)
			} else {
				fmt.Println("DEBUG: Adding point to current stroke")
				ww.canvas.AddPointToCurrentStroke(point)
			}
		} else if ww.canvas.CurrentStroke != nil {
			// Pen lifted or button released, finish stroke
			fmt.Println("DEBUG: Pen lifted or button released, finishing stroke")
			ww.canvas.FinishStroke()
		}

		// Refresh drawing area
		ww.drawingArea.Refresh()
	}

	fmt.Println("DEBUG: Tablet input processing stopped")
}

// Show displays the window
func (ww *WhiteboardWindow) Show() {
	ww.window.ShowAndRun()
}

// Close closes the window and disconnects tablet
func (ww *WhiteboardWindow) Close() {
	if ww.tablet != nil {
		ww.tablet.Disconnect()
	}
	ww.app.Quit()
}
