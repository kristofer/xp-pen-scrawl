package tablet

import "xp-pen-controller/internal/drawing"

// CoordinateMapper handles transformation between tablet and screen coordinates
type CoordinateMapper struct {
	tabletMaxX, tabletMaxY    int     // Tablet coordinate bounds
	screenWidth, screenHeight float64 // Screen dimensions
}

// NewCoordinateMapper creates a new coordinate mapper
func NewCoordinateMapper(tabletMaxX, tabletMaxY int, screenWidth, screenHeight float64) *CoordinateMapper {
	return &CoordinateMapper{
		tabletMaxX:   tabletMaxX,
		tabletMaxY:   tabletMaxY,
		screenWidth:  screenWidth,
		screenHeight: screenHeight,
	}
}

// TabletToScreen converts tablet coordinates to normalized screen coordinates
func (cm *CoordinateMapper) TabletToScreen(tabletX, tabletY int) (float64, float64) {
	// Normalize tablet coordinates to 0.0-1.0 range
	normalizedX := float64(tabletX) / float64(cm.tabletMaxX)
	normalizedY := float64(tabletY) / float64(cm.tabletMaxY)

	// Clamp to valid range
	if normalizedX < 0 {
		normalizedX = 0
	} else if normalizedX > 1 {
		normalizedX = 1
	}

	if normalizedY < 0 {
		normalizedY = 0
	} else if normalizedY > 1 {
		normalizedY = 1
	}

	return normalizedX, normalizedY
}

// NormalizePressure converts raw pressure to normalized 0.0-1.0 range
func (cm *CoordinateMapper) NormalizePressure(rawPressure int) float64 {
	// XP-Pen tablets typically have pressure range 0-8191
	const maxPressure = 8191

	pressure := float64(rawPressure) / float64(maxPressure)

	// Clamp to valid range
	if pressure < 0 {
		pressure = 0
	} else if pressure > 1 {
		pressure = 1
	}

	return pressure
}

// PenDataToPoint converts raw pen data to a drawing point
func (cm *CoordinateMapper) PenDataToPoint(penData *PenData) drawing.Point {
	x, y := cm.TabletToScreen(penData.X, penData.Y)
	pressure := cm.NormalizePressure(penData.Pressure)

	return drawing.Point{
		X:        x,
		Y:        y,
		Pressure: pressure,
	}
}
