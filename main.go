package main

import (
	"log"

	"xp-pen-controller/internal/ui"
)

func main() {
	// Create the whiteboard window
	window := ui.NewWhiteboardWindow()

	// Try to connect to the tablet
	err := window.ConnectTablet()
	if err != nil {
		log.Printf("Warning: Failed to connect to XP-Pen tablet: %v", err)
		log.Println("The application will still work, but tablet input will not be available.")
	} else {
		log.Println("Successfully connected to XP-Pen tablet")
	}

	// Show the window (this blocks until the window is closed)
	window.Show()
}
