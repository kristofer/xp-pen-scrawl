package tablet

import (
	"fmt"

	"github.com/karalabe/hid"
)

const (
	// XP-Pen Star G640 identifiers
	VendorID  = 0x28bd // XP-Pen vendor ID
	ProductID = 0x0914 // 6 inch PenTablet product ID (was 0x0094)
)

// PenData represents the state of the pen at a given moment
type PenData struct {
	X        int  // X coordinate (0-32767)
	Y        int  // Y coordinate (0-32767)
	Pressure int  // Pressure level (0-8191)
	PenDown  bool // Whether pen is touching the tablet
	InRange  bool // Whether pen is in proximity to tablet
	Button1  bool // First pen button pressed
	Button2  bool // Second pen button pressed
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
		return fmt.Errorf("XP-Pen tablet not found")
	}

	fmt.Printf("DEBUG: Found %d tablet devices, trying to connect...\n", len(devices))

	// Print detailed info about each device
	for i, deviceInfo := range devices {
		fmt.Printf("DEBUG: Device %d details:\n", i+1)
		fmt.Printf("  VendorID: 0x%04x\n", deviceInfo.VendorID)
		fmt.Printf("  ProductID: 0x%04x\n", deviceInfo.ProductID)
		fmt.Printf("  Manufacturer: '%s'\n", deviceInfo.Manufacturer)
		fmt.Printf("  Product: '%s'\n", deviceInfo.Product)
		fmt.Printf("  Serial: '%s'\n", deviceInfo.Serial)
		fmt.Printf("  Path: '%s'\n", deviceInfo.Path)
		fmt.Printf("  Interface: %d\n", deviceInfo.Interface)
		fmt.Printf("  Usage Page: 0x%04x\n", deviceInfo.UsagePage)
		fmt.Printf("  Usage: 0x%04x\n", deviceInfo.Usage)
		fmt.Println()
	}

	// Try to open each device until we find one that works
	// Prefer digitizer devices (Usage Page 0x000d) first
	var lastError error

	// First pass: try digitizer devices only
	for i, deviceInfo := range devices {
		if deviceInfo.UsagePage == 0x000d { // Digitizer usage page
			fmt.Printf("DEBUG: Trying DIGITIZER device %d (Interface %d, UsagePage 0x%04x, Usage 0x%04x)...\n",
				i+1, deviceInfo.Interface, deviceInfo.UsagePage, deviceInfo.Usage)

			device, err := deviceInfo.Open()
			if err != nil {
				fmt.Printf("DEBUG: Failed to open digitizer device %d: %v\n", i+1, err)
				lastError = err
				continue
			}

			// Successfully opened
			tc.device = device
			tc.active = true
			fmt.Printf("DEBUG: Successfully connected to DIGITIZER device %d\n", i+1)
			return nil
		}
	}

	// Second pass: try any remaining devices
	for i, deviceInfo := range devices {
		if deviceInfo.UsagePage != 0x000d { // Skip digitizer devices (already tried)
			fmt.Printf("DEBUG: Trying OTHER device %d (Interface %d, UsagePage 0x%04x, Usage 0x%04x)...\n",
				i+1, deviceInfo.Interface, deviceInfo.UsagePage, deviceInfo.Usage)

			device, err := deviceInfo.Open()
			if err != nil {
				fmt.Printf("DEBUG: Failed to open device %d: %v\n", i+1, err)
				lastError = err
				continue
			}

			// Successfully opened
			tc.device = device
			tc.active = true
			fmt.Printf("DEBUG: Successfully connected to device %d\n", i+1)
			return nil
		}
	}

	return fmt.Errorf("failed to open any tablet device: %w", lastError)
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

	// Debug: Print raw data bytes
	fmt.Printf("DEBUG: Raw tablet data (%d bytes): %02x %02x %02x %02x %02x %02x %02x %02x\n",
		n, data[0], data[1], data[2], data[3], data[4], data[5], data[6], data[7])

	// Parse the pen data based on XP-Pen protocol
	// Note: This parsing may need adjustment based on actual tablet output
	penData := &PenData{
		X:        int(data[2]) | (int(data[3]) << 8),
		Y:        int(data[4]) | (int(data[5]) << 8),
		Pressure: int(data[6]) | (int(data[7]) << 8),
		PenDown:  (data[1] & 0x01) != 0,
		InRange:  (data[1] & 0x02) != 0,
		Button1:  (data[1] & 0x04) != 0, // Check bit 2 for button 1
		Button2:  (data[1] & 0x08) != 0, // Check bit 3 for button 2
	}

	return penData, nil
}

// GetTabletDimensions returns the tablet's maximum coordinates
func (tc *TabletController) GetTabletDimensions() (int, int) {
	// XP-Pen Star G640 specifications
	return 32767, 32767 // Maximum X, Y coordinates
}
