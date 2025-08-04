package tablet

import (
	"fmt"

	"github.com/karalabe/hid"
)

const (
	// XP-Pen Star G640 identifiers
	VendorID  = 0x28bd // XP-Pen vendor ID
	ProductID = 0x0094 // Star G640 product ID
)

// PenData represents the state of the pen at a given moment
type PenData struct {
	X        int  // X coordinate (0-32767)
	Y        int  // Y coordinate (0-32767)
	Pressure int  // Pressure level (0-8191)
	PenDown  bool // Whether pen is touching the tablet
	InRange  bool // Whether pen is in proximity to tablet
}

// TabletController handles communication with the XP-Pen tablet
type TabletController struct {
	device *hid.Device
	active bool
}

// NewTabletController creates a new tablet controller
func NewTabletController() *TabletController {
	return &TabletController{
		active: false,
	}
}

// Connect establishes connection to the XP-Pen tablet
func (tc *TabletController) Connect() error {
	devices := hid.Enumerate(VendorID, ProductID)
	if len(devices) == 0 {
		return fmt.Errorf("XP-Pen Star G640 not found")
	}

	device, err := devices[0].Open()
	if err != nil {
		return fmt.Errorf("failed to open device: %w", err)
	}

	tc.device = device
	tc.active = true
	return nil
}

// Disconnect closes the connection to the tablet
func (tc *TabletController) Disconnect() error {
	if tc.device != nil {
		tc.active = false
		return tc.device.Close()
	}
	return nil
}

// IsConnected returns whether the tablet is currently connected
func (tc *TabletController) IsConnected() bool {
	return tc.active && tc.device != nil
}

// ReadPenData reads the current pen state from the tablet
func (tc *TabletController) ReadPenData() (*PenData, error) {
	if !tc.IsConnected() {
		return nil, fmt.Errorf("tablet not connected")
	}

	// Read raw data from the tablet
	data := make([]byte, 12) // XP-Pen reports are typically 8-12 bytes
	n, err := tc.device.Read(data)
	if err != nil {
		return nil, fmt.Errorf("failed to read from tablet: %w", err)
	}

	if n < 8 {
		return nil, fmt.Errorf("insufficient data received: %d bytes", n)
	}

	// Parse the pen data based on XP-Pen protocol
	// Note: This parsing may need adjustment based on actual tablet output
	penData := &PenData{
		X:        int(data[2]) | (int(data[3]) << 8),
		Y:        int(data[4]) | (int(data[5]) << 8),
		Pressure: int(data[6]) | (int(data[7]) << 8),
		PenDown:  (data[1] & 0x01) != 0,
		InRange:  (data[1] & 0x02) != 0,
	}

	return penData, nil
}

// GetTabletDimensions returns the tablet's maximum coordinates
func (tc *TabletController) GetTabletDimensions() (int, int) {
	// XP-Pen Star G640 specifications
	return 32767, 32767 // Maximum X, Y coordinates
}
