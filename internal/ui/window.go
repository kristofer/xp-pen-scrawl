package ui

import (
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
	clearButton := widget.NewButton("Clear", func() {
		ww.canvas.Clear()
		ww.drawingArea.Refresh()
	})

	saveButton := widget.NewButton("Save", func() {
		// TODO: Implement save functionality
		// For now, just show a placeholder
		dialog.ShowInformation("Save", "Save functionality coming soon!", ww.window)
	})

	quitButton := widget.NewButton("Quit", func() {
		ww.Close()
	})

	// Create toolbar with minimal buttons
	toolbar := container.NewHBox(
		clearButton,
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
	for ww.tablet.IsConnected() {
		penData, err := ww.tablet.ReadPenData()
		if err != nil {
			continue // Skip errors and keep trying
		}

		// Convert to drawing point
		point := ww.mapper.PenDataToPoint(penData)

		// Handle pen input
		if penData.PenDown {
			if ww.canvas.CurrentStroke == nil {
				ww.canvas.StartStroke(point)
			} else {
				ww.canvas.AddPointToCurrentStroke(point)
			}
		} else if ww.canvas.CurrentStroke != nil {
			// Pen lifted, finish stroke
			ww.canvas.FinishStroke()
		}

		// Refresh drawing area
		ww.drawingArea.Refresh()
	}
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
