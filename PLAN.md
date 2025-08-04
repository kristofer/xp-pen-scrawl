# XP-Pen Simple Whiteboard App - Implementation Plan

## Current State Analysis

### What's Already Implemented
- Basic XP-Pen Star G640 device connection using HID library
- Device enumeration and connection logic
- Basic Go module structure with `github.com/karalabe/hid` dependency

### What's Missing (Based on SPEC.md Requirements)
The current `main.go` only handles device connection. We need to build the entire whiteboard application on top of this foundation.

## Implementation Plan

### Phase 1: Core Drawing Engine
- [x] **Add GUI Framework**: Choose and integrate a cross-platform GUI library
  - ✅ Integrated Fyne for cross-platform support with canvas drawing and real-time updates
- [x] **Implement Canvas**: Create drawing surface that fills the screen
  - ✅ Created 1200x900 windowed drawing surface with white background
- [x] **Add Pen Input Handling**: Process tablet input with pressure sensitivity
  - ✅ Implemented pen data processing with pressure and coordinate mapping
  - ✅ Added mouse/stylus GUI input handling (only draws on contact, not hover)
- [x] **Implement Stroke Rendering**: Draw smooth lines with variable thickness
  - ✅ Created stroke rendering system with pressure-based thickness

### Phase 2: Tablet Input Processing
- [x] **Complete TabletController**: Extend with methods to:
  - ✅ `ReadPenData()` - Read raw input from tablet
  - ✅ `ParsePosition()` - Extract X, Y coordinates  
  - ✅ `GetPressure()` - Extract pressure values
  - ✅ `IsPenDown()` - Detect pen contact state
- [x] **Coordinate Mapping**: Transform tablet coordinates to screen coordinates
  - ✅ Implemented CoordinateMapper with normalization
- [x] **Input Event Loop**: Continuous reading and processing of tablet data
  - ✅ Background goroutine processes tablet input when connected

### Phase 3: Drawing System
- [x] **Create Stroke Type**: Define structure for drawing paths
  ```go
  type Stroke struct {
      Points []Point
      Thickness float64
      Color color.Color // Always black per spec
  }
  ```
- [x] **Real-time Rendering**: Update canvas as user draws
  - ✅ Custom DrawingArea widget with real-time refresh
- [x] **Pressure Mapping**: Convert pressure values to line thickness
  - ✅ Pressure affects line width from min to max thickness
- [ ] **Line Smoothing**: Interpolate between points for smooth curves
  - ⚠️ Basic line rendering implemented, smoothing can be enhanced

### Phase 4: User Interface
- [x] **Minimal GUI Design**: 
  - ✅ 1200x900 windowed canvas (primary area)
  - ✅ Small clear button
  - ✅ Save button
  - ✅ Quit button
  - ✅ Minimal toolbar
- [ ] **Keyboard Shortcuts**:
  - Ctrl+Z: Undo
  - Ctrl+Y: Redo
  - Ctrl+S: Save
  - Ctrl+N: Clear canvas
- [x] **No Tool Palette**: Single black marker only
  - ✅ Only black strokes, no tool selection

### Phase 5: File Operations
- [ ] **Save Functionality**: Export drawings as PNG/JPG
- [ ] **Canvas Clearing**: Single-click clear with confirmation
- [ ] **Basic Undo/Redo**: Store stroke history
- [ ] **Auto-save**: Periodic backup of work

### Phase 6: Application Structure
- [ ] **Main Function**: Initialize GUI and start application loop
- [ ] **Application Lifecycle**: Handle startup, shutdown, and errors
- [ ] **Error Handling**: User-friendly error messages
- [ ] **Performance Optimization**: Ensure responsive drawing

## Required Dependencies

### Add to go.mod
```go
require (
    fyne.io/fyne/v2 v2.4.0           // Cross-platform GUI
    github.com/karalabe/hid v1.0.0    // Already present - tablet communication
    image v0.0.0                      // Image processing for saves
    image/png v0.0.0                  // PNG export
    image/jpeg v0.0.0                 // JPEG export
)
```

## Proposed Code Organization

### Current Structure
```
main.go          (single file with basic tablet connection)
```

### Target Structure
```
main.go                    (application entry point)
internal/
├── tablet/
│   ├── controller.go      (tablet communication)
│   ├── input.go          (input processing)
│   └── coordinates.go    (coordinate transformation)
├── drawing/
│   ├── canvas.go         (drawing surface)
│   ├── stroke.go         (stroke representation)
│   └── renderer.go       (drawing operations)
├── ui/
│   ├── window.go         (main window)
│   ├── toolbar.go        (minimal controls)
│   └── events.go         (event handling)
└── file/
    ├── save.go           (export functionality)
    └── autosave.go       (automatic backups)
```

## Key Technical Challenges

### 1. Real-time Performance
- **Challenge**: Drawing must feel responsive with minimal latency
- **Solution**: Efficient rendering pipeline, optimize drawing operations

### 2. Pressure Mapping
- **Challenge**: Convert raw tablet pressure to meaningful line thickness
- **Solution**: Calibration system, smooth pressure curves

### 3. Coordinate Transformation
- **Challenge**: Map tablet coordinates accurately to screen space
- **Solution**: Proper scaling and offset calculations

### 4. Cross-platform Compatibility
- **Challenge**: Ensure consistent behavior on macOS, Windows, Linux
- **Solution**: Use cross-platform GUI framework, test on all platforms

## Implementation Phases

### Phase 1 Priority (MVP)
1. Choose and integrate GUI framework
2. Extend TabletController to read pen data
3. Create basic canvas that responds to pen input
4. Implement simple line drawing

### Phase 2 Priority (Core Features)
1. Add pressure sensitivity
2. Implement clear and save functions
3. Add basic undo functionality
4. Polish user interface

### Phase 3 Priority (Polish)
1. Add keyboard shortcuts
2. Implement auto-save
3. Performance optimizations
4. Error handling improvements

## Success Metrics

The implementation will be considered successful when:
- [ ] User can draw immediately after launching the app
- [ ] Drawing feels natural and responsive (< 10ms latency)
- [ ] Interface stays out of the way (minimal UI)
- [ ] Saving and clearing works intuitively
- [ ] App launches quickly (< 3 seconds)
- [ ] Works consistently across target platforms

## Next Steps

1. **Choose GUI Framework**: Research and select the best cross-platform option
2. **Set up Development Environment**: Configure dependencies and build system
3. **Implement Basic Canvas**: Start with simple drawing functionality
4. **Integrate Tablet Input**: Connect existing tablet code to canvas
5. **Iterate and Test**: Continuous testing with actual XP-Pen tablet
