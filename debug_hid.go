package main

import (
	"fmt"

	"github.com/karalabe/hid"
)

func debugMain() {
	fmt.Println("Scanning for all HID devices...")

	// List all HID devices
	devices := hid.Enumerate(0, 0) // 0, 0 means all devices

	fmt.Printf("Found %d HID devices:\n", len(devices))

	for i, device := range devices {
		fmt.Printf("%d. VendorID: 0x%04x, ProductID: 0x%04x\n", i+1, device.VendorID, device.ProductID)
		fmt.Printf("   Manufacturer: %s\n", device.Manufacturer)
		fmt.Printf("   Product: %s\n", device.Product)
		fmt.Printf("   Serial: %s\n", device.Serial)
		fmt.Printf("   Path: %s\n", device.Path)
		fmt.Println()
	}

	// Specifically look for XP-Pen devices
	fmt.Println("Looking specifically for XP-Pen devices...")
	xpPenDevices := hid.Enumerate(0x28bd, 0) // XP-Pen vendor ID, any product

	if len(xpPenDevices) > 0 {
		fmt.Printf("Found %d XP-Pen devices:\n", len(xpPenDevices))
		for i, device := range xpPenDevices {
			fmt.Printf("%d. ProductID: 0x%04x\n", i+1, device.ProductID)
			fmt.Printf("   Product: %s\n", device.Product)
		}
	} else {
		fmt.Println("No XP-Pen devices found with vendor ID 0x28bd")
	}
}
